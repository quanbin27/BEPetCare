package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/payOSHQ/payos-lib-golang"
	"log"
	"strconv"
	"time"
)

type PaymentServiceImpl struct {
	store PaymentStore
}

func NewPaymentService(store PaymentStore) PaymentService {
	return &PaymentServiceImpl{store: store}
}

func (s *PaymentServiceImpl) CreatePayment(ctx context.Context, orderID, appointmentID int32, amount float32, description string, method PaymentMethod) (int32, error) {
	if orderID == 0 && appointmentID == 0 {
		return 0, fmt.Errorf("phải cung cấp ít nhất OrderID hoặc AppointmentID")
	}
	if amount <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}

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

func (s *PaymentServiceImpl) CreatePaymentURL(ctx context.Context, paymentID int32) (string, string, error) {
	payment, err := s.store.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get payment info: %w", err)
	}
	if payment.Method != PaymentMethodBank {
		return "", "", fmt.Errorf("phương thức thanh toán phải là BANK để tạo link PayOS")
	}
	if payment.OrderCode != 0 && payment.PaymentLinkID != "" {
		// Kiểm tra trạng thái link cũ
		data, err := payos.GetPaymentLinkInformation(strconv.FormatInt(payment.OrderCode, 10))
		if err == nil && data.Status == "PENDING" {
			// Link cũ vẫn hợp lệ, trả về PaymentLinkID và CheckoutUrl hiện tại
			return payment.PaymentLinkID, payment.CheckoutURL, nil
		}
		// Nếu link cũ không hợp lệ (hết hạn, bị hủy, hoặc thất bại), tạo link mới
	}
	if payment.Status != PaymentStatusPending && payment.Status != PaymentStatusFailed {
		return "", "", fmt.Errorf("payment phải ở trạng thái PENDING hoặc FAILED để tạo link thanh toán")
	}
	var orderCode int64
	if payment.OrderCode == 0 {
		if payment.Method == PaymentMethodBank {
			orderCode = generateOrderCode(payment.OrderID, payment.AppointmentID)
		}
		if err := s.store.UpdateOrderCode(ctx, paymentID, orderCode); err != nil {
			return "", "", fmt.Errorf("lỗi cập nhật OrderCode: %v", err)
		}
	}
	log.Printf("Creating payment link for OrderCode: %d, Amount: %.2f, Description: %s", orderCode, payment.Amount, payment.Description)
	body := payos.CheckoutRequestType{
		OrderCode:   orderCode,
		Amount:      int(payment.Amount),
		Description: payment.Description,
		CancelUrl:   "http://26.30.229.237:8080/payments/cancel",
		ReturnUrl:   "http://26.30.229.237:8080/payments/success",
	}
	data, err := payos.CreatePaymentLink(body)
	if err != nil {
		return "", "", fmt.Errorf("lỗi tạo link thanh toán PayOS: %v", err)
	}
	err = s.store.UpdateCheckoutURL(ctx, paymentID, data.CheckoutUrl, data.PaymentLinkId)
	if err != nil {
		return "", "", fmt.Errorf("lỗi cập nhật URL thanh toán: %v", err)
	}
	return data.CheckoutUrl, data.PaymentLinkId, nil
}

func (s *PaymentServiceImpl) CancelPaymentLink(ctx context.Context, paymentID int32, cancellationReason string) (string, error) {
	payment, err := s.store.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return "", fmt.Errorf("lỗi lấy thông tin payment: %v", err)
	}

	// Chỉ áp dụng cho BANK
	if payment.Method != PaymentMethodBank {
		return "", fmt.Errorf("phương thức thanh toán phải là BANK để hủy link PayOS")
	}
	if payment.Status != PaymentStatusPending {
		return "", fmt.Errorf("payment không ở trạng thái PENDING")
	}
	if payment.OrderCode == 0 {
		return "", fmt.Errorf("payment không có OrderCode hợp lệ")
	}
	_, err = payos.CancelPaymentLink(strconv.FormatInt(payment.OrderCode, 10), &cancellationReason)
	if err != nil {
		return "", fmt.Errorf("lỗi hủy link thanh toán PayOS: %v", err)
	}
	err = s.store.UpdatePaymentStatus(ctx, paymentID, PaymentStatusCancelled)
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
	payment, err := s.store.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return "Thất bại", fmt.Errorf("lỗi lấy thông tin payment: %v", err)
	}
	if payment.PaymentLinkID != "" {
		return "Thất bại", fmt.Errorf("không thể cập nhật phương thức cho payment đã có link PayOS")
	}

	if err := s.store.UpdatePaymentMethod(ctx, paymentID, method); err != nil {
		return "Thất bại", fmt.Errorf("lỗi cập nhật phương thức thanh toán: %v", err)
	}
	return "Thành công", nil
}
func (s *PaymentServiceImpl) UpdateBankPaymentStatus(ctx context.Context, orderCode int64, status PaymentStatus) (string, error) {
	payment, err := s.store.GetPaymentByOrderCode(ctx, orderCode)
	if err != nil {
		return "Thất bại", fmt.Errorf("lỗi lấy thông tin payment: %v", err)
	}
	if payment.Status != "PENDING" {
		return "Thất bại", fmt.Errorf("không thể cập nhật CHO giao dich đã hoàn thành hoặc thất bại")
	}

	if err := s.store.UpdatePaymentStatus(ctx, payment.ID, status); err != nil {
		return "Thất bại", fmt.Errorf("lỗi cập nhật phương thức thanh toán: %v", err)
	}
	return "Thành công", nil
}

func (s *PaymentServiceImpl) UpdatePaymentAmount(ctx context.Context, paymentID int32, amount float32) (string, error) {
	payment, err := s.store.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return "Thất bại", fmt.Errorf("lỗi lấy thông tin payment: %v", err)
	}
	if payment.PaymentLinkID != "" {
		return "Thất bại", fmt.Errorf("không thể cập nhật số tiền cho payment đã có link PayOS")
	}
	if amount <= 0 {
		return "Thất bại", fmt.Errorf("số tiền không hợp lệ")
	}

	if err := s.store.UpdatePaymentAmount(ctx, paymentID, amount); err != nil {
		return "Thất bại", fmt.Errorf("lỗi cập nhật số tiền: %v", err)
	}
	return "Thành công", nil
}
func generateOrderCode(orderID, appointmentID int32) int64 {
	timestamp := time.Now().UnixMilli() % 1_000_000_000 // ví dụ: 1749538578566

	// Limit orderID and appointmentID to 3 digits
	suffix := int64(orderID%1000)*1000 + int64(appointmentID%1000)

	return timestamp*1_000_000 + suffix // timestamp (13 digits) + suffix (6 digits) = max 19 digits
}
