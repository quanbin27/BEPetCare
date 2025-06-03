package main

import (
	"context"
	"time"
)

// EmailNotification đại diện cho thông tin email cần gửi
type EmailNotification struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Subject   string    `gorm:"type:varchar(255);not null"`
	Body      string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Status    string    `gorm:"type:varchar(50);not null"`
}

// For later update
//type Notification struct {
//	ID        primitive.ObjectID `bson:"_id,omitempty"`        // MongoDB ObjectID
//	UserID    string             `bson:"user_id"`              // Người nhận
//	Type      string             `bson:"type"`                 // "email", "push", "sms"
//	Email     *EmailContent      `bson:"email,omitempty"`      // Chỉ có nếu là email
//	Push      *PushContent       `bson:"push,omitempty"`       // Chỉ có nếu là push
//	Status    string             `bson:"status"`               // pending, sent, read, failed
//	CreatedAt time.Time          `bson:"created_at"`           // Tự gán khi insert
//	UpdatedAt time.Time          `bson:"updated_at,omitempty"` // Tự cập nhật khi update
//}
//type EmailContent struct {
//	To      string `bson:"to"`      // Email address
//	Subject string `bson:"subject"` // Tiêu đề
//	Body    string `bson:"body"`    // Nội dung
//}
//type PushContent struct {
//	Title   string `bson:"title"`   // Tiêu đề push
//	Message string `bson:"message"` // Nội dung push
//	Token   string `bson:"token"`   // Firebase device token
//}

// NotificationStore interface
type NotificationStore interface {
	SaveNotification(ctx context.Context, notification *EmailNotification) error
	GetNotification(ctx context.Context, id string) (*EmailNotification, error)
	UpdateNotificationStatus(ctx context.Context, id, status string) error
}

// NotificationService interface
type NotificationService interface {
	SendResetPasswordEmail(ctx context.Context, email, token, baseURL string) (string, error)
	SendVerificationEmail(ctx context.Context, email, token, baseURL string) (string, error)
}
type OrderData struct {
	OrderID    int32       `json:"order_id"`
	CustomerID int32       `json:"customer_id"`
	BranchID   int32       `json:"branch_id"`
	Items      []OrderItem `json:"items"`
	Status     string      `json:"status"`
	Email      string      `json:"email"`
}
type OrderItem struct {
	ProductID   int32   `json:"product_id"`
	Quantity    int32   `json:"quantity"`
	UnitPrice   float32 `json:"unit_price"`
	ProductType string  `json:"product_type"`
	ProductName string  `json:"product_name"`
}
