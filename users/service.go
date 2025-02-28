package main

import (
	"context"
	"errors"
	"github.com/quanbin27/commons/auth"
	"github.com/quanbin27/commons/config"
	"github.com/quanbin27/commons/genproto/users"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	userStore UserStore
}

func NewService(userStore UserStore) *Service {
	return &Service{userStore: userStore}
}

func (s *Service) Register(ctx context.Context, user *users.RegisterRequest) (*users.RegisterResponse, error) {
	_, err := s.userStore.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return &users.RegisterResponse{Status: "Failed"}, errors.New("User already exists")
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return &users.RegisterResponse{Status: "Failed"}, errors.New("Failed to hash password")
	}
	return &users.RegisterResponse{Status: "Success"}, s.userStore.CreateUser(ctx, &User{Name: user.Name, Email: user.Email, Password: hashedPassword})
}
func (s *Service) Login(ctx context.Context, login *users.LoginRequest) (*users.LoginResponse, error) {
	u, err := s.userStore.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return &users.LoginResponse{Status: "Failed"}, errors.New("not found, invalid email")
	}
	if !auth.CheckPassword(u.Password, []byte(login.Password)) {
		return &users.LoginResponse{Status: "Failed"}, errors.New("invalid password")
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		return &users.LoginResponse{Status: "Failed"}, errors.New("Failed to create JWT")
	}
	return &users.LoginResponse{Status: "Success", Token: token}, nil
}
func (s *Service) ChangeInfo(ctx context.Context, update *users.ChangeInfoRequest) (*users.ChangeInfoResponse, error) {
	updatedData := map[string]interface{}{
		"name":  update.Name,
		"email": update.Email,
	}
	err := s.userStore.UpdateInfo(ctx, update.Id, updatedData)
	if err != nil {
		return &users.ChangeInfoResponse{Email: update.Email, Name: update.Name, Status: "Failed"}, errors.New("Failed to update user")
	}
	return &users.ChangeInfoResponse{Email: update.Email, Name: update.Name, Status: "Success"}, nil
}
func (s *Service) ChangePassword(ctx context.Context, update *users.ChangePasswordRequest) (*users.ChangePasswordResponse, error) {
	if update.NewPassword == "" {
		return &users.ChangePasswordResponse{Status: "Failed"}, errors.New("Invalid password")
	}
	user, err := s.userStore.GetUserByID(ctx, update.Id)
	if err != nil {
		return &users.ChangePasswordResponse{Status: "Failed"}, errors.New("User not found")
	}
	if !auth.CheckPassword(user.Password, []byte(update.OldPassword)) {
		return &users.ChangePasswordResponse{Status: "Failed"}, errors.New("Invalid old password")
	}
	password, err := auth.HashPassword(update.NewPassword)
	err = s.userStore.UpdatePassword(ctx, user.ID, password)
	if err != nil {
		return &users.ChangePasswordResponse{Status: "Failed"}, errors.New("Failed to update user")
	}
	return &users.ChangePasswordResponse{Status: "Success"}, nil
}
func (s *Service) GetUserInfo(ctx context.Context, id *users.GetUserInfoRequest) (*users.User, error) {
	user, err := s.userStore.GetUserByID(ctx, id.ID)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return &users.User{
		Name:      user.Name,
		Email:     user.Email,
		ID:        user.ID,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}
func (s *Service) GetUserInfoByEmail(ctx context.Context, email *users.GetUserInfoByEmailRequest) (*users.User, error) {
	user, err := s.userStore.GetUserByEmail(ctx, email.Email)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return &users.User{
		Name:      user.Name,
		Email:     user.Email,
		ID:        user.ID,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}
