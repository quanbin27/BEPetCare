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
	CreatedAt time.Time `gorm:"not null"`
	SentAt    time.Time // Có thể null trong MySQL
	Status    string    `gorm:"type:varchar(50);not null"`
}

// NotificationStore interface
type NotificationStore interface {
	SaveNotification(ctx context.Context, notification *EmailNotification) error
	GetNotification(ctx context.Context, id string) (*EmailNotification, error)
	UpdateNotificationStatus(ctx context.Context, id, status string) error
}

// NotificationService interface
type NotificationService interface {
	SendVerificationEmail(ctx context.Context, email, token, baseURL string) (string, error)
}
