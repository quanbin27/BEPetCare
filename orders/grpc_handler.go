package main

import (
	"context"
	"encoding/json"
	"github.com/quanbin27/commons/config"
	pb "github.com/quanbin27/commons/genproto/orders"
	pbProduct "github.com/quanbin27/commons/genproto/products"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"time"
)

type OrderGrpcHandler struct {
	orderService OrderService
	pb.UnimplementedOrderServiceServer
	kafkaWriter   *kafka.Writer
	productClient pbProduct.ProductServiceClient
}

func NewOrderGrpcHandler(grpcServer *grpc.Server, orderService OrderService, kafkaAddr string) *OrderGrpcHandler {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaAddr),
		Topic:    "order_topic",
		Balancer: &kafka.LeastBytes{}, // hoặc RoundRobin nếu bạn muốn chia đều
		Async:    false,               // true nếu bạn chấp nhận gửi async
	}
	productsServiceAddr := config.Envs.ProductsGrpcAddr
	productsConn, err := grpc.NewClient(productsServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}
	handler := &OrderGrpcHandler{
		orderService:  orderService,
		kafkaWriter:   writer,
		productClient: pbProduct.NewProductServiceClient(productsConn),
	}

	pb.RegisterOrderServiceServer(grpcServer, handler)
	return handler
}

func (h *OrderGrpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Chuyển đổi OrderItem từ protobuf sang nội bộ
	items := make([]OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = OrderItem{
			ProductID:   item.ProductId,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			ProductType: item.ProductType,
			ProductName: item.ProductName,
		}
	}

	var pickupTime *time.Time
	if req.PickupTime != "" {
		t, err := time.Parse(time.RFC3339, req.PickupTime)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid pickup_time format: %v", err)
		}
		pickupTime = &t
	}

	// ====== Gọi ReserveProduct cho từng item ======
	var reservedItems []*pbProduct.ReserveProductRequest
	for _, item := range req.Items {
		reserveReq := &pbProduct.ReserveProductRequest{
			ProductId:   item.ProductId,
			ProductType: item.ProductType,
			Quantity:    item.Quantity,
			BranchId:    req.BranchId,
		}
		_, err := h.productClient.ReserveProduct(ctx, reserveReq)
		if err != nil {
			// Rollback các item đã reserve trước đó
			for _, reserved := range reservedItems {
				_, _ = h.productClient.ReleaseReservation(ctx, &pbProduct.ReleaseReservationRequest{
					ProductId:   reserved.ProductId,
					ProductType: reserved.ProductType,
					Quantity:    reserved.Quantity,
					BranchId:    reserved.BranchId,
				})
			}
			return nil, status.Errorf(codes.Internal, "failed to reserve product: %v", err)
		}
		reservedItems = append(reservedItems, reserveReq)
	}

	// ====== Gọi orderService để tạo order ======
	orderID, statusMsg, err := h.orderService.CreateOrder(ctx, req.CustomerId, req.BranchId, items, &req.AppointmentId, pickupTime)
	if err != nil {
		// Rollback nếu tạo order thất bại
		for _, reserved := range reservedItems {
			_, _ = h.productClient.ReleaseReservation(ctx, &pbProduct.ReleaseReservationRequest{
				ProductId:   reserved.ProductId,
				ProductType: reserved.ProductType,
				Quantity:    reserved.Quantity,
				BranchId:    reserved.BranchId,
			})
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// ====== Gửi Kafka async ======
	go func() {
		orderData := map[string]interface{}{
			"order_id":    orderID,
			"customer_id": req.CustomerId,
			"branch_id":   req.BranchId,
			"items":       items,
			"status":      statusMsg,
			"email":       req.Email,
		}
		orderJSON, err := json.Marshal(orderData)
		if err != nil {
			log.Printf("Failed to marshal order data: %v", err)
			return
		}
		msg := kafka.Message{
			Key:   []byte(strconv.FormatInt(int64(orderID), 10)),
			Value: orderJSON,
		}
		if err := h.kafkaWriter.WriteMessages(context.Background(), msg); err != nil {
			log.Printf("Failed to write message to Kafka: %v", err)
		} else {
			log.Printf("Order %d sent to Kafka", orderID)
		}
	}()

	return &pb.CreateOrderResponse{OrderId: orderID, Status: statusMsg}, nil
}

func (h *OrderGrpcHandler) Close() {
	if err := h.kafkaWriter.Close(); err != nil {
		log.Printf("Failed to close Kafka writer: %v", err)
	}
}

func (h *OrderGrpcHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := h.orderService.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &pb.GetOrderResponse{Order: toPbOrder(order)}, nil
}
func (h *OrderGrpcHandler) GetOrderByAppointmentID(ctx context.Context, req *pb.GetOrderByAppointmentIDRequest) (*pb.GetOrderByAppointmentIDResponse, error) {
	order, err := h.orderService.GetOrderByAppointmentID(ctx, req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &pb.GetOrderByAppointmentIDResponse{Order: toPbOrder(order)}, nil
}
func (h *OrderGrpcHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	statusMsg, err := h.orderService.UpdateOrderStatus(ctx, req.OrderId, fromPbOrderStatus(req.Status))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdateOrderStatusResponse{Status: statusMsg}, nil
}

func (h *OrderGrpcHandler) GetOrderItems(ctx context.Context, req *pb.GetOrderItemsRequest) (*pb.GetOrderItemsResponse, error) {
	items, err := h.orderService.GetOrderItems(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pbItems := make([]*pb.OrderItem, len(items))
	for i, item := range items {
		pbItems[i] = toPbOrderItem(item)
	}
	return &pb.GetOrderItemsResponse{Items: pbItems}, nil
}
func (h *OrderGrpcHandler) GetOrdersByCustomerID(ctx context.Context, req *pb.GetOrdersByCustomerIDRequest) (*pb.GetOrdersByCustomerIDResponse, error) {
	orders, err := h.orderService.GetOrdersByCustomerID(ctx, req.CustomerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pbOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		pbOrders[len(orders)-1-i] = toPbOrder(&order)
	}
	return &pb.GetOrdersByCustomerIDResponse{Orders: pbOrders}, nil
}
