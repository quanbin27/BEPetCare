package main

import (
	"context"

	pb "github.com/quanbin27/commons/genproto/payments"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentGrpcHandler struct {
	paymentService PaymentService
	pb.UnimplementedPaymentServiceServer
}

func NewPaymentGrpcHandler(grpc *grpc.Server, paymentService PaymentService) {
	grpcHandler := &PaymentGrpcHandler{
		paymentService: paymentService,
	}
	pb.RegisterPaymentServiceServer(grpc, grpcHandler)
}

func (h *PaymentGrpcHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	paymentID, err := h.paymentService.CreatePayment(ctx, req.OrderId, req.AppointmentId, req.Amount, req.Description, fromProtoPaymentMethod(req.Method))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreatePaymentResponse{PaymentId: paymentID}, nil
}

func (h *PaymentGrpcHandler) GetPaymentInfo(ctx context.Context, req *pb.GetPaymentInfoRequest) (*pb.GetPaymentInfoResponse, error) {
	payment, err := h.paymentService.GetPaymentInfo(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return toProtoPayment(payment), nil
}

func (h *PaymentGrpcHandler) CreatePaymentURL(ctx context.Context, req *pb.CreatePaymentURLRequest) (*pb.CreatePaymentURLResponse, error) {
	paymentLinkID, checkoutURL, err := h.paymentService.CreatePaymentURL(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreatePaymentURLResponse{
		PaymentLinkId: paymentLinkID,
		CheckoutUrl:   checkoutURL,
	}, nil
}

func (h *PaymentGrpcHandler) CancelPaymentLink(ctx context.Context, req *pb.CancelPaymentLinkRequest) (*pb.CancelPaymentLinkResponse, error) {
	statusMsg, err := h.paymentService.CancelPaymentLink(ctx, req.PaymentId, req.CancellationReason)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CancelPaymentLinkResponse{Status: statusMsg}, nil
}

func (h *PaymentGrpcHandler) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.UpdatePaymentStatusResponse, error) {
	statusMsg, err := h.paymentService.UpdatePaymentStatus(ctx, req.PaymentId, fromProtoPaymentStatus(req.Status))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdatePaymentStatusResponse{Status: statusMsg}, nil
}
func (h *PaymentGrpcHandler) UpdateBankPaymentStatus(ctx context.Context, req *pb.UpdateBankPaymentStatusRequest) (*pb.UpdatePaymentStatusResponse, error) {
	statusMsg, err := h.paymentService.UpdateBankPaymentStatus(ctx, req.OrderCode, fromProtoPaymentStatus(req.Status))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdatePaymentStatusResponse{Status: statusMsg}, nil
}

func (h *PaymentGrpcHandler) UpdatePaymentMethod(ctx context.Context, req *pb.UpdatePaymentMethodRequest) (*pb.UpdatePaymentMethodResponse, error) {
	statusMsg, err := h.paymentService.UpdatePaymentMethod(ctx, req.PaymentId, fromProtoPaymentMethod(req.Method))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdatePaymentMethodResponse{Status: statusMsg}, nil
}

func (h *PaymentGrpcHandler) UpdatePaymentAmount(ctx context.Context, req *pb.UpdatePaymentAmountRequest) (*pb.UpdatePaymentAmountResponse, error) {
	statusMsg, err := h.paymentService.UpdatePaymentAmount(ctx, req.PaymentId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &pb.UpdatePaymentAmountResponse{Status: statusMsg}, nil
}
