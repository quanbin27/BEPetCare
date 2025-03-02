package main

import (
	"context"
	pb "github.com/quanbin27/commons/genproto/payments"
	"google.golang.org/grpc"
)

// PaymentGRPCHandler - Triển khai gRPC server cho PaymentService
type PaymentGRPCHandler struct {
	pb.UnimplementedPaymentServiceServer
	service PaymentService
}

// NewPaymentGRPCHandler - Khởi tạo handler
func NewGrpcPaymentHandler(grpcServer *grpc.Server, service PaymentService) {
	pb.RegisterPaymentServiceServer(grpcServer, &PaymentGRPCHandler{service: service})
}

// CreatePayment - Xử lý yêu cầu tạo thanh toán mới
func (h *PaymentGRPCHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	return h.service.CreatePayment(ctx, req)
}

// GetPaymentInfo - Xử lý yêu cầu lấy thông tin thanh toán
func (h *PaymentGRPCHandler) GetPaymentInfo(ctx context.Context, req *pb.GetPaymentInfoRequest) (*pb.GetPaymentInfoResponse, error) {
	return h.service.GetPaymentInfo(ctx, req)
}

// CreatePaymentURL - Xử lý yêu cầu tạo URL thanh toán PayOS
func (h *PaymentGRPCHandler) CreatePaymentURL(ctx context.Context, req *pb.CreatePaymentURLRequest) (*pb.CreatePaymentURLResponse, error) {
	return h.service.CreatePaymentURL(ctx, req)
}

// CancelPaymentLink - Xử lý yêu cầu hủy link thanh toán
func (h *PaymentGRPCHandler) CancelPaymentLink(ctx context.Context, req *pb.CancelPaymentLinkRequest) (*pb.CancelPaymentLinkResponse, error) {
	return h.service.CancelPaymentLink(ctx, req)
}

// UpdatePaymentStatus - Xử lý yêu cầu cập nhật trạng thái thanh toán
func (h *PaymentGRPCHandler) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.UpdatePaymentStatusResponse, error) {
	return h.service.UpdatePaymentStatus(ctx, req)
}

// UpdatePaymentMethod - Xử lý yêu cầu cập nhật phương thức thanh toán
func (h *PaymentGRPCHandler) UpdatePaymentMethod(ctx context.Context, req *pb.UpdatePaymentMethodRequest) (*pb.UpdatePaymentMethodResponse, error) {
	return h.service.UpdatePaymentMethod(ctx, req)
}

// UpdatePaymentAmount - Xử lý yêu cầu cập nhật số tiền thanh toán
func (h *PaymentGRPCHandler) UpdatePaymentAmount(ctx context.Context, req *pb.UpdatePaymentAmountRequest) (*pb.UpdatePaymentAmountResponse, error) {
	return h.service.UpdatePaymentAmount(ctx, req)
}
