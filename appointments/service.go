package main

import (
	"context"
	"time"

	pb "github.com/quanbin27/commons/genproto/appointments"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AppService triển khai các chức năng quản lý lịch hẹn & dịch vụ
type AppService struct {
	store AppointmentStore
}

// NewAppService khởi tạo service
func NewAppointmentService(store AppointmentStore) AppointmentService {
	return &AppService{store: store}
}

// --- LỊCH HẸN ---
// Tạo lịch hẹn
func (s *AppService) CreateAppointment(ctx context.Context, req *pb.CreateAppointmentRequest) (*pb.CreateAppointmentResponse, error) {
	appointment := &Appointment{
		CustomerID:      req.CustomerId,
		EmployeeID:      req.EmployeeId,
		ScheduledTime:   req.ScheduledTime.AsTime(),
		CustomerAddress: req.CustomerAddress,
		Status:          StatusPending,
		CreatedAt:       time.Now(),
	}

	// Lưu vào DB
	if err := s.store.CreateAppointment(ctx, appointment, req.ServiceIds); err != nil {
		return nil, err
	}

	return &pb.CreateAppointmentResponse{AppointmentId: appointment.ID, Status: "Success"}, nil
}

// Lấy lịch hẹn theo khách hàng
func (s *AppService) GetAppointmentsByCustomer(ctx context.Context, req *pb.GetAppointmentsByCustomerRequest) (*pb.GetAppointmentsResponse, error) {
	appointments, err := s.store.GetAppointmentsByCustomer(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}

	var pbAppointments []*pb.Appointment
	for _, a := range appointments {
		pbAppointments = append(pbAppointments, &pb.Appointment{
			Id:              a.ID,
			CustomerId:      a.CustomerID,
			EmployeeId:      a.EmployeeID,
			ScheduledTime:   timestamppb.New(a.ScheduledTime),
			Status:          toPbAppointmentStatus(a.Status),
			CustomerAddress: a.CustomerAddress,
		})
	}

	return &pb.GetAppointmentsResponse{Appointments: pbAppointments}, nil
}

// Lấy lịch hẹn theo nhân viên
func (s *AppService) GetAppointmentsByEmployee(ctx context.Context, req *pb.GetAppointmentsByEmployeeRequest) (*pb.GetAppointmentsResponse, error) {
	appointments, err := s.store.GetAppointmentsByEmployee(ctx, req.EmployeeId)
	if err != nil {
		return nil, err
	}

	var pbAppointments []*pb.Appointment
	for _, a := range appointments {
		pbAppointments = append(pbAppointments, &pb.Appointment{
			Id:              a.ID,
			CustomerId:      a.CustomerID,
			EmployeeId:      a.EmployeeID,
			ScheduledTime:   timestamppb.New(a.ScheduledTime),
			Status:          toPbAppointmentStatus(a.Status),
			CustomerAddress: a.CustomerAddress,
		})
	}

	return &pb.GetAppointmentsResponse{Appointments: pbAppointments}, nil
}

// Cập nhật trạng thái lịch hẹn
func (s *AppService) UpdateAppointmentStatus(ctx context.Context, req *pb.UpdateAppointmentStatusRequest) (*pb.UpdateAppointmentStatusResponse, error) {
	newStatus := fromPbAppointmentStatus(req.Status)
	if err := s.store.UpdateAppointmentStatus(ctx, req.AppointmentId, newStatus); err != nil {
		return nil, err
	}
	return &pb.UpdateAppointmentStatusResponse{Status: "Success"}, nil
}

// Lấy chi tiết lịch hẹn
func (s *AppService) GetAppointmentDetails(ctx context.Context, req *pb.GetAppointmentDetailsRequest) (*pb.GetAppointmentDetailsResponse, error) {
	appointment, details, err := s.store.GetAppointmentDetails(ctx, req.AppointmentId)
	if err != nil {
		return nil, err
	}

	var pbDetails []*pb.AppointmentDetail
	for _, d := range details {
		pbDetails = append(pbDetails, &pb.AppointmentDetail{
			AppointmentId: d.AppointmentID,
			ServiceId:     d.ServiceID,
			ServicePrice:  d.ServicePrice,
		})
	}

	return &pb.GetAppointmentDetailsResponse{
		Appointment: &pb.Appointment{
			Id:              appointment.ID,
			CustomerId:      appointment.CustomerID,
			EmployeeId:      appointment.EmployeeID,
			ScheduledTime:   timestamppb.New(appointment.ScheduledTime),
			Status:          toPbAppointmentStatus(appointment.Status),
			CustomerAddress: appointment.CustomerAddress,
		},
		Details: pbDetails,
	}, nil
}

// --- DỊCH VỤ ---
// Tạo dịch vụ
func (s *AppService) CreateService(ctx context.Context, req *pb.CreateServiceRequest) (*pb.CreateServiceResponse, error) {
	service := &Service{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CreatedAt:   time.Now(),
	}

	if err := s.store.CreateService(ctx, service); err != nil {
		return nil, err
	}

	return &pb.CreateServiceResponse{ServiceId: service.ID, Status: "Success"}, nil
}

// Lấy danh sách dịch vụ
func (s *AppService) GetServices(ctx context.Context, req *pb.GetServicesRequest) (*pb.GetServicesResponse, error) {
	services, err := s.store.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	var pbServices []*pb.Service
	for _, svc := range services {
		pbServices = append(pbServices, &pb.Service{
			Id:          svc.ID,
			Name:        svc.Name,
			Description: svc.Description,
			Price:       svc.Price,
		})
	}

	return &pb.GetServicesResponse{Services: pbServices}, nil
}

// Cập nhật dịch vụ
func (s *AppService) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest) (*pb.UpdateServiceResponse, error) {
	service := &Service{
		ID:          req.ServiceId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := s.store.UpdateService(ctx, service); err != nil {
		return nil, err
	}

	return &pb.UpdateServiceResponse{Status: "Success"}, nil
}

// Xóa dịch vụ
func (s *AppService) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest) (*pb.DeleteServiceResponse, error) {
	if err := s.store.DeleteService(ctx, req.ServiceId); err != nil {
		return nil, err
	}

	return &pb.DeleteServiceResponse{Status: "Success"}, nil
}
