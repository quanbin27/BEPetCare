package main

import (
	"context"
	"time"

	"github.com/quanbin27/commons/genproto/orders"
)

// OrderStore interface làm việc với database
type OrderStore interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, orderID int32) (*Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int32, status OrderStatus) error
	GetOrderItems(ctx context.Context, orderID int32) ([]OrderItem, error)
}

// OrderService interface cho gRPC
type OrderService interface {
	CreateOrder(ctx context.Context, req *orders.CreateOrderRequest) (*orders.CreateOrderResponse, error)
	GetOrder(ctx context.Context, req *orders.GetOrderRequest) (*orders.GetOrderResponse, error)
	UpdateOrderStatus(ctx context.Context, req *orders.UpdateOrderStatusRequest) (*orders.UpdateOrderStatusResponse, error)
	GetOrderItems(ctx context.Context, req *orders.GetOrderItemsRequest) (*orders.GetOrderItemsResponse, error)
}

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
	ID         int32       `gorm:"primaryKey"`
	CustomerID int32       `gorm:"index;not null"`
	BranchID   int32       `gorm:"index;not null"`
	TotalPrice float32     `gorm:"not null"`
	Status     OrderStatus `gorm:"not null"`
	CreatedAt  time.Time   `gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime"`
	Items      []OrderItem `gorm:"foreignKey:OrderID"`
}

// OrderItem đại diện cho sản phẩm trong đơn hàng
type OrderItem struct {
	ID        int32   `gorm:"primaryKey"`
	OrderID   int32   `gorm:"index;not null"`
	ProductID int32   `gorm:"index;not null"`
	Quantity  int32   `gorm:"not null"`
	UnitPrice float32 `gorm:"not null"`
	Total     float32 `gorm:"not null"`
}
