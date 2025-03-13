package main

import (
	"context"

	pb "github.com/quanbin27/commons/genproto/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGrpcHandler struct {
	orderService OrderService
	pb.UnimplementedOrderServiceServer
}

func NewOrderGrpcHandler(grpc *grpc.Server, orderService OrderService) {
	grpcHandler := &OrderGrpcHandler{
		orderService: orderService,
	}
	pb.RegisterOrderServiceServer(grpc, grpcHandler)
}

func (h *OrderGrpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	// Chuyển đổi OrderItem từ protobuf sang nội bộ
	items := make([]OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	orderID, statusMsg, err := h.orderService.CreateOrder(ctx, req.CustomerId, req.BranchId, items, &req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateOrderResponse{OrderId: orderID, Status: statusMsg}, nil
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
