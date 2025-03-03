package main

import (
	"context"
	"time"

	pb "github.com/quanbin27/commons/genproto/payments"
)

// PaymentStatus - Trạng thái thanh toán trong database
type PaymentStatus string
type PaymentMethod string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

const (
	PaymentMethodCash PaymentMethod = "CASH"
	PaymentMethodBank PaymentMethod = "BANK"
)

// Payment - Bảng lưu trữ thông tin thanh toán
type Payment struct {
	ID            int32 `gorm:"primaryKey"`
	OrderID       int32
	AppointmentID int32
	Amount        float32       `gorm:"not null"`
	Description   string        `gorm:"type:text"`
	Status        PaymentStatus `gorm:"type:varchar(20);not null"`
	Method        PaymentMethod `gorm:"type:varchar(20);not null"`
	CreatedAt     time.Time     `gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime"`
}

// PaymentStore - Interface tương tác với database
type PaymentStore interface {
	CreatePayment(ctx context.Context, payment *Payment) (int32, error)
	GetPaymentByID(ctx context.Context, paymentID int32) (*Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID int32, status PaymentStatus) error
	UpdatePaymentMethod(ctx context.Context, paymentID int32, method PaymentMethod) error
	UpdatePaymentAmount(ctx context.Context, paymentID int32, amount float32) error
}

// PaymentService - Interface cho logic xử lý thanh toán với dữ liệu nội bộ
type PaymentService interface {
	CreatePayment(ctx context.Context, orderID, appointmentID int32, amount float32, description string, method PaymentMethod) (int32, error)
	GetPaymentInfo(ctx context.Context, paymentID int32) (*Payment, error)
	CreatePaymentURL(ctx context.Context, paymentID int32, amount float32, description string) (string, string, error) // Trả về payment_link_id, checkout_url, error
	CancelPaymentLink(ctx context.Context, paymentID int32, cancellationReason string) (string, error)                 // Trả về status
	UpdatePaymentStatus(ctx context.Context, paymentID int32, status PaymentStatus) (string, error)                    // Trả về status
	UpdatePaymentMethod(ctx context.Context, paymentID int32, method PaymentMethod) (string, error)                    // Trả về status
	UpdatePaymentAmount(ctx context.Context, paymentID int32, amount float32) (string, error)                          // Trả về status
}

// Helper functions to convert between internal types and protobuf types
func toProtoPaymentStatus(status PaymentStatus) pb.PaymentStatus {
	switch status {
	case PaymentStatusPending:
		return pb.PaymentStatus_PENDING
	case PaymentStatusCompleted:
		return pb.PaymentStatus_COMPLETED
	case PaymentStatusFailed:
		return pb.PaymentStatus_FAILED
	case PaymentStatusCancelled:
		return pb.PaymentStatus_CANCELLED
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

func fromProtoPaymentStatus(status pb.PaymentStatus) PaymentStatus {
	switch status {
	case pb.PaymentStatus_PENDING:
		return PaymentStatusPending
	case pb.PaymentStatus_COMPLETED:
		return PaymentStatusCompleted
	case pb.PaymentStatus_FAILED:
		return PaymentStatusFailed
	case pb.PaymentStatus_CANCELLED:
		return PaymentStatusCancelled
	default:
		return PaymentStatusPending // Mặc định là PENDING nếu không xác định
	}
}

func toProtoPaymentMethod(method PaymentMethod) pb.PaymentMethod {
	switch method {
	case PaymentMethodCash:
		return pb.PaymentMethod_CASH
	case PaymentMethodBank:
		return pb.PaymentMethod_BANK
	default:
		return pb.PaymentMethod_CASH // Mặc định là CASH nếu không xác định
	}
}

func fromProtoPaymentMethod(method pb.PaymentMethod) PaymentMethod {
	switch method {
	case pb.PaymentMethod_CASH:
		return PaymentMethodCash
	case pb.PaymentMethod_BANK:
		return PaymentMethodBank
	default:
		return PaymentMethodCash // Mặc định là CASH nếu không xác định
	}
}

func toProtoPayment(p *Payment) *pb.GetPaymentInfoResponse {
	return &pb.GetPaymentInfoResponse{
		Status:        toProtoPaymentStatus(p.Status),
		Method:        toProtoPaymentMethod(p.Method),
		OrderId:       p.OrderID,
		AppointmentId: p.AppointmentID,
		Amount:        p.Amount,
		Description:   p.Description,
	}
}
