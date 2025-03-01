package main

import (
	"context"
	"time"

	pb "github.com/quanbin27/commons/genproto/payments"
	"gorm.io/gorm"
)

// PaymentStatus - Trạng thái thanh toán trong database

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusCancelled PaymentStatus = "CANCELLED"
)

// Payment - Bảng lưu trữ thông tin thanh toán

type Payment struct {
	ID            string         `gorm:"primaryKey;type:varchar(36)"`
	CustomerID    string         `gorm:"index;not null"`
	OrderID       *string        `gorm:"index"`
	AppointmentID *string        `gorm:"index"`
	Amount        float32        `gorm:"not null"`
	Description   string         `gorm:"type:text"`
	Status        PaymentStatus  `gorm:"type:varchar(20);not null"`
	TransactionID *string        `gorm:"index"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// PaymentStore - Interface tương tác với database

type PaymentStore interface {
	CreatePayment(ctx context.Context, payment *Payment) error
	GetPaymentByID(ctx context.Context, paymentID string) (*Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID string, status PaymentStatus, transactionID *string) error
}

// PaymentService - Interface cho service xử lý thanh toán

type PaymentService interface {
	CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error)
	GetPaymentStatus(ctx context.Context, req *pb.GetPaymentStatusRequest) (*pb.GetPaymentStatusResponse, error)
	CancelPayment(ctx context.Context, req *pb.CancelPaymentRequest) (*pb.CancelPaymentResponse, error)
}
