package main

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	store OrderStore
}

func NewOrderService(store OrderStore) OrderService {
	return &Service{store: store}
}

// CreateOrder tạo đơn hàng mới
func (s *Service) CreateOrder(ctx context.Context, customerID, branchID int32, items []OrderItem, appointmentID *int32, pickupTime *time.Time) (int32, string, error) {
	order := &Order{
		CustomerID:    customerID,
		BranchID:      branchID,
		AppointmentID: appointmentID,
		Status:        OrderStatusPending,
		CreatedAt:     time.Now(),
		Items:         items,
		PickupTime:    pickupTime,
	}
	var total float32 = 0
	for _, item := range items {
		total += float32(item.Quantity) * item.UnitPrice
	}
	order.TotalPrice = total
	println("total price:", total)
	// Lưu vào database
	if err := s.store.CreateOrder(ctx, order); err != nil {
		return 0, "Failed", err
	}

	return order.ID, "Success", nil
}

// GetOrder lấy đơn hàng theo ID
func (s *Service) GetOrder(ctx context.Context, orderID int32) (*Order, error) {
	order, err := s.store.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}
func (s *Service) GetOrderByAppointmentID(ctx context.Context, appointmentID int32) (*Order, error) {
	order, err := s.store.GetOrderByAppointmentID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

// UpdateOrderStatus cập nhật trạng thái đơn hàng
func (s *Service) UpdateOrderStatus(ctx context.Context, orderID int32, status OrderStatus) (string, error) {
	if err := s.store.UpdateOrderStatus(ctx, orderID, status); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

// GetOrderItems lấy danh sách sản phẩm trong đơn hàng
func (s *Service) GetOrderItems(ctx context.Context, orderID int32) ([]OrderItem, error) {
	return s.store.GetOrderItems(ctx, orderID)
}
func (s *Service) GetOrdersByCustomerID(ctx context.Context, customerID int32) ([]Order, error) {
	orders, err := s.store.GetOrdersByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
