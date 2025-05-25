package main

import (
	"context"
	"fmt"
	"log"

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
func (s *Service) SendResetPasswordEmail(ctx context.Context, email, token, baseURL string) (string, error) {
	// Tạo thông báo
	verifyURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)
	notification := &EmailNotification{
		Email:   email,
		Subject: "Reset Password",
		Body:    fmt.Sprintf("Please reset your password by clicking this link: %s", verifyURL),
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

// SendOrderConfirmationEmail gửi email xác nhận đơn hàng
func (s *Service) SendOrderConfirmationEmail(ctx context.Context, email string, orderID int32, items []OrderItem) (string, error) {
	body := fmt.Sprintf("Your order #%d has been placed successfully. Items:\n", orderID)
	for _, item := range items {
		log.Printf("ProductName: %q, Quantity: %d, UnitPrice: %f", item.ProductName, item.Quantity, item.UnitPrice)
		body += fmt.Sprintf("- %s (Qty: %d, Price: %.2f)\n", item.ProductName, item.Quantity, item.UnitPrice)
	}
	notification := &EmailNotification{
		Email:   email,
		Subject: fmt.Sprintf("Order #%d Confirmation", orderID),
		Body:    body,
	}

	err := s.store.SaveNotification(ctx, notification)
	if err != nil {
		return "", fmt.Errorf("failed to save notification: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "votrungquan2002@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", notification.Subject)
	m.SetBody("text/plain", notification.Body)

	err = s.mailDialer.DialAndSend(m)
	if err != nil {
		if updateErr := s.store.UpdateNotificationStatus(ctx, notification.ID, "failed"); updateErr != nil {
			return "", fmt.Errorf("failed to send email: %v, and failed to update status: %v", err, updateErr)
		}
		return "", fmt.Errorf("failed to send email: %v", err)
	}

	err = s.store.UpdateNotificationStatus(ctx, notification.ID, "sent")
	if err != nil {
		return "", fmt.Errorf("failed to update status to sent: %v", err)
	}

	return "Order confirmation email sent", nil
}
