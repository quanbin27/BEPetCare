package main

import (
	"context"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) PaymentStore {
	return &Store{db: db}
}

// CreatePayment lưu thông tin thanh toán vào cơ sở dữ liệu
func (s *Store) CreatePayment(ctx context.Context, payment *Payment) error {
	return s.db.WithContext(ctx).Create(payment).Error
}

// GetPaymentByID lấy thông tin thanh toán theo ID
func (s *Store) GetPaymentByID(ctx context.Context, id string) (*Payment, error) {
	var payment Payment
	if err := s.db.WithContext(ctx).First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// UpdatePaymentStatus cập nhật trạng thái thanh toán
func (s *Store) UpdatePaymentStatus(ctx context.Context, transactionID string, status string) error {
	return s.db.WithContext(ctx).Model(&Payment{}).
		Where("transaction_id = ?", transactionID).
		Update("status", status).Error
}

// GetPaymentsByUser lấy danh sách thanh toán theo ID khách hàng
func (s *Store) GetPaymentsByUser(ctx context.Context, userID int32) ([]Payment, error) {
	var payments []Payment
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// GetPaymentsByOrder lấy danh sách thanh toán theo ID đơn hàng
func (s *Store) GetPaymentsByOrder(ctx context.Context, orderID int32) ([]Payment, error) {
	var payments []Payment
	if err := s.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// GetPaymentsByAppointment lấy danh sách thanh toán theo ID lịch hẹn
func (s *Store) GetPaymentsByAppointment(ctx context.Context, appointmentID int32) ([]Payment, error) {
	var payments []Payment
	if err := s.db.WithContext(ctx).Where("appointment_id = ?", appointmentID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}
