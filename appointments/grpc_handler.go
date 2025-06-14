package main

import (
	"context"
	"fmt"
	"github.com/quanbin27/commons/config"
	pb "github.com/quanbin27/commons/genproto/appointments"
	pbUser "github.com/quanbin27/commons/genproto/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
)

type AppointmentGrpcHandler struct {
	appointmentService AppointmentService
	pb.UnimplementedAppointmentServiceServer
	userClient pbUser.UserServiceClient
}

func NewAppointmentGrpcHandler(grpcServer *grpc.Server, appointmentService AppointmentService) {

	usersServiceAddr := config.Envs.UsersGrpcAddr
	usersConn, err := grpc.NewClient(usersServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial user server: %v", err)
	}
	grpcHandler := &AppointmentGrpcHandler{
		appointmentService: appointmentService,
		userClient:         pbUser.NewUserServiceClient(usersConn),
	}
	pb.RegisterAppointmentServiceServer(grpcServer, grpcHandler)
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

	log.Printf("Appointments fetched from service for customerID %d: %+v", req.CustomerId, appointments)

	pbAppointments := make([]*pb.Appointment, len(appointments))
	for i, a := range appointments {
		pbAppointments[len(appointments)-1-i] = toProtoAppointment(&a)
	}

	return &pb.GetAppointmentsResponse{Appointments: pbAppointments}, nil
}

func (h *AppointmentGrpcHandler) GetAppointmentsByEmployee(ctx context.Context, req *pb.GetAppointmentsByEmployeeRequest) (*pb.GetAppointmentsResponse, error) {
	appointments, err := h.appointmentService.GetAppointmentsByEmployee(ctx, req.EmployeeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	fmt.Println("---", appointments)
	pbAppointments := make([]*pb.Appointment, len(appointments))
	for i, a := range appointments {
		pbAppointments[len(appointments)-1-i] = toProtoAppointment(&a)
	}
	return &pb.GetAppointmentsResponse{Appointments: pbAppointments}, nil
}
func (h *AppointmentGrpcHandler) GetAppointmentsByBranch(ctx context.Context, req *pb.GetAppointmentsByBranchRequest) (*pb.GetAppointmentsResponse, error) {
	appointments, err := h.appointmentService.GetAppointmentsByBranch(ctx, req.BranchId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pbAppointments := make([]*pb.Appointment, len(appointments))
	for i, a := range appointments {
		pbAppointments[len(appointments)-1-i] = toProtoAppointment(&a)
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
func (h *AppointmentGrpcHandler) UpdateEmployeeForAppointment(ctx context.Context, req *pb.UpdateEmployeeForAppointmentRequest) (*pb.UpdateEmployeeForAppointmentResponse, error) {
	statusMsg, err := h.appointmentService.UpdateEmployeeForAppointment(ctx, req.AppointmentId, req.EmployeeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdateEmployeeForAppointmentResponse{Status: statusMsg}, nil
}

// Updated GetAppointmentDetails handler to include service names
func (h *AppointmentGrpcHandler) GetAppointmentDetails(ctx context.Context, req *pb.GetAppointmentDetailsRequest) (*pb.GetAppointmentDetailsResponse, error) {
	appointment, detailsWithService, err := h.appointmentService.GetAppointmentDetails(ctx, req.AppointmentId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	pbDetailsWithService := make([]*pb.AppointmentDetailWithService, len(detailsWithService))
	for i, d := range detailsWithService {
		pbDetailsWithService[i] = toProtoAppointmentDetailWithService(&d.AppointmentDetail, d.ServiceName)
	}

	return &pb.GetAppointmentDetailsResponse{
		Appointment: toProtoAppointment(appointment),
		Details:     pbDetailsWithService,
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
func (h *AppointmentGrpcHandler) GetAllAppointments(ctx context.Context, req *pb.GetAllAppointmentsRequest) (*pb.GetAllAppointmentsResponse, error) {
	appointments, err := h.appointmentService.GetAllAppointments(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	pbAppointments := make([]*pb.AppointmentWithCustomerName, len(appointments))
	for i, a := range appointments {
		customer, err := h.userClient.GetUserInfo(ctx, &pbUser.GetUserInfoRequest{ID: a.CustomerID})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get customer info: %v", err)
		}
		employee, err := h.userClient.GetUserInfo(ctx, &pbUser.GetUserInfoRequest{ID: a.EmployeeID})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get employee info: %v", err)
		}
		pbAppointments[len(appointments)-1-i] = &pb.AppointmentWithCustomerName{
			Appointment:   toProtoAppointment(&a),
			CustomerName:  customer.Name,
			CustomerEmail: customer.Email,
			EmployeeEmail: employee.Email,
			EmployeeName:  employee.Name,
		}
	}
	return &pb.GetAllAppointmentsResponse{Appointments: pbAppointments}, nil
}
