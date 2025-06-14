package main

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// CreatePayment lưu thông tin thanh toán vào cơ sở dữ liệu
func (s *Store) CreatePayment(ctx context.Context, payment *Payment) (int32, error) {
	if payment == nil {
		return 0, errors.New("payment cannot be nil")
	}
	err := s.db.WithContext(ctx).Create(payment).Error
	if err != nil {
		return 0, err
	}
	return payment.ID, nil
}
func (s *Store) GetPaymentByID(ctx context.Context, paymentID int32) (*Payment, error) {
	var payment Payment
	err := s.db.WithContext(ctx).Where("id = ?", paymentID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}
func (s *Store) GetPaymentByOrderCode(ctx context.Context, orderCode int64) (*Payment, error) {
	var payment Payment
	err := s.db.WithContext(ctx).Where("order_code = ?", orderCode).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

// UpdatePaymentStatus - Cập nhật trạng thái thanh toán
func (s *Store) UpdatePaymentStatus(ctx context.Context, paymentID int32, status PaymentStatus) error {
	return s.db.WithContext(ctx).
		Model(&Payment{}).
		Where("id = ?", paymentID).
		Update("status", status).
		Error
}

// UpdatePaymentMethod - Cập nhật phương thức thanh toán
func (s *Store) UpdatePaymentMethod(ctx context.Context, paymentID int32, method PaymentMethod) error {
	return s.db.WithContext(ctx).
		Model(&Payment{}).
		Where("id = ?", paymentID).
		Update("method", method).
		Error
}
func (s *Store) UpdateOrderCode(ctx context.Context, paymentID int32, orderCode int64) error {
	return s.db.WithContext(ctx).
		Model(&Payment{}).
		Where("id = ?", paymentID).
		Update("order_code", orderCode).
		Error
}

// Update CheckoutURL - Cập nhật URL thanh toán
func (s *Store) UpdateCheckoutURL(ctx context.Context, paymentID int32, checkoutURL, paymentLinkID string) error {
	return s.db.WithContext(ctx).
		Model(&Payment{}).
		Where("id = ?", paymentID).
		Updates(map[string]interface{}{
			"payment_link_id": paymentLinkID,
			"checkout_url":    checkoutURL,
		}).
		Error
}

// UpdatePaymentAmount - Cập nhật số tiền thanh toán
func (s *Store) UpdatePaymentAmount(ctx context.Context, paymentID int32, amount float32) error {
	return s.db.WithContext(ctx).
		Model(&Payment{}).
		Where("id = ?", paymentID).
		Update("amount", amount).
		Error
}
