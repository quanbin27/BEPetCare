package main

import (
	"context"

	pb "github.com/quanbin27/commons/genproto/appointments"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppointmentGrpcHandler struct {
	appointmentService AppointmentService
	pb.UnimplementedAppointmentServiceServer
}

func NewAppointmentGrpcHandler(grpc *grpc.Server, appointmentService AppointmentService) {
	grpcHandler := &AppointmentGrpcHandler{
		appointmentService: appointmentService,
	}
	pb.RegisterAppointmentServiceServer(grpc, grpcHandler)
}

// --- LỊCH HẸN ---
func (h *AppointmentGrpcHandler) CreateAppointment(ctx context.Context, req *pb.CreateAppointmentRequest) (*pb.CreateAppointmentResponse, error) {
	Items := make([]AppointmentDetail, len(req.Detail))
	for i, item := range req.Detail {
		Items[i] = AppointmentDetail{
			Quantity:  item.Quantity,
			ServiceID: item.ServiceId,
		}
	}
	appointmentID, statusMsg, err := h.appointmentService.CreateAppointment(ctx, req.CustomerId, req.CustomerAddress, req.ScheduledTime.AsTime(), Items, req.Note, req.BranchId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateAppointmentResponse{AppointmentId: appointmentID, Status: statusMsg}, nil
}

func (h *AppointmentGrpcHandler) GetAppointmentsByCustomer(ctx context.Context, req *pb.GetAppointmentsByCustomerRequest) (*pb.GetAppointmentsResponse, error) {
	appointments, err := h.appointmentService.GetAppointmentsByCustomer(ctx, req.CustomerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pbAppointments := make([]*pb.Appointment, len(appointments))
	for i, a := range appointments {
		pbAppointments[i] = toProtoAppointment(&a)
	}
	return &pb.GetAppointmentsResponse{Appointments: pbAppointments}, nil
}

func (h *AppointmentGrpcHandler) GetAppointmentsByEmployee(ctx context.Context, req *pb.GetAppointmentsByEmployeeRequest) (*pb.GetAppointmentsResponse, error) {
	appointments, err := h.appointmentService.GetAppointmentsByEmployee(ctx, req.EmployeeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pbAppointments := make([]*pb.Appointment, len(appointments))
	for i, a := range appointments {
		pbAppointments[i] = toProtoAppointment(&a)
	}
	return &pb.GetAppointmentsResponse{Appointments: pbAppointments}, nil
}

func (h *AppointmentGrpcHandler) UpdateAppointmentStatus(ctx context.Context, req *pb.UpdateAppointmentStatusRequest) (*pb.UpdateAppointmentStatusResponse, error) {
	statusMsg, err := h.appointmentService.UpdateAppointmentStatus(ctx, req.AppointmentId, fromPbAppointmentStatus(req.Status))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdateAppointmentStatusResponse{Status: statusMsg}, nil
}

func (h *AppointmentGrpcHandler) GetAppointmentDetails(ctx context.Context, req *pb.GetAppointmentDetailsRequest) (*pb.GetAppointmentDetailsResponse, error) {
	appointment, details, err := h.appointmentService.GetAppointmentDetails(ctx, req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	pbDetails := make([]*pb.AppointmentDetail, len(details))
	for i, d := range details {
		pbDetails[i] = toProtoAppointmentDetail(&d)
	}
	return &pb.GetAppointmentDetailsResponse{
		Appointment: toProtoAppointment(appointment),
		Details:     pbDetails,
	}, nil
}

// --- DỊCH VỤ ---
func (h *AppointmentGrpcHandler) CreateService(ctx context.Context, req *pb.CreateServiceRequest) (*pb.CreateServiceResponse, error) {
	serviceID, statusMsg, err := h.appointmentService.CreateService(ctx, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateServiceResponse{ServiceId: serviceID, Status: statusMsg}, nil
}

func (h *AppointmentGrpcHandler) GetServices(ctx context.Context, req *pb.GetServicesRequest) (*pb.GetServicesResponse, error) {
	services, err := h.appointmentService.GetServices(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pbServices := make([]*pb.Service, len(services))
	for i, s := range services {
		pbServices[i] = toProtoService(&s)
	}
	return &pb.GetServicesResponse{Services: pbServices}, nil
}

func (h *AppointmentGrpcHandler) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest) (*pb.UpdateServiceResponse, error) {
	statusMsg, err := h.appointmentService.UpdateService(ctx, req.ServiceId, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdateServiceResponse{Status: statusMsg}, nil
}

func (h *AppointmentGrpcHandler) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest) (*pb.DeleteServiceResponse, error) {
	statusMsg, err := h.appointmentService.DeleteService(ctx, req.ServiceId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteServiceResponse{Status: statusMsg}, nil
}
