package main

import (
	"context"
	"github.com/quanbin27/commons/genproto/users"
	"google.golang.org/grpc"
)

type UsersGrpcHandler struct {
	userService UserService
	users.UnimplementedUserServiceServer
}

func NewGrpcUsersHandler(grpc *grpc.Server, userService UserService) {
	grpcHandler := &UsersGrpcHandler{
		userService: userService,
	}
	users.RegisterUserServiceServer(grpc, grpcHandler)
}
func (h *UsersGrpcHandler) Register(ctx context.Context, req *users.RegisterRequest) (*users.RegisterResponse, error) {
	err := h.userService.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &users.RegisterResponse{
		Status: "success",
	}
	return res, nil
}
func (h *UsersGrpcHandler) Login(ctx context.Context, req *users.LoginRequest) (*users.LoginResponse, error) {
	res, err := h.userService.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (h *UsersGrpcHandler) ChangeInfo(ctx context.Context, req *users.ChangeInfoRequest) (*users.ChangeInfoResponse, error) {
	err := h.userService.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &users.ChangeInfoResponse{
		Status: "success",
		Name:   req.Name,
		Email:  req.Email,
	}
	return res, nil
}
func (h *UsersGrpcHandler) ChangePassword(ctx context.Context, req *users.ChangePasswordRequest) (*users.ChangePasswordResponse, error) {
	err := h.userService.UpdatePassword(ctx, req)
	if err != nil {
		return nil, err
	}
	res := &users.ChangePasswordResponse{
		Status: "success",
	}
	return res, nil
}
func (h *UsersGrpcHandler) GetUserInfo(ctx context.Context, req *users.GetUserInfoRequest) (*users.User, error) {
	res, err := h.userService.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (h *UsersGrpcHandler) GetUserInfoByEmail(ctx context.Context, req *users.GetUserInfoByEmailRequest) (*users.User, error) {
	res, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	return res, nil
}
