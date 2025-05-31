package main

import (
	"context"
	"errors"

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
	statusMsg, err := h.userService.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, errors.New("user already exists")) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.RegisterResponse{Status: statusMsg}, nil
}
func (h *UsersGrpcHandler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	id, err := h.userService.VerifyEmail(ctx, req.Token)
	if err != nil {
		if errors.Is(err, errors.New("invalid or expired token")) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.VerifyEmailResponse{Id: id}, nil
}
func (h *UsersGrpcHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	stt, token, err := h.userService.Login(ctx, req.Email, req.Password, req.RememberMe)
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
	stt, email, name, address, phoneNumber, err := h.userService.ChangeInfo(ctx, req.Id, req.Email, req.Name, req.Address, req.PhoneNumber)
	if err != nil {
		// Ánh xạ lỗi từ Service sang mã gRPC
		if stt == "Failed" {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ChangeInfoResponse{Status: stt, Email: email, Name: name, PhoneNumber: phoneNumber, Address: address}, nil
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

func (h *UsersGrpcHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	err := h.userService.ForgotPassword(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ForgotPasswordResponse{Status: "sent reset password email"}, nil
}
func (h *UsersGrpcHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	err := h.userService.ResetPassword(ctx, req.UserID, req.Token, req.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ResetPasswordResponse{Status: "reset password success"}, nil
}
func (h *UsersGrpcHandler) GetAllCustomers(ctx context.Context, req *pb.GetAllCustomersRequest) (*pb.GetAllCustomersResponse, error) {
	users, err := h.userService.GetAllCustomers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(&user)
	}
	return &pb.GetAllCustomersResponse{Users: protoUsers}, nil
}

func (h *UsersGrpcHandler) GetCustomersPaginated(ctx context.Context, req *pb.GetCustomersPaginatedRequest) (*pb.GetCustomersPaginatedResponse, error) {
	users, total, err := h.userService.GetCustomersPaginated(ctx, req.Page, req.PageSize)
	if err != nil {
		if errors.Is(err, errors.New("invalid pagination parameters")) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(&user)
	}
	return &pb.GetCustomersPaginatedResponse{Users: protoUsers, Total: total}, nil
}

func (h *UsersGrpcHandler) GetCustomersByName(ctx context.Context, req *pb.GetCustomersByNameRequest) (*pb.GetCustomersByNameResponse, error) {
	users, err := h.userService.GetCustomersByName(ctx, req.NameFilter)
	if err != nil {
		if errors.Is(err, errors.New("name filter cannot be empty")) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(&user)
	}
	return &pb.GetCustomersByNameResponse{Users: protoUsers}, nil
}
