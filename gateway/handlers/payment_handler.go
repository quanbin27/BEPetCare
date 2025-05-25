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

// CreatePayment creates a new payment
// @Summary Create a new payment
// @Description Creates a new payment record with details like order ID, appointment ID, amount, description, and payment method
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body object{order_id=integer,appointment_id=integer,amount=number,description=string,method=string} true "Payment details"
// @Success 200 {object} object{payment_id=integer} "Payment created successfully"
// @Failure 400 {object} object{error=string} "Invalid request or invalid payment method"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /payments [post]
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

// GetPaymentInfo retrieves payment information by ID
// @Summary Get payment by ID
// @Description Retrieves a payment record using its unique ID
// @Tags Payments
// @Produce json
// @Param payment_id path int true "Payment ID"
// @Success 200 {object} object{payment_id=integer,order_id=integer,appointment_id=integer,amount=number,description=string,method=string,status=string} "Payment details"
// @Failure 400 {object} object{error=string} "Payment ID is required or invalid format"
// @Failure 404 {object} object{error=string} "Payment not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /payments/{payment_id} [get]
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

// CreatePaymentURL creates a payment URL
// @Summary Create a payment URL
// @Description Creates a payment URL for a specific payment with details like payment ID, amount, and description
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body object{payment_id=string,amount=number,description=string} true "Payment URL details"
// @Success 200 {object} object{payment_link_id=string,checkout_url=string} "Payment URL created successfully"
// @Failure 400 {object} object{error=string} "Invalid request or invalid payment_id format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /payments/url [post]
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

// CancelPaymentLink cancels a payment link
// @Summary Cancel a payment link
// @Description Cancels a payment link for a specific payment ID with a cancellation reason
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body object{payment_id=string,cancellation_reason=string} true "Cancellation details"
// @Success 200 {object} object{status=string} "Payment link cancelled successfully"
// @Failure 400 {object} object{error=string} "Invalid request or invalid payment_id format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /payments/cancel [post]
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

// UpdatePaymentStatus updates the status of a payment
// @Summary Update payment status
// @Description Updates the status of a payment record for a specific payment ID
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body object{payment_id=string,status=string} true "Payment status update details"
// @Success 200 {object} object{status=string} "Payment status updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request, invalid payment_id format, or invalid payment status"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /payments/update-status [put]
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

// UpdatePaymentMethod updates the payment method
// @Summary Update payment method
// @Description Updates the payment method for a specific payment ID
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body object{payment_id=string,method=string} true "Payment method update details"
// @Success 200 {object} object{status=string} "Payment method updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request, invalid payment_id format, or invalid payment method"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /payments/update-method [put]
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

// UpdatePaymentAmount updates the payment amount
// @Summary Update payment amount
// @Description Updates the amount for a specific payment ID
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body object{payment_id=string,amount=number} true "Payment amount update details"
// @Success 200 {object} object{status=string} "Payment amount updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request or invalid payment_id format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /payments/update-amount [put]
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
