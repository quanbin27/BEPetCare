package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/orders"
	pbProduct "github.com/quanbin27/commons/genproto/products"
	pbUser "github.com/quanbin27/commons/genproto/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	client        pb.OrderServiceClient
	userClient    pbUser.UserServiceClient
	productClient pbProduct.ProductServiceClient
}

func NewOrderHandler(client pb.OrderServiceClient, userClient pbUser.UserServiceClient, productClient pbProduct.ProductServiceClient) *OrderHandler {
	return &OrderHandler{client: client, userClient: userClient, productClient: productClient}
}

// RegisterRoutes đăng ký các route cho Orders service với tiền tố "/orders"
func (h *OrderHandler) RegisterRoutes(e *echo.Group) {
	e.POST("/orders", h.CreateOrder)
	e.GET("/orders/:order_id", h.GetOrder)
	e.GET("/orders/appointment/:appointment_id", h.GetOrderByAppointmentID)
	e.PUT("/orders/update-status", h.UpdateOrderStatus)
	e.GET("/orders/:order_id/items", h.GetOrderItems)
	e.GET("/orders/customer/:customer_id", h.GetOrderByCustomerID)
}

// CreateOrder xử lý yêu cầu tạo đơn hàng
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	fmt.Println("CREATE ORDERS")
	type OrderItemReq struct {
		ProductID   int32   `json:"product_id"`
		ProductType string  `json:"product_type"`
		Quantity    int32   `json:"quantity"`
		UnitPrice   float32 `json:"unit_price"`
		ProductName string  `json:"product_name"`
	}
	var req struct {
		CustomerID    int32          `json:"customer_id"`
		BranchID      int32          `json:"branch_id"`
		AppointmentID int32          `json:"appointment_id"`
		Items         []OrderItemReq `json:"items"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi items từ HTTP request sang gRPC request
	pbItems := make([]*pb.OrderItem, len(req.Items))
	for i, item := range req.Items {
		pbItems[i] = &pb.OrderItem{
			ProductId:   item.ProductID,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			ProductType: item.ProductType,
			ProductName: item.ProductName,
		}
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()
	user, err := h.userClient.GetUserInfo(ctx, &pbUser.GetUserInfoRequest{
		ID: req.CustomerID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	resp, err := h.client.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId:    req.CustomerID,
		BranchId:      req.BranchID,
		Items:         pbItems,
		AppointmentId: req.AppointmentID,
		Email:         user.Email,
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
		"order_id": resp.OrderId,
		"status":   resp.Status,
	})
}
func (h *OrderHandler) GetOrderByCustomerID(c echo.Context) error {
	customerIDStr := c.Param("customer_id")
	if customerIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Customer ID is required"})
	}

	customerID, err := strconv.ParseInt(customerIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid customer_id format, must be an integer"})
	}

	ctx := c.Request().Context()

	resp, err := h.client.GetOrdersByCustomerID(ctx, &pb.GetOrdersByCustomerIDRequest{CustomerId: int32(customerID)})
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

	return c.JSON(http.StatusOK, resp.Orders)
}

// GetOrder xử lý yêu cầu lấy thông tin đơn hàng
func (h *OrderHandler) GetOrder(c echo.Context) error {
	orderIDStr := c.Param("order_id")
	if orderIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Order ID is required"})
	}

	orderID, err := strconv.ParseInt(orderIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order_id format, must be an integer"})
	}

	ctx := c.Request().Context()

	resp, err := h.client.GetOrder(ctx, &pb.GetOrderRequest{OrderId: int32(orderID)})
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

	return c.JSON(http.StatusOK, resp.Order)
}
func (h *OrderHandler) GetOrderByAppointmentID(c echo.Context) error {
	appointmentIDStr := c.Param("appointment_id")
	if appointmentIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "appointment ID is required"})
	}

	appointmentID, err := strconv.ParseInt(appointmentIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid appointment_id format, must be an integer"})
	}

	ctx := c.Request().Context()

	resp, err := h.client.GetOrderByAppointmentID(ctx, &pb.GetOrderByAppointmentIDRequest{AppointmentId: int32(appointmentID)})
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
	return c.JSON(http.StatusOK, resp.Order)
}

// UpdateOrderStatus xử lý yêu cầu cập nhật trạng thái đơn hàng
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	var req struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi order_id từ string sang int32
	orderID, err := strconv.ParseInt(req.OrderID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order_id format, must be an integer"})
	}

	// Kiểm tra và chuyển đổi status sang pb.OrderStatus
	pbStatus, ok := pb.OrderStatus_value[req.Status]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid status value"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.UpdateOrderStatus(ctx, &pb.UpdateOrderStatusRequest{
		OrderId: int32(orderID),
		Status:  pb.OrderStatus(pbStatus),
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

// GetOrderItems xử lý yêu cầu lấy danh sách items của đơn hàng
func (h *OrderHandler) GetOrderItems(c echo.Context) error {
	orderIDStr := c.Param("order_id")
	if orderIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Order ID is required"})
	}

	// Chuyển đổi order_id từ string sang int32
	orderID, err := strconv.ParseInt(orderIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid order_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.GetOrderItems(ctx, &pb.GetOrderItemsRequest{OrderId: int32(orderID)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Items)
}
