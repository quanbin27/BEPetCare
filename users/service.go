package main

import (
	"context"
	"errors"

	"github.com/quanbin27/commons/auth"
	"github.com/quanbin27/commons/config"
)

type Service struct {
	userStore UserStore
}

func NewService(userStore UserStore) *Service {
	return &Service{userStore: userStore}
}

// Register creates a new user
func (s *Service) Register(ctx context.Context, email, password, name string) (string, error) {
	_, err := s.userStore.GetUserByEmail(ctx, email)
	if err == nil {
		return "Failed", errors.New("user already exists")
	}
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "Failed", errors.New("failed to hash password")
	}
	user := &User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
	}

	userId, err := s.userStore.CreateUser(ctx, user)
	if err != nil {
		return "Failed", err
	}
	err = s.userStore.CreateRole(ctx, userId, 3)
	if err != nil {
		return "Failed", err
	}
	return "Success", nil
}

// Login authenticates a user and generates a JWT token
func (s *Service) Login(ctx context.Context, email, password string) (string, string, error) {
	u, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return "Failed", "", errors.New("not found, invalid email")
	}
	if !auth.CheckPassword(u.Password, []byte(password)) {
		return "Failed", "", errors.New("invalid password")
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		return "Failed", "", errors.New("failed to create JWT")
	}
	return "Success", token, nil
}

// ChangeInfo updates user info
func (s *Service) ChangeInfo(ctx context.Context, userID int32, email, name string) (string, string, string, error) {
	updatedData := make(map[string]interface{})
	if email != "" {
		updatedData["email"] = email
	}
	if name != "" {
		updatedData["name"] = name
	}
	if len(updatedData) == 0 {
		return "Failed", "", "", errors.New("no data to update")
	}
	err := s.userStore.UpdateInfo(ctx, userID, updatedData)
	if err != nil {
		return "Failed", "", "", errors.New("failed to update user")
	}
	// Lấy thông tin user sau khi cập nhật để trả về
	user, err := s.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return "Failed", "", "", err
	}
	return "Success", user.Email, user.Name, nil
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
