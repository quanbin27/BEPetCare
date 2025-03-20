package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"time"

	"github.com/quanbin27/commons/auth"
	"github.com/quanbin27/commons/config"
	pbNotification "github.com/quanbin27/commons/genproto/notifications"
)

type Service struct {
	userStore       UserStore
	redis           *redis.Client
	notificationSvc pbNotification.NotificationServiceClient
	baseURL         string
}

func NewService(store UserStore, redisClient *redis.Client, notificationConn *grpc.ClientConn, baseURL string) *Service {
	return &Service{
		userStore:       store,
		redis:           redisClient,
		notificationSvc: pbNotification.NewNotificationServiceClient(notificationConn),
		baseURL:         baseURL,
	}
}

// Register creates a new user
func (s *Service) Register(ctx context.Context, email, password, name string) (string, error) {
	_, err := s.userStore.GetUserByEmail(ctx, email)
	if err == nil {
		return "Failed", errors.New("user already exists")
	}
	if s.isPendingEmail(ctx, email) {
		return "", errors.New("email is already pending verification")
	}
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "Failed", errors.New("failed to hash password")
	}
	token, err := generateToken()
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	pu := &PendingUser{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
		Token:    token,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	err = s.savePendingUser(ctx, pu)
	if err != nil {
		return "", err
	}
	_, err = s.notificationSvc.SendVerificationEmail(ctx, &pbNotification.SendVerificationEmailRequest{
		Email:   email,
		Token:   token,
		BaseUrl: s.baseURL,
	})
	if err != nil {
		s.redis.Del(ctx, "pending:"+token)
		return "", err
	}
	return "Verification email sent", nil
}
func (s *Service) VerifyEmail(ctx context.Context, token string) (int32, error) {
	// Lấy thông tin từ Redis
	pu, err := s.getPendingUser(ctx, token)
	if err != nil {
		return 0, err
	}

	// Kiểm tra lại email trong DB (tránh race condition)
	if _, err := s.userStore.GetUserByEmail(ctx, pu.Email); err == nil {
		s.redis.Del(ctx, "pending:"+token)
		return 0, errors.New("user already exists")
	}
	// Tạo user để lưu vào DB
	user := &User{
		Email:    pu.Email,
		Password: pu.Password,
		Name:     pu.Name,
	}
	userId, err := s.userStore.CreateUser(ctx, user)
	if err != nil {
		return 0, errors.New("failed to create user: " + err.Error())
	}

	// Xóa dữ liệu tạm trong Redis
	s.redis.Del(ctx, "pending:"+token)
	err = s.userStore.CreateRole(ctx, userId, 3)
	if err != nil {
		return userId, err
	}
	return userId, nil
}

// Login authenticates a user and generates a JWT token
func (s *Service) Login(ctx context.Context, email, password string, rememberMe bool) (string, string, error) {
	u, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return "Failed", "", errors.New("not found, invalid email")
	}
	if !auth.CheckPassword(u.Password, []byte(password)) {
		return "Failed", "", errors.New("invalid password")
	}

	secret := []byte(config.Envs.JWTSecret)
	roleId, err := s.userStore.GetRole(ctx, u.ID)
	if err != nil {
		return "Failed", "", err
	}
	var token string
	if rememberMe {
		token, err = auth.CreateJWT(secret, u.ID, config.Envs.JWTExpirationInSeconds, roleId)
	} else {
		token, err = auth.CreateJWT(secret, u.ID, 3600, roleId)
	}
	if err != nil {
		return "Failed", "", errors.New("failed to create JWT")
	}
	return "Success", token, nil
}

// ChangeInfo updates user info
func (s *Service) ChangeInfo(ctx context.Context, userID int32, email, name, address, phoneNumber string) (string, string, string, string, string, error) {
	updatedData := make(map[string]interface{})
	if email != "" {
		updatedData["email"] = email
	}
	if name != "" {
		updatedData["name"] = name
	}
	if phoneNumber != "" {
		updatedData["phone_number"] = phoneNumber
	}
	if address != "" {
		updatedData["address"] = address
	}
	if len(updatedData) == 0 {
		return "Failed", "", "", "", "", errors.New("no data to update")
	}
	err := s.userStore.UpdateInfo(ctx, userID, updatedData)
	if err != nil {
		return "Failed", "", "", "", "", errors.New("failed to update user")
	}
	// Lấy thông tin user sau khi cập nhật để trả về
	user, err := s.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return "Failed", "", "", "", "", err
	}
	return "Success", user.Email, user.Name, user.Address, user.PhoneNumber, nil
}

// ChangePassword updates user password
func (s *Service) ChangePassword(ctx context.Context, userID int32, oldPassword, newPassword string) (string, error) {
	if newPassword == "" {
		return "Failed", errors.New("invalid password")
	}
	user, err := s.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return "Failed", errors.New("user not found")
	}
	if !auth.CheckPassword(user.Password, []byte(oldPassword)) {
		return "Failed", errors.New("invalid old password")
	}
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return "Failed", errors.New("failed to hash password")
	}
	err = s.userStore.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return "Failed", errors.New("failed to update user")
	}
	return "Success", nil
}

// GetUserInfo retrieves user info by ID
func (s *Service) GetUserInfo(ctx context.Context, id int32) (*User, error) {
	user, err := s.userStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserInfoByEmail retrieves user info by email
func (s *Service) GetUserInfoByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
func (s *Service) isPendingEmail(ctx context.Context, email string) bool {
	iter := s.redis.Keys(ctx, "pending:*").Val()
	for _, key := range iter {
		data, _ := s.redis.HGetAll(ctx, key).Result()
		if data["email"] == email {
			return true
		}
	}
	return false
}
func generateToken() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func (s *Service) savePendingUser(ctx context.Context, pu *PendingUser) error {
	key := "pending:" + pu.Token
	data := map[string]interface{}{
		"email":    pu.Email,
		"password": pu.Password,
		"name":     pu.Name,
		"token":    pu.Token,
		"expires":  pu.Expires.Format(time.RFC3339),
	}
	err := s.redis.HMSet(ctx, key, data).Err()
	if err != nil {
		return err
	}
	s.redis.Expire(ctx, key, 24*time.Hour)
	return nil
}
func (s *Service) getPendingUser(ctx context.Context, token string) (*PendingUser, error) {
	key := "pending:" + token
	data, err := s.redis.HGetAll(ctx, key).Result()
	if err != nil || len(data) == 0 {
		return nil, errors.New("invalid or expired token")
	}

	expires, err := time.Parse(time.RFC3339, data["expires"])
	if err != nil || time.Now().After(expires) {
		s.redis.Del(ctx, key)
		return nil, errors.New("token expired")
	}

	return &PendingUser{
		Email:    data["email"],
		Password: data["password"],
		Name:     data["name"],
		Token:    data["token"],
		Expires:  expires,
	}, nil
}
