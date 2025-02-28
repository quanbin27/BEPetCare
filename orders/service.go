package main

import (
	"context"
	"errors"
	"time"

	pb "github.com/quanbin27/commons/genproto/orders"
)

// Service triển khai OrderService
type Service struct {
	store OrderStore
}

// NewOrderService khởi tạo OrderService
func NewOrderService(store OrderStore) OrderService {
	return &Service{store: store}
}

// PlaceOrder tạo đơn hàng mới
func (s *Service) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := &Order{
		CustomerID: req.CustomerId,
		BranchID:   req.BranchId,
		Status:     OrderStatusPending,
		CreatedAt:  time.Now(),
		Items:      []OrderItem{},
	}
	var total float32 = 0
	for _, item := range req.Items {
		order.Items = append(order.Items, OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
		total += float32(item.Quantity) * item.UnitPrice
	}
	order.TotalPrice = total

	// Lưu vào database
	if err := s.store.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{OrderId: order.ID, Status: "Success"}, nil
}
func toPbOrderStatus(status OrderStatus) pb.OrderStatus {
	switch status {
	case OrderStatusPending:
		return pb.OrderStatus_PENDING
	case OrderStatusPaid:
		return pb.OrderStatus_PAID
	case OrderStatusCompleted:
		return pb.OrderStatus_COMPLETED
	case OrderStatusCancelled:
		return pb.OrderStatus_CANCELLED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED // Giá trị mặc định nếu không khớp
	}
}
func fromPbOrderStatus(pbStatus pb.OrderStatus) OrderStatus {
	switch pbStatus {
	case pb.OrderStatus_PENDING:
		return OrderStatusPending
	case pb.OrderStatus_PAID:
		return OrderStatusPaid
	case pb.OrderStatus_COMPLETED:
		return OrderStatusCompleted
	case pb.OrderStatus_CANCELLED:
		return OrderStatusCancelled
	default:
		return "unknown order status"
	}
}

// GetOrder lấy đơn hàng theo ID
func (s *Service) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.store.GetOrderByID(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	// Chuyển đổi sang protobuf
	pbOrder := &pb.Order{
		Id:         order.ID,
		CustomerId: order.CustomerID,
		BranchId:   order.BranchID,
		TotalPrice: order.TotalPrice,
		Status:     toPbOrderStatus(order.Status),
	}

	for _, item := range order.Items {
		pbOrder.Items = append(pbOrder.Items, &pb.OrderItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	return &pb.GetOrderResponse{Order: pbOrder}, nil
}

// UpdateOrderStatus cập nhật trạng thái đơn hàng
func (s *Service) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	newStatus := fromPbOrderStatus(req.Status)
	if err := s.store.UpdateOrderStatus(ctx, req.OrderId, newStatus); err != nil {
		return &pb.UpdateOrderStatusResponse{Status: "Failed"}, err
	}
	return &pb.UpdateOrderStatusResponse{Status: "Success"}, nil
}

// GetOrderItems lấy danh sách sản phẩm trong đơn hàng
func (s *Service) GetOrderItems(ctx context.Context, req *pb.GetOrderItemsRequest) (*pb.GetOrderItemsResponse, error) {
	items, err := s.store.GetOrderItems(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	var pbItems []*pb.OrderItem
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			ProductId: item.ProductID,
			OrderId:   item.OrderID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	return &pb.GetOrderItemsResponse{Items: pbItems}, nil
}
