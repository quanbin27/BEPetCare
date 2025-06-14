package main

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Store triển khai OrderStore
type Store struct {
	db *gorm.DB
}

// NewOrderStore khởi tạo OrderStore
func NewOrderStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// CreateOrder thêm đơn hàng vào DB
func (s *Store) CreateOrder(ctx context.Context, order *Order) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetOrderByID lấy thông tin đơn hàng theo ID
func (s *Store) GetOrderByID(ctx context.Context, orderID int32) (*Order, error) {
	var order Order
	if err := s.db.WithContext(ctx).Preload("Items").First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}
func (s *Store) GetOrderByAppointmentID(ctx context.Context, appointmentID int32) (*Order, error) {
	var order Order
	if err := s.db.WithContext(ctx).
		Preload("Items").
		Where("appointment_id = ?", appointmentID).
		First(&order).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}
func (s *Store) GetOrdersByCustomerID(ctx context.Context, customerID int32) ([]Order, error) {
	var orders []Order
	if err := s.db.WithContext(ctx).
		Preload("Items").
		Where("customer_id = ?", customerID).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
func (s *Store) GetAllOrders(ctx context.Context) ([]Order, error) {
	var orders []Order
	if err := s.db.WithContext(ctx).
		Preload("Items").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateOrderStatus cập nhật trạng thái đơn hàng
func (s *Store) UpdateOrderStatus(ctx context.Context, orderID int32, status OrderStatus) error {
	return s.db.WithContext(ctx).Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// GetOrderItems lấy danh sách sản phẩm trong đơn hàng
func (s *Store) GetOrderItems(ctx context.Context, orderID int32) ([]OrderItem, error) {
	var items []OrderItem
	if err := s.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
