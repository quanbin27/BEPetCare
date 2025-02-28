package main

import (
	"context"
	"google.golang.org/grpc"

	pb "github.com/quanbin27/commons/genproto/orders"
)

// OrdersGrpcHandler triển khai gRPC handler cho OrderService
type OrdersGrpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
}

// NewGrpcOrderHandler đăng ký gRPC handler vào server
func NewGrpcOrderHandler(grpcServer *grpc.Server, service OrderService) {
	pb.RegisterOrderServiceServer(grpcServer, &OrdersGrpcHandler{service: service})
}

// CreateOrder xử lý tạo đơn hàng
func (h *OrdersGrpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	return h.service.CreateOrder(ctx, req)
}

// GetOrder xử lý lấy đơn hàng theo ID
func (h *OrdersGrpcHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	return h.service.GetOrder(ctx, req)
}

// UpdateOrderStatus xử lý cập nhật trạng thái đơn hàng
func (h *OrdersGrpcHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	return h.service.UpdateOrderStatus(ctx, req)
}

// GetOrderItems xử lý lấy danh sách sản phẩm trong đơn hàng
func (h *OrdersGrpcHandler) GetOrderItems(ctx context.Context, req *pb.GetOrderItemsRequest) (*pb.GetOrderItemsResponse, error) {
	return h.service.GetOrderItems(ctx, req)
}
