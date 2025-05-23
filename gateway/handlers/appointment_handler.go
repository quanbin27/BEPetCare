package handlers

import (
	"github.com/quanbin27/commons/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/appointments"
	pbOrder "github.com/quanbin27/commons/genproto/orders"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AppointmentHandler struct {
	client      pb.AppointmentServiceClient
	orderClient pbOrder.OrderServiceClient
}

func NewAppointmentHandler(client pb.AppointmentServiceClient, orderClient pbOrder.OrderServiceClient) *AppointmentHandler {
	return &AppointmentHandler{client: client, orderClient: orderClient}
}

// RegisterRoutes đăng ký các route cho Appointments service
func (h *AppointmentHandler) RegisterRoutes(e *echo.Group) {
	// Routes cho lịch hẹn
	e.POST("/appointments", h.CreateAppointment, auth.WithJWTAuth())
	e.GET("/appointments/customer/:customer_id", h.GetAppointmentsByCustomer)
	e.GET("/appointments/employee/:employee_id", h.GetAppointmentsByEmployee)
	e.PUT("/appointments/update-status", h.UpdateAppointmentStatus)
	e.GET("/appointments/:appointment_id", h.GetAppointmentDetails)

	// Routes cho dịch vụ
	e.POST("/services", h.CreateService)
	e.GET("/services", h.GetServices)
	e.PUT("/services", h.UpdateService)
	e.DELETE("/services/:service_id", h.DeleteService)
}

// --- Lịch hẹn ---

// CreateAppointment xử lý yêu cầu tạo lịch hẹn
func (h *AppointmentHandler) CreateAppointment(c echo.Context) error {
	type ServiceItemReq struct {
		ServiceId int32 `json:"service_id"`
		Quantity  int32 `json:"quantity"`
	}
	var req struct {
		CustomerID      int32            `json:"customer_id"`
		CustomerAddress string           `json:"customer_address"`
		ScheduledTime   string           `json:"scheduled_time"` // RFC3339 format (e.g., "2025-03-07T10:00:00Z")
		Detail          []ServiceItemReq `json:"services"`
		Note            string           `json:"note"`
		BranchID        int32            `json:"branch_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi ScheduledTime từ string sang timestamppb
	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledTime)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid scheduled_time format, must be RFC3339"})
	}
	// Lấy context từ Echo request
	ctx := c.Request().Context()
	pbItems := make([]*pb.AppointmentDetail, len(req.Detail))
	for i, item := range req.Detail {
		pbItems[i] = &pb.AppointmentDetail{
			Quantity:  item.Quantity,
			ServiceId: item.ServiceId,
		}
	}
	resp, err := h.client.CreateAppointment(ctx, &pb.CreateAppointmentRequest{
		CustomerId:      req.CustomerID,
		CustomerAddress: req.CustomerAddress,
		ScheduledTime:   timestamppb.New(scheduledTime),
		Detail:          pbItems,
		Note:            req.Note,
		BranchId:        req.BranchID,
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
		"appointment_id": resp.AppointmentId,
		"status":         resp.Status,
	})
}

// GetAppointmentsByCustomer xử lý yêu cầu lấy danh sách lịch hẹn theo customer_id
func (h *AppointmentHandler) GetAppointmentsByCustomer(c echo.Context) error {
	customerIDStr := c.Param("customer_id")
	if customerIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Customer ID is required"})
	}

	// Chuyển đổi customer_id từ string sang int32
	customerID, err := strconv.ParseInt(customerIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid customer_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.GetAppointmentsByCustomer(ctx, &pb.GetAppointmentsByCustomerRequest{CustomerId: int32(customerID)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	responses := make([]AppointmentResponse, 0, len(resp.Appointments))
	for _, pbApp := range resp.Appointments {
		responses = append(responses, toAppointmentResponse(pbApp))
	}

	return c.JSON(http.StatusOK, responses)

}

// GetAppointmentsByEmployee xử lý yêu cầu lấy danh sách lịch hẹn theo employee_id
func (h *AppointmentHandler) GetAppointmentsByEmployee(c echo.Context) error {
	employeeIDStr := c.Param("employee_id")
	if employeeIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Employee ID is required"})
	}

	// Chuyển đổi employee_id từ string sang int32
	employeeID, err := strconv.ParseInt(employeeIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid employee_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.GetAppointmentsByEmployee(ctx, &pb.GetAppointmentsByEmployeeRequest{EmployeeId: int32(employeeID)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	responses := make([]AppointmentResponse, 0, len(resp.Appointments))
	for _, pbApp := range resp.Appointments {
		responses = append(responses, toAppointmentResponse(pbApp))
	}

	return c.JSON(http.StatusOK, responses)

}

// UpdateAppointmentStatus xử lý yêu cầu cập nhật trạng thái lịch hẹn
func (h *AppointmentHandler) UpdateAppointmentStatus(c echo.Context) error {
	var req struct {
		AppointmentID string `json:"appointment_id"`
		Status        string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi appointment_id từ string sang int32
	appointmentID, err := strconv.ParseInt(req.AppointmentID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid appointment_id format, must be an integer"})
	}

	// Kiểm tra và chuyển đổi status sang pb.AppointmentStatus
	pbStatus, ok := pb.AppointmentStatus_value[req.Status]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid status value"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.UpdateAppointmentStatus(ctx, &pb.UpdateAppointmentStatusRequest{
		AppointmentId: int32(appointmentID),
		Status:        pb.AppointmentStatus(pbStatus),
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

// GetAppointmentDetails xử lý yêu cầu lấy chi tiết lịch hẹn
// GetAppointmentDetails xử lý yêu cầu lấy chi tiết lịch hẹn
func (h *AppointmentHandler) GetAppointmentDetails(c echo.Context) error {
	appointmentIDStr := c.Param("appointment_id")
	if appointmentIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Appointment ID is required"})
	}

	// Chuyển đổi appointment_id từ string sang int32
	appointmentID, err := strconv.ParseInt(appointmentIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid appointment_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	// Gọi AppointmentService
	resp, err := h.client.GetAppointmentDetails(ctx, &pb.GetAppointmentDetailsRequest{AppointmentId: int32(appointmentID)})
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

	// Convert Appointment
	appointment := toAppointmentResponse(resp.Appointment)

	// Gọi OrderService để tìm đơn hàng gắn với lịch hẹn
	orderResp, err := h.orderClient.GetOrderByAppointmentID(ctx, &pbOrder.GetOrderByAppointmentIDRequest{
		AppointmentId: int32(appointmentID),
	})

	var order interface{} = nil
	if err == nil && orderResp != nil && orderResp.Order != nil {
		order = orderResp.Order
	} else if grpcErr, ok := status.FromError(err); ok && grpcErr.Code() != codes.NotFound {
		// Nếu không phải lỗi "not found" thì trả về lỗi
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
	}

	// Trả về JSON gồm appointment, details và (nếu có) order
	return c.JSON(http.StatusOK, map[string]interface{}{
		"appointment": appointment,
		"details":     resp.Details,
		"order":       order, // null nếu không có
	})
}

// --- Dịch vụ ---

// CreateService xử lý yêu cầu tạo dịch vụ
func (h *AppointmentHandler) CreateService(c echo.Context) error {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.CreateService(ctx, &pb.CreateServiceRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
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
		"service_id": resp.ServiceId,
		"status":     resp.Status,
	})
}

// GetServices xử lý yêu cầu lấy danh sách dịch vụ
func (h *AppointmentHandler) GetServices(c echo.Context) error {
	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.GetServices(ctx, &pb.GetServicesRequest{})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": grpcErr.Message()})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Chuyển đổi dữ liệu từ resp.Services sang ServiceResponse
	var services []ServiceResponse
	for _, svc := range resp.Services {
		services = append(services, ServiceResponse{
			ServiceID:   svc.Id,
			Name:        svc.Name,
			Description: svc.Description,
			Price:       svc.Price,
			ImgURL:      svc.Imgurl,
		})
	}

	return c.JSON(http.StatusOK, services)
}

// UpdateService xử lý yêu cầu cập nhật dịch vụ
func (h *AppointmentHandler) UpdateService(c echo.Context) error {
	var req struct {
		ServiceID   string  `json:"service_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Chuyển đổi service_id từ string sang int32
	serviceID, err := strconv.ParseInt(req.ServiceID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid service_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.UpdateService(ctx, &pb.UpdateServiceRequest{
		ServiceId:   int32(serviceID),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
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

// DeleteService xử lý yêu cầu xóa dịch vụ
func (h *AppointmentHandler) DeleteService(c echo.Context) error {
	serviceIDStr := c.Param("service_id")
	if serviceIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Service ID is required"})
	}

	// Chuyển đổi service_id từ string sang int32
	serviceID, err := strconv.ParseInt(serviceIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid service_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.DeleteService(ctx, &pb.DeleteServiceRequest{ServiceId: int32(serviceID)})
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
