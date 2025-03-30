package main

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
	"sync"

	pb "github.com/quanbin27/commons/genproto/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGrpcHandler struct {
	orderService OrderService
	pb.UnimplementedOrderServiceServer
	kafkaConn *kafka.Conn
	kafkaMu   sync.Mutex
}

func NewOrderGrpcHandler(grpc *grpc.Server, orderService OrderService, kafkaAddr string) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddr, "order_topic", 0)
	if err != nil {
		log.Fatalf("Failed to dial Kafka leader: %v", err)
	}
	grpcHandler := &OrderGrpcHandler{
		orderService: orderService,
		kafkaConn:    conn,
	}
	pb.RegisterOrderServiceServer(grpc, grpcHandler)
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
		}
	}

	orderID, statusMsg, err := h.orderService.CreateOrder(ctx, req.CustomerId, req.BranchId, items, &req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	go func() {
		h.kafkaMu.Lock()
		defer h.kafkaMu.Unlock()

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

		_, err = h.kafkaConn.WriteMessages(
			kafka.Message{
				Key:   []byte(strconv.FormatInt(int64(orderID), 10)),
				Value: orderJSON,
			},
		)
		if err != nil {
			log.Printf("Failed to write message to Kafka: %v", err)
		} else {
			log.Printf("Order %d sent to Kafka", orderID)
		}
	}()
	return &pb.CreateOrderResponse{OrderId: orderID, Status: statusMsg}, nil
}
func (h *OrderGrpcHandler) Close() {
	if err := h.kafkaConn.Close(); err != nil {
		log.Printf("Failed to close Kafka connection: %v", err)
	}
}
func (h *OrderGrpcHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := h.orderService.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &pb.GetOrderResponse{Order: toPbOrder(order)}, nil
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
