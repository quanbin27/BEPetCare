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

func (s *Service) CreateUser(ctx context.Context, user *users.RegisterRequest) error {
	_, err := s.userStore.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("User already exists")
	}
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return errors.New("Failed to hash password")
	}
	return s.userStore.CreateUser(ctx, &User{Name: user.Name, Email: user.Email, Password: hashedPassword})
}
func (s *Service) Login(ctx context.Context, login *users.LoginRequest) (*users.LoginResponse, error) {
	u, err := s.userStore.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return &users.LoginResponse{Status: "not found, invalid email"}, errors.New("not found, invalid email")
	}
	if !auth.CheckPassword(u.Password, []byte(login.Password)) {
		return &users.LoginResponse{Status: "invalid passwor"}, errors.New("invalid password")
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		return &users.LoginResponse{Status: "Failed to create JWT"}, errors.New("Failed to create JWT")
	}
	return &users.LoginResponse{Status: "Success", Token: token}, nil
}
func (s *Service) UpdateUser(ctx context.Context, update *users.ChangeInfoRequest) error {
	updatedData := map[string]interface{}{
		"name":  update.Name,
		"email": update.Email,
	}
	err := s.userStore.UpdateInfo(ctx, update.Id, updatedData)
	if err != nil {
		return errors.New("Failed to update user")
	}
	return nil
}
func (s *Service) UpdatePassword(ctx context.Context, update *users.ChangePasswordRequest) error {
	if update.NewPassword == "" {
		return errors.New("Invalid password")
	}
	user, err := s.userStore.GetUserByID(ctx, update.Id)
	if err != nil {
		return errors.New("User not found")
	}
	if !auth.CheckPassword(user.Password, []byte(update.OldPassword)) {
		return errors.New("Invalid old password")
	}
	password, err := auth.HashPassword(update.NewPassword)
	err = s.userStore.UpdatePassword(ctx, user.ID, password)
	if err != nil {
		return errors.New("Failed to update user")
	}
	return nil
}
func (s *Service) GetUserByID(ctx context.Context, id int32) (*users.User, error) {
	user, err := s.userStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return &users.User{
		Name:      user.Name,
		Email:     user.Email,
		ID:        id,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	user, err := s.userStore.GetUserByEmail(ctx, email)
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
