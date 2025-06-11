package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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
	updatedData["phone_number"] = phoneNumber
	updatedData["address"] = address
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
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, user.ID, 3600, 3)
	if err != nil {
		return errors.New("failed to generate token")
	}
	err = s.storeResetToken(ctx, user.ID, token, 24*time.Hour)
	if err != nil {
		return errors.New("failed to store ResetToken in Redis")
	}
	_, err = s.notificationSvc.SendResetPasswordEmail(ctx, &pbNotification.SendVerificationEmailRequest{
		Email:   email,
		Token:   token,
		BaseUrl: s.baseURL,
	})
	if err != nil {
		return errors.New("failed to send verification email")
	}
	return nil
}
func (s *Service) storeResetToken(ctx context.Context, userID int32, token string, ttl time.Duration) error {
	key := fmt.Sprintf("reset_token:%d", userID)
	return s.redis.Set(ctx, key, token, ttl).Err()
}
func (s *Service) VerifyResetToken(ctx context.Context, userID int32, token string) (bool, error) {
	key := fmt.Sprintf("reset_token:%d", userID)
	storedToken, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return storedToken == token, nil
}
func (s *Service) ResetPassword(ctx context.Context, userID int32, token, newPassword string) error {
	valid, err := s.VerifyResetToken(ctx, userID, token)
	if err != nil || !valid {
		return errors.New("invalid or expired token")
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	err = s.userStore.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return errors.New("failed to update password")
	}
	err = s.DeleteResetToken(ctx, userID)
	if err != nil {
		return errors.New("failed to delete reset token")
	}

	return nil
}

func (s *Service) DeleteResetToken(ctx context.Context, userID int32) error {
	key := fmt.Sprintf("reset_token:%d", userID)
	return s.redis.Del(ctx, key).Err()
}

// GetAllCustomers retrieves all users with role_id = 1
func (s *Service) GetAllCustomers(ctx context.Context) ([]User, error) {
	users, err := s.userStore.GetAllCustomers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}
	return users, nil
}

// GetAllUsers retrieves all users
func (s *Service) GetAllUsers(ctx context.Context) ([]UserWithRole, error) {
	users, err := s.userStore.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}
func (s *Service) EditUser(ctx context.Context, input UserWithRole) error {
	if err := s.userStore.UpdateUser(ctx, input); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *Service) GetBranchByEmployeeID(ctx context.Context, employeeID int32) (int32, error) {
	branchID, err := s.userStore.GetBranchByEmployeeID(ctx, employeeID)
	if err != nil {
		return 0, fmt.Errorf("failed to get branch for employee: %w", err)
	}
	return branchID, nil
}

// GetCustomersPaginated retrieves customers with pagination
func (s *Service) GetCustomersPaginated(ctx context.Context, page int32, pageSize int32) ([]User, int64, error) {
	if page < 1 || pageSize < 1 {
		return nil, 0, errors.New("invalid pagination parameters")
	}
	users, total, err := s.userStore.GetCustomersPaginated(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get paginated customers: %w", err)
	}
	return users, total, nil
}

// GetCustomersByName retrieves customers filtered by name
func (s *Service) GetCustomersByName(ctx context.Context, nameFilter string) ([]User, error) {
	if nameFilter == "" {
		return nil, errors.New("name filter cannot be empty")
	}
	users, err := s.userStore.GetCustomersByName(ctx, nameFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers by name: %w", err)
	}
	return users, nil
}
func (s *Service) CreateUser(ctx context.Context, email, name, phoneNumber string) (int32, error) {
	if _, err := s.userStore.GetUserByEmail(ctx, email); err == nil {
		return 0, errors.New("user already exists")
	}
	userID, err := s.userStore.CreateUser(ctx, &User{
		Email:       email,
		Name:        name,
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	err = s.userStore.CreateRole(ctx, userID, 1)
	if err != nil {
		return userID, err
	}

	return userID, nil
}
