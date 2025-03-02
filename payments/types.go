package main

import (
	"context"
	pb "github.com/quanbin27/commons/genproto/payments"
	"time"
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

// PaymentService - Interface cho logic xử lý thanh toán
type PaymentService interface {
	CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error)
	GetPaymentInfo(ctx context.Context, req *pb.GetPaymentInfoRequest) (*pb.GetPaymentInfoResponse, error)
	CreatePaymentURL(ctx context.Context, req *pb.CreatePaymentURLRequest) (*pb.CreatePaymentURLResponse, error)
	CancelPaymentLink(ctx context.Context, req *pb.CancelPaymentLinkRequest) (*pb.CancelPaymentLinkResponse, error)
	UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.UpdatePaymentStatusResponse, error)
	UpdatePaymentMethod(ctx context.Context, req *pb.UpdatePaymentMethodRequest) (*pb.UpdatePaymentMethodResponse, error)
	UpdatePaymentAmount(ctx context.Context, req *pb.UpdatePaymentAmountRequest) (*pb.UpdatePaymentAmountResponse, error)
}
