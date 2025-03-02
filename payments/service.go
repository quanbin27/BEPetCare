package main

import (
	"context"
	"errors"
	pb "github.com/quanbin27/commons/genproto/payments"
)

type paymentService struct {
	store PaymentStore
}

// NewPaymentService - Khởi tạo service
func NewPaymentService(store PaymentStore) PaymentService {
	return &paymentService{store: store}
}

// CreatePayment - Tạo thanh toán mới
func (s *paymentService) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	payment := &Payment{
		OrderID:       req.OrderId,
		AppointmentID: req.AppointmentId,
		Amount:        req.Amount,
		Description:   req.Description,
		Status:        PaymentStatusPending,
		Method:        fromPbPaymentMethod(req.Method),
	}

	paymentId, err := s.store.CreatePayment(ctx, payment)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		PaymentId: paymentId,
	}, nil
}

// GetPaymentInfo - Lấy thông tin thanh toán
func (s *paymentService) GetPaymentInfo(ctx context.Context, req *pb.GetPaymentInfoRequest) (*pb.GetPaymentInfoResponse, error) {
	payment, err := s.store.GetPaymentByID(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.New("payment not found")
	}

	return &pb.GetPaymentInfoResponse{
		Status:        toPbPaymentStatus(payment.Status),
		Method:        toPbPaymentMethod(payment.Method),
		OrderId:       payment.OrderID,
		AppointmentId: payment.AppointmentID,
		Amount:        payment.Amount,
		Description:   payment.Description,
	}, nil
}

// CreatePaymentURL - Tạo URL thanh toán (tích hợp PayOS)
func (s *paymentService) CreatePaymentURL(ctx context.Context, req *pb.CreatePaymentURLRequest) (*pb.CreatePaymentURLResponse, error) {
	// Gửi request đến PayOS API để tạo link thanh toán
	paymentLinkID := "test"
	checkoutURL := "https://payos.vn/checkout/" + paymentLinkID

	return &pb.CreatePaymentURLResponse{
		PaymentLinkId: paymentLinkID,
		CheckoutUrl:   checkoutURL,
	}, nil
}

// CancelPaymentLink - Hủy link thanh toán
func (s *paymentService) CancelPaymentLink(ctx context.Context, req *pb.CancelPaymentLinkRequest) (*pb.CancelPaymentLinkResponse, error) {
	// Giả lập hủy thanh toán trên PayOS
	return &pb.CancelPaymentLinkResponse{Status: "Cancelled"}, nil
}

// UpdatePaymentStatus - Cập nhật trạng thái thanh toán
func (s *paymentService) UpdatePaymentStatus(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.UpdatePaymentStatusResponse, error) {
	err := s.store.UpdatePaymentStatus(ctx, req.PaymentId, fromPbPaymentStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePaymentStatusResponse{Status: "Updated"}, nil
}

// UpdatePaymentMethod - Cập nhật phương thức thanh toán
func (s *paymentService) UpdatePaymentMethod(ctx context.Context, req *pb.UpdatePaymentMethodRequest) (*pb.UpdatePaymentMethodResponse, error) {
	err := s.store.UpdatePaymentMethod(ctx, req.PaymentId, fromPbPaymentMethod(req.Method))
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePaymentMethodResponse{Status: "Updated"}, nil
}

// UpdatePaymentAmount - Cập nhật số tiền thanh toán
func (s *paymentService) UpdatePaymentAmount(ctx context.Context, req *pb.UpdatePaymentAmountRequest) (*pb.UpdatePaymentAmountResponse, error) {
	err := s.store.UpdatePaymentAmount(ctx, req.PaymentId, req.Amount)
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePaymentAmountResponse{Status: "Updated"}, nil
}

// toPbPaymentStatus chuyển đổi từ PaymentStatus sang pb.PaymentStatus
func toPbPaymentStatus(status PaymentStatus) pb.PaymentStatus {
	switch status {
	case PaymentStatusPending:
		return pb.PaymentStatus_PENDING
	case PaymentStatusCompleted:
		return pb.PaymentStatus_COMPLETED
	case PaymentStatusFailed:
		return pb.PaymentStatus_FAILED
	case PaymentStatusCancelled:
		return pb.PaymentStatus_CANCELLED
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

// fromPbPaymentStatus chuyển đổi từ pb.PaymentStatus sang PaymentStatus
func fromPbPaymentStatus(pbStatus pb.PaymentStatus) PaymentStatus {
	switch pbStatus {
	case pb.PaymentStatus_PENDING:
		return PaymentStatusPending
	case pb.PaymentStatus_COMPLETED:
		return PaymentStatusCompleted
	case pb.PaymentStatus_FAILED:
		return PaymentStatusFailed
	case pb.PaymentStatus_CANCELLED:
		return PaymentStatusCancelled
	default:
		return "unknown"
	}
}
func fromPbPaymentMethod(pbMethod pb.PaymentMethod) PaymentMethod {
	switch pbMethod {
	case pb.PaymentMethod_BANK:
		return PaymentMethodBank
	case pb.PaymentMethod_CASH:
		return PaymentMethodCash
	default:
		return "unknown"
	}
}
func toPbPaymentMethod(paymentMethod PaymentMethod) pb.PaymentMethod {
	switch paymentMethod {
	case PaymentMethodBank:
		return pb.PaymentMethod_BANK
	case PaymentMethodCash:
		return pb.PaymentMethod_CASH
	default:
		return pb.PaymentMethod_CASH
	}
}
