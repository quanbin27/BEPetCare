package main

import (
	"context"
	"encoding/json"
	pb "github.com/quanbin27/commons/genproto/orders"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

type OrderGrpcHandler struct {
	orderService OrderService
	pb.UnimplementedOrderServiceServer
	kafkaWriter *kafka.Writer
}

func NewOrderGrpcHandler(grpc *grpc.Server, orderService OrderService, kafkaAddr string) *OrderGrpcHandler {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaAddr),
		Topic:    "order_topic",
		Balancer: &kafka.LeastBytes{}, // hoặc RoundRobin nếu bạn muốn chia đều
		Async:    false,               // true nếu bạn chấp nhận gửi async
	}

	handler := &OrderGrpcHandler{
		orderService: orderService,
		kafkaWriter:  writer,
	}

	pb.RegisterOrderServiceServer(grpc, handler)
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

	orderID, statusMsg, err := h.orderService.CreateOrder(ctx, req.CustomerId, req.BranchId, items, &req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

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
		pbOrders[i] = toPbOrder(&order)
	}
	return &pb.GetOrdersByCustomerIDResponse{Orders: pbOrders}, nil
}
