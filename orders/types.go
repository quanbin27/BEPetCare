package main

import (
	"context"
	"time"

	pb "github.com/quanbin27/commons/genproto/orders"
)

// OrderStatus định nghĩa trạng thái đơn hàng
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

// Order đại diện cho đơn hàng
type Order struct {
	ID            int32 `gorm:"primaryKey"`
	CustomerID    int32 `gorm:"index;not null"`
	BranchID      int32 `gorm:"index;not null"`
	AppointmentID *int32
	TotalPrice    float32     `gorm:"not null"`
	Status        OrderStatus `gorm:"not null"`
	CreatedAt     time.Time   `gorm:"autoCreateTime"`
	UpdatedAt     time.Time   `gorm:"autoUpdateTime"`
	Items         []OrderItem `gorm:"foreignKey:OrderID"`
}

// OrderItem đại diện cho sản phẩm trong đơn hàng
type OrderItem struct {
	OrderID     int32   `gorm:"primaryKey"`
	ProductID   int32   `gorm:"primaryKey"`
	ProductType string  `gorm:"primaryKey"`
	Quantity    int32   `gorm:"not null"`
	UnitPrice   float32 `gorm:"not null"`
}

// OrderStore interface làm việc với database
type OrderStore interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, orderID int32) (*Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int32, status OrderStatus) error
	GetOrderItems(ctx context.Context, orderID int32) ([]OrderItem, error)
}

// OrderService interface cho logic xử lý với dữ liệu nội bộ
type OrderService interface {
	CreateOrder(ctx context.Context, customerID, branchID int32, items []OrderItem, appointmentID *int32) (int32, string, error) // Trả về orderID, status
	GetOrder(ctx context.Context, orderID int32) (*Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int32, status OrderStatus) (string, error) // Trả về status
	GetOrderItems(ctx context.Context, orderID int32) ([]OrderItem, error)
}

// Helper functions to convert between internal types and protobuf types
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
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
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
		return OrderStatusPending // Mặc định là PENDING nếu không xác định
	}
}

func toPbOrder(o *Order) *pb.Order {
	pbItems := make([]*pb.OrderItem, len(o.Items))
	for i, item := range o.Items {
		pbItems[i] = &pb.OrderItem{
			OrderId:     item.OrderID,
			ProductId:   item.ProductID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			ProductType: item.ProductType,
		}
	}
	return &pb.Order{
		Id:         o.ID,
		CustomerId: o.CustomerID,
		BranchId:   o.BranchID,
		TotalPrice: o.TotalPrice,
		Status:     toPbOrderStatus(o.Status),
		CreatedAt:  o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  o.UpdatedAt.Format(time.RFC3339),
		Items:      pbItems,
	}
}

func toPbOrderItem(item OrderItem) *pb.OrderItem {
	return &pb.OrderItem{
		OrderId:     item.OrderID,
		ProductId:   item.ProductID,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		ProductType: item.ProductType,
	}
}
