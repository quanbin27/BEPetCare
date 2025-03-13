package main

import (
	"context"
	"fmt"

	"gopkg.in/gomail.v2"
)

// Service triển khai NotificationService
type Service struct {
	store      NotificationStore
	mailDialer *gomail.Dialer
}

// NewService tạo instance mới
func NewService(store NotificationStore, mailDialer *gomail.Dialer) *Service {
	return &Service{
		store:      store,
		mailDialer: mailDialer,
	}
}

// SendVerificationEmail gửi email xác minh
func (s *Service) SendVerificationEmail(ctx context.Context, email, token, baseURL string) (string, error) {
	// Tạo thông báo
	verifyURL := fmt.Sprintf("%s/verify?token=%s", baseURL, token)
	notification := &EmailNotification{
		Email:   email,
		Subject: "Verify Your Email",
		Body:    fmt.Sprintf("Please verify your email by clicking this link: %s", verifyURL),
	}

	// Lưu thông báo vào store
	err := s.store.SaveNotification(ctx, notification)
	if err != nil {
		return "", fmt.Errorf("failed to save notification: %v", err)
	}

	// Gửi email
	m := gomail.NewMessage()
	m.SetHeader("From", "votrungquan2002@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", notification.Subject)
	m.SetBody("text/plain", notification.Body)

	err = s.mailDialer.DialAndSend(m)
	if err != nil {
		// Cập nhật trạng thái thất bại
		if updateErr := s.store.UpdateNotificationStatus(ctx, notification.ID, "failed"); updateErr != nil {
			return "", fmt.Errorf("failed to send email: %v, and failed to update status: %v", err, updateErr)
		}
		return "", fmt.Errorf("failed to send email: %v", err)
	}

	// Cập nhật trạng thái thành công
	err = s.store.UpdateNotificationStatus(ctx, notification.ID, "sent")
	if err != nil {
		return "", fmt.Errorf("failed to update status to sent: %v", err)
	}

	return "Email sent", nil
}
