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
	e.GET("/appointments/branch/:branch_id", h.GetAppointmentsByBranch)
	e.GET("/appointments/employee/:employee_id", h.GetAppointmentsByEmployee)
	e.PUT("/appointments/update-status", h.UpdateAppointmentStatus)
	e.PUT("/appointments/update-employee", h.UpdateEmployeeForAppointment)
	e.GET("/appointments/:appointment_id", h.GetAppointmentDetails)

	// Routes cho dịch vụ
	e.POST("/services", h.CreateService)
	e.GET("/services", h.GetServices)
	e.PUT("/services", h.UpdateService)
	e.DELETE("/services/:service_id", h.DeleteService)
}

// --- Lịch hẹn ---

// CreateAppointment xử lý yêu cầu tạo lịch hẹn
// CreateAppointment creates a new appointment
// @Summary Create a new appointment
// @Description Creates a new appointment with details like customer ID, address, scheduled time, services, note, and branch ID
// @Tags Appointments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{customer_id=integer,customer_address=string,scheduled_time=string,services=array{service_id=integer,quantity=integer},note=string,branch_id=integer} true "Appointment details"
// @Success 200 {object} object{appointment_id=integer,status=string} "Appointment created successfully"
// @Failure 400 {object} object{error=string} "Invalid request or invalid scheduled_time format (must be RFC3339)"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /appointments [post]
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

// GetAppointmentsByCustomer retrieves appointments by customer ID
// @Summary Get appointments by customer
// @Description Retrieves a list of appointment records for a specific customer ID
// @Tags Appointments
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Success 200 {array} object{id=integer,customer_id=integer,employee_id=integer,branch_id=integer,scheduled_time=string,status=string,note=string,customer_address=string} "List of appointments"
// @Failure 400 {object} object{error=string} "Customer ID is required or invalid format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /appointments/customer/{customer_id} [get]
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

// GetAppointmentsByEmployee retrieves appointments by employee ID
// @Summary Get appointments by employee
// @Description Retrieves a list of appointment records for a specific employee ID
// @Tags Appointments
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Success 200 {array} object{id=integer,customer_id=integer,employee_id=integer,branch_id=integer,scheduled_time=string,status=string,note=string,customer_address=string} "List of appointments"
// @Failure 400 {object} object{error=string} "Employee ID is required or invalid format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /appointments/employee/{employee_id} [get]
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

// GetAppointmentsByBranch retrieves appointments by branch ID
// @Summary Get appointments by branch
// @Description Retrieves a list of appointment records for a specific branch ID
// @Tags Appointments
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Success 200 {array} object{id=integer,customer_id=integer,employee_id=integer,branch_id=integer,scheduled_time=string,status=string,note=string,customer_address=string} "List of appointments"
// @Failure 400 {object} object{error=string} "Branch ID is required or invalid format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /appointments/branch/{branch_id} [get]
func (h *AppointmentHandler) GetAppointmentsByBranch(c echo.Context) error {
	branchIDStr := c.Param("branch_id")
	if branchIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "branchI ID is required"})
	}

	// Chuyển đổi branchIDStrtừ string sang int32
	branchID, err := strconv.ParseInt(branchIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid branch_id format, must be an integer"})
	}

	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.GetAppointmentsByBranch(ctx, &pb.GetAppointmentsByBranchRequest{BranchId: int32(branchID)})
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

// UpdateAppointmentStatus updates the status of an appointment
// @Summary Update appointment status
// @Description Updates the status of an appointment for a specific appointment ID
// @Tags Appointments
// @Accept json
// @Produce json
// @Param request body object{appointment_id=string,status=string} true "Appointment status update details"
// @Success 200 {object} object{status=string} "Appointment status updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request, invalid appointment_id format, or invalid status"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /appointments/update-status [put]
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

// UpdateEmployeeForAppointment updates the employee assigned to an appointment
// @Summary Update employee for appointment
// @Description Updates the employee assigned to a specific appointment ID
// @Tags Appointments
// @Accept json
// @Produce json
// @Param request body object{appointment_id=string,employee_id=string} true "Appointment employee update details"
// @Success 200 {object} object{status=string} "Employee updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request, invalid appointment_id or employee_id format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /appointments/update-employee [put]
func (h *AppointmentHandler) UpdateEmployeeForAppointment(c echo.Context) error {
	var req struct {
		AppointmentID string `json:"appointment_id"`
		EmployeeID    string `json:"employee_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	appointmentID, err := strconv.ParseInt(req.AppointmentID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid appointment_id format, must be an integer"})
	}
	employeeID, err := strconv.ParseInt(req.EmployeeID, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid employee_id format, must be an integer"})
	}
	// Lấy context từ Echo request
	ctx := c.Request().Context()

	resp, err := h.client.UpdateEmployeeForAppointment(ctx, &pb.UpdateEmployeeForAppointmentRequest{
		AppointmentId: int32(appointmentID),
		EmployeeId:    int32(employeeID),
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

// GetAppointmentDetails retrieves appointment details by ID
// @Summary Get appointment details
// @Description Retrieves detailed information for a specific appointment ID, including associated order if available
// @Tags Appointments
// @Produce json
// @Param appointment_id path int true "Appointment ID"
// @Success 200 {object} object{appointment=object{id=integer,customer_id=integer,employee_id=integer,branch_id=integer,scheduled_time=string,status=string,note=string,customer_address=string},details=array{service_id=integer,quantity=integer},order=object} "Appointment details"
// @Failure 400 {object} object{error=string} "Appointment ID is required or invalid format"
// @Failure 404 {object} object{error=string} "Appointment not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /appointments/{appointment_id} [get]
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

// CreateService creates a new service
// @Summary Create a new service
// @Description Creates a new service with details like name, description, and price
// @Tags Services
// @Accept json
// @Produce json
// @Param request body object{name=string,description=string,price=number} true "Service details"
// @Success 200 {object} object{service_id=integer,status=string} "Service created successfully"
// @Failure 400 {object} object{error=string} "Invalid request"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /services [post]
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

// GetServices retrieves all services
// @Summary List all services
// @Description Retrieves a list of all service records
// @Tags Services
// @Produce json
// @Success 200 {array} object{id=integer,name=string,description=string,price=number,imgurl=string} "List of services"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /services [get]
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

// UpdateService updates an existing service
// @Summary Update a service
// @Description Updates a service with details like service ID, name, description, and price
// @Tags Services
// @Accept json
// @Produce json
// @Param request body object{service_id=string,name=string,description=string,price=number} true "Updated service details"
// @Success 200 {object} object{status=string} "Service updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request or invalid service_id format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /services [put]
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

// DeleteService deletes a service
// @Summary Delete a service
// @Description Deletes a service by its unique ID
// @Tags Services
// @Produce json
// @Param service_id path int true "Service ID"
// @Success 200 {object} object{status=string} "Service deleted successfully"
// @Failure 400 {object} object{error=string} "Service ID is required or invalid format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /services/{service_id} [delete]
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
