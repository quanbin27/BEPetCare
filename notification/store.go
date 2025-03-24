package main

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

// MySQLNotificationStore triển khai NotificationStore với GORM và MySQL
type MySQLNotificationStore struct {
	db *gorm.DB
}

// NewMySQLNotificationStore tạo instance mới
func NewMySQLNotificationStore(db *gorm.DB) *MySQLNotificationStore {
	// AutoMigrate để tạo bảng nếu chưa tồn tại
	err := db.AutoMigrate(&EmailNotification{})
	if err != nil {
		panic("failed to migrate notifications table: " + err.Error())
	}
	return &MySQLNotificationStore{db: db}
}

// SaveNotification lưu thông báo vào cơ sở dữ liệu
func (s *MySQLNotificationStore) SaveNotification(ctx context.Context, notification *EmailNotification) error {
	if notification == nil {
		return errors.New("notification cannot be nil")
	}
	// Nếu ID trống, gán giá trị mặc định
	if notification.ID == "" {
		notification.ID = generateUUID()
		notification.CreatedAt = time.Now()
		notification.Status = "pending"
	}

	err := s.db.WithContext(ctx).Create(notification).Error
	if err != nil {
		return err
	}
	return nil
}

// GetNotification lấy thông báo theo ID
func (s *MySQLNotificationStore) GetNotification(ctx context.Context, id string) (*EmailNotification, error) {
	var notification EmailNotification
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&notification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Không tìm thấy thì trả nil
		}
		return nil, err
	}
	return &notification, nil
}

// UpdateNotificationStatus cập nhật trạng thái thông báo
func (s *MySQLNotificationStore) UpdateNotificationStatus(ctx context.Context, id, status string) error {
	err := s.db.WithContext(ctx).
		Model(&EmailNotification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": status,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

func generateUUID() string {
	return uuid.New().String()
}
