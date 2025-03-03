package main

import (
	"context"

	pb "github.com/quanbin27/commons/genproto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersGrpcHandler struct {
	userService UserService
	pb.UnimplementedUserServiceServer
}

func NewGrpcUsersHandler(grpc *grpc.Server, userService UserService) {
	grpcHandler := &UsersGrpcHandler{
		userService: userService,
	}
	pb.RegisterUserServiceServer(grpc, grpcHandler)
}

func (h *UsersGrpcHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	stt, err := h.userService.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		// Nếu Service trả về lỗi, ánh xạ thành mã gRPC phù hợp
		if stt == "Failed" {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.RegisterResponse{Status: stt}, nil
}

func (h *UsersGrpcHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	stt, token, err := h.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		// Ánh xạ lỗi từ Service sang mã gRPC
		if stt == "Failed" {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LoginResponse{Status: stt, Token: token}, nil
}

func (h *UsersGrpcHandler) ChangeInfo(ctx context.Context, req *pb.ChangeInfoRequest) (*pb.ChangeInfoResponse, error) {
	stt, email, name, err := h.userService.ChangeInfo(ctx, req.Id, req.Email, req.Name)
	if err != nil {
		// Ánh xạ lỗi từ Service sang mã gRPC
		if stt == "Failed" {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ChangeInfoResponse{Status: stt, Email: email, Name: name}, nil
}

func (h *UsersGrpcHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	stt, err := h.userService.ChangePassword(ctx, req.Id, req.OldPassword, req.NewPassword)
	if err != nil {
		// Ánh xạ lỗi từ Service sang mã gRPC
		if stt == "Failed" {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.ChangePasswordResponse{Status: stt}, nil
}

func (h *UsersGrpcHandler) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.User, error) {
	user, err := h.userService.GetUserInfo(ctx, req.ID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return toProtoUser(user), nil
}

func (h *UsersGrpcHandler) GetUserInfoByEmail(ctx context.Context, req *pb.GetUserInfoByEmailRequest) (*pb.User, error) {
	user, err := h.userService.GetUserInfoByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return toProtoUser(user), nil
}
