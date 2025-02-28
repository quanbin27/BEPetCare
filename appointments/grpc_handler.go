package main

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/quanbin27/commons/genproto/appointments"
)

// AppointmentGrpcHandler triển khai gRPC server cho AppointmentService
type AppointmentGrpcHandler struct {
	pb.UnimplementedAppointmentServiceServer
	service AppointmentService
}

// NewGrpcAppointmentHandler đăng ký gRPC handler vào server
func NewGrpcAppointmentHandler(grpc *grpc.Server, service AppointmentService) {
	pb.RegisterAppointmentServiceServer(grpc, &AppointmentGrpcHandler{service: service})
}

// --- 🗓 LỊCH HẸN --- //

// Tạo lịch hẹn mới
func (h *AppointmentGrpcHandler) CreateAppointment(ctx context.Context, req *pb.CreateAppointmentRequest) (*pb.CreateAppointmentResponse, error) {
	return h.service.CreateAppointment(ctx, req)
}

// Lấy danh sách lịch hẹn theo khách hàng
func (h *AppointmentGrpcHandler) GetAppointmentsByCustomer(ctx context.Context, req *pb.GetAppointmentsByCustomerRequest) (*pb.GetAppointmentsResponse, error) {
	return h.service.GetAppointmentsByCustomer(ctx, req)
}

// Lấy danh sách lịch hẹn theo nhân viên
func (h *AppointmentGrpcHandler) GetAppointmentsByEmployee(ctx context.Context, req *pb.GetAppointmentsByEmployeeRequest) (*pb.GetAppointmentsResponse, error) {
	return h.service.GetAppointmentsByEmployee(ctx, req)
}

// Cập nhật trạng thái lịch hẹn
func (h *AppointmentGrpcHandler) UpdateAppointmentStatus(ctx context.Context, req *pb.UpdateAppointmentStatusRequest) (*pb.UpdateAppointmentStatusResponse, error) {
	return h.service.UpdateAppointmentStatus(ctx, req)
}

// Lấy chi tiết lịch hẹn
func (h *AppointmentGrpcHandler) GetAppointmentDetails(ctx context.Context, req *pb.GetAppointmentDetailsRequest) (*pb.GetAppointmentDetailsResponse, error) {
	return h.service.GetAppointmentDetails(ctx, req)
}

// --- 🛠 DỊCH VỤ --- //

// Tạo dịch vụ thú cưng mới
func (h *AppointmentGrpcHandler) CreateService(ctx context.Context, req *pb.CreateServiceRequest) (*pb.CreateServiceResponse, error) {
	return h.service.CreateService(ctx, req)
}

// Lấy danh sách dịch vụ
func (h *AppointmentGrpcHandler) GetServices(ctx context.Context, req *pb.GetServicesRequest) (*pb.GetServicesResponse, error) {
	return h.service.GetServices(ctx, req)
}

// Cập nhật thông tin dịch vụ
func (h *AppointmentGrpcHandler) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest) (*pb.UpdateServiceResponse, error) {
	return h.service.UpdateService(ctx, req)
}

// Xóa dịch vụ
func (h *AppointmentGrpcHandler) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest) (*pb.DeleteServiceResponse, error) {
	return h.service.DeleteService(ctx, req)
}
