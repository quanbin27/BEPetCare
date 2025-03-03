package main

import (
	"context"
	"errors"
)

type PaymentServiceImpl struct {
	store PaymentStore
}

func NewPaymentService(store PaymentStore) PaymentService {
	return &PaymentServiceImpl{store: store}
}

func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, orderID, appointmentID int32, amount float32, description string, method PaymentMethod) (int32, error) {
	payment := &Payment{
		OrderID:       orderID,
		AppointmentID: appointmentID,
		Amount:        amount,
		Description:   description,
		Status:        PaymentStatusPending, // Mặc định là PENDING
		Method:        method,
	}
	return s.store.CreatePayment(ctx, payment)
}

func (s *PaymentServiceImpl) GetPaymentInfo(ctx context.Context, paymentID int32) (*Payment, error) {
	return s.store.GetPaymentByID(ctx, paymentID)
}

func (s *PaymentServiceImpl) CreatePaymentURL(ctx context.Context, paymentID int32, amount float32, description string) (string, string, error) {
	// Giả lập logic tạo URL thanh toán (có thể tích hợp PayOS hoặc dịch vụ khác)
	// Đây là giả lập, trả về payment_link_id và checkout_url
	paymentLinkID := "fake-link-id-" + string(rune(paymentID))
	checkoutURL := "https://fake-checkout-url.com/" + paymentLinkID
	return paymentLinkID, checkoutURL, nil
}

func (s *PaymentServiceImpl) CancelPaymentLink(ctx context.Context, paymentID int32, cancellationReason string) (string, error) {
	// Giả lập logic hủy link thanh toán
	// Cập nhật trạng thái thanh toán thành CANCELLED
	err := s.store.UpdatePaymentStatus(ctx, paymentID, PaymentStatusCancelled)
	if err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *PaymentServiceImpl) UpdatePaymentStatus(ctx context.Context, paymentID int32, status PaymentStatus) (string, error) {
	err := s.store.UpdatePaymentStatus(ctx, paymentID, status)
	if err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *PaymentServiceImpl) UpdatePaymentMethod(ctx context.Context, paymentID int32, method PaymentMethod) (string, error) {
	err := s.store.UpdatePaymentMethod(ctx, paymentID, method)
	if err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *PaymentServiceImpl) UpdatePaymentAmount(ctx context.Context, paymentID int32, amount float32) (string, error) {
	if amount <= 0 {
		return "Failed", errors.New("invalid amount")
	}
	err := s.store.UpdatePaymentAmount(ctx, paymentID, amount)
	if err != nil {
		return "Failed", err
	}
	return "Success", nil
}
