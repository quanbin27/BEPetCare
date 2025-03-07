package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/payments"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	client pb.PaymentServiceClient
}

func NewPaymentHandler(client pb.PaymentServiceClient) *PaymentHandler {
	return &PaymentHandler{client: client}
}

// RegisterRoutes đăng ký các route cho Payments service với tiền tố "/payments"
func (h *PaymentHandler) RegisterRoutes(e *echo.Group) {
	e.POST("/payments", h.CreatePayment)
	e.GET("/payments/:payment_id", h.GetPaymentInfo)
	e.POST("/payments/url", h.CreatePaymentURL)
	e.POST("/payments/cancel", h.CancelPaymentLink)
	e.PUT("/payments/update-status", h.UpdatePaymentStatus)
	e.PUT("/payments/update-method", h.UpdatePaymentMethod)
	e.PUT("/payments/update-amount", h.UpdatePaymentAmount)
}

// CreatePayment xử lý yêu cầu tạo thanh toán
func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	var req struct {
		OrderID       int32   `json:"order_id"`
		AppointmentID int32   `json:"appointment_id"`
		Amount        float32 `json:"amount"`
		Description   string  `json:"description"`
		Method        string  `json:"method"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Kiểm tra và chuyển đổi PaymentMethod từ string sang pb.PaymentMethod
	pbMethod, ok := pb.PaymentMethod_value[req.Method]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment method"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.CreatePayment(ctx, &pb.CreatePaymentRequest{
		OrderId:       req.OrderID,
		AppointmentId: req.AppointmentID,
		Amount:        req.Amount,
		Description:   req.Description,
		Method:        pb.PaymentMethod(pbMethod),
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int32{
		"payment_id": resp.PaymentId,
	})
}

// GetPaymentInfo xử lý yêu cầu lấy thông tin thanh toán
func (h *PaymentHandler) GetPaymentInfo(c echo.Context) error {
	paymentIDStr := c.Param("payment_id")
	if paymentIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Payment ID is required"})
	}

	// Chuyển đổi payment_id từ string sang int32
	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.GetPaymentInfo(ctx, &pb.GetPaymentInfoRequest{PaymentId: int32(paymentID)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// CreatePaymentURL xử lý yêu cầu tạo URL thanh toán
func (h *PaymentHandler) CreatePaymentURL(c echo.Context) error {
	var req struct {
		PaymentID   string  `json:"payment_id"`
		Amount      float32 `json:"amount"`
		Description string  `json:"description"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi payment_id từ string sang int32
	paymentID, err := strconv.ParseInt(req.PaymentID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.CreatePaymentURL(ctx, &pb.CreatePaymentURLRequest{
		PaymentId:   int32(paymentID),
		Amount:      req.Amount,
		Description: req.Description,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"payment_link_id": resp.PaymentLinkId,
		"checkout_url":    resp.CheckoutUrl,
	})
}

// CancelPaymentLink xử lý yêu cầu hủy link thanh toán
func (h *PaymentHandler) CancelPaymentLink(c echo.Context) error {
	var req struct {
		PaymentID          string `json:"payment_id"`
		CancellationReason string `json:"cancellation_reason"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi payment_id từ string sang int32
	paymentID, err := strconv.ParseInt(req.PaymentID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.CancelPaymentLink(ctx, &pb.CancelPaymentLinkRequest{
		PaymentId:          int32(paymentID),
		CancellationReason: req.CancellationReason,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": resp.Status,
	})
}

// UpdatePaymentStatus xử lý yêu cầu cập nhật trạng thái thanh toán
func (h *PaymentHandler) UpdatePaymentStatus(c echo.Context) error {
	var req struct {
		PaymentID string `json:"payment_id"`
		Status    string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi payment_id từ string sang int32
	paymentID, err := strconv.ParseInt(req.PaymentID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment_id format, must be an integer"})
	}

	// Kiểm tra và chuyển đổi PaymentStatus từ string sang pb.PaymentStatus
	pbStatus, ok := pb.PaymentStatus_value[req.Status]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment status"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.UpdatePaymentStatus(ctx, &pb.UpdatePaymentStatusRequest{
		PaymentId: int32(paymentID),
		Status:    pb.PaymentStatus(pbStatus),
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": resp.Status,
	})
}

// UpdatePaymentMethod xử lý yêu cầu cập nhật phương thức thanh toán
func (h *PaymentHandler) UpdatePaymentMethod(c echo.Context) error {
	var req struct {
		PaymentID string `json:"payment_id"`
		Method    string `json:"method"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi payment_id từ string sang int32
	paymentID, err := strconv.ParseInt(req.PaymentID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment_id format, must be an integer"})
	}

	// Kiểm tra và chuyển đổi PaymentMethod từ string sang pb.PaymentMethod
	pbMethod, ok := pb.PaymentMethod_value[req.Method]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment method"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.UpdatePaymentMethod(ctx, &pb.UpdatePaymentMethodRequest{
		PaymentId: int32(paymentID),
		Method:    pb.PaymentMethod(pbMethod),
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": resp.Status,
	})
}

// UpdatePaymentAmount xử lý yêu cầu cập nhật số tiền thanh toán
func (h *PaymentHandler) UpdatePaymentAmount(c echo.Context) error {
	var req struct {
		PaymentID string  `json:"payment_id"`
		Amount    float32 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi payment_id từ string sang int32
	paymentID, err := strconv.ParseInt(req.PaymentID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payment_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.UpdatePaymentAmount(ctx, &pb.UpdatePaymentAmountRequest{
		PaymentId: int32(paymentID),
		Amount:    req.Amount,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": resp.Status,
	})
}
