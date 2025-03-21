package main

import (
	"context"
	"google.golang.org/grpc"

	pb "github.com/quanbin27/commons/genproto/notifications"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCHandler xử lý các yêu cầu gRPC
type GRPCHandler struct {
	pb.UnimplementedNotificationServiceServer
	service NotificationService
}

// NewGRPCHandler tạo instance mới
func NewGRPCHandler(grpc *grpc.Server, service NotificationService) {
	grpcHandler := &GRPCHandler{
		service: service,
	}
	pb.RegisterNotificationServiceServer(grpc, grpcHandler)
}

// SendVerificationEmail xử lý yêu cầu gửi email xác minh
func (h *GRPCHandler) SendVerificationEmail(ctx context.Context, req *pb.SendVerificationEmailRequest) (*pb.SendVerificationEmailResponse, error) {
	// Gọi service để gửi email
	statusMsg, err := h.service.SendVerificationEmail(ctx, req.Email, req.Token, req.BaseUrl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send verification email: %v", err)
	}
	return &pb.SendVerificationEmailResponse{Status: statusMsg}, nil
}
func (h *GRPCHandler) SendResetPasswordEmail(ctx context.Context, req *pb.SendVerificationEmailRequest) (*pb.SendVerificationEmailResponse, error) {
	statusMsg, err := h.service.SendResetPasswordEmail(ctx, req.Email, req.Token, req.BaseUrl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send reset password email: %v", err)
	}
	return &pb.SendVerificationEmailResponse{Status: statusMsg}, nil
}
