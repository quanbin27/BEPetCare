package main

import (
	"context"
	"time"

	pb "github.com/quanbin27/commons/genproto/appointments"
)

// --- ENUM TRẠNG THÁI LỊCH HẸN ---
type AppointmentStatus string

const (
	StatusPending    AppointmentStatus = "pending"
	StatusInProgress AppointmentStatus = "in_progress"
	StatusCompleted  AppointmentStatus = "completed"
	StatusCancelled  AppointmentStatus = "cancelled"
)

// --- CHUYỂN ĐỔI ENUM PROTO <-> GO ---
func toPbAppointmentStatus(status AppointmentStatus) pb.AppointmentStatus {
	switch status {
	case StatusPending:
		return pb.AppointmentStatus_PENDING
	case StatusInProgress:
		return pb.AppointmentStatus_IN_PROGRESS
	case StatusCompleted:
		return pb.AppointmentStatus_COMPLETED
	case StatusCancelled:
		return pb.AppointmentStatus_CANCELLED
	default:
		return pb.AppointmentStatus_UNSPECIFIED
	}
}

func fromPbAppointmentStatus(pbStatus pb.AppointmentStatus) AppointmentStatus {
	switch pbStatus {
	case pb.AppointmentStatus_PENDING:
		return StatusPending
	case pb.AppointmentStatus_IN_PROGRESS:
		return StatusInProgress
	case pb.AppointmentStatus_COMPLETED:
		return StatusCompleted
	case pb.AppointmentStatus_CANCELLED:
		return StatusCancelled
	default:
		return "unknown_status"
	}
}

// --- BẢNG DỊCH VỤ ---
type Service struct {
	ID          int32   `gorm:"primaryKey"`
	Name        string  `gorm:"not null"`
	Description string  `gorm:"type:text"`
	Price       float32 `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// --- BẢNG LỊCH HẸN ---
type Appointment struct {
	ID              int32             `gorm:"primaryKey"`
	CustomerID      int32             `gorm:"not null;index"`
	EmployeeID      int32             `gorm:"not null;index"`
	CustomerAddress string            `gorm:"type:text;not null"`
	ScheduledTime   time.Time         `gorm:"not null"`
	Status          AppointmentStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Details []AppointmentDetail `gorm:"foreignKey:AppointmentID"`
}

// --- BẢNG CHI TIẾT LỊCH HẸN ---
type AppointmentDetail struct {
	AppointmentID int32   `gorm:"primaryKey"`
	ServiceID     int32   `gorm:"primaryKey"`
	ServicePrice  float32 `gorm:"not null"`

	Appointment Appointment `gorm:"foreignKey:AppointmentID;constraint:OnDelete:CASCADE"`
	Service     Service     `gorm:"foreignKey:ServiceID"`
}

// --- INTERFACE CHO APPOINTMENT STORE ---
type AppointmentStore interface {
	// Lịch hẹn
	CreateAppointment(ctx context.Context, appointment *Appointment, services []int32) error
	GetAppointmentsByCustomer(ctx context.Context, customerID int32) ([]Appointment, error)
	GetAppointmentsByEmployee(ctx context.Context, employeeID int32) ([]Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, appointmentID int32, status AppointmentStatus) error
	GetAppointmentDetails(ctx context.Context, appointmentID int32) (*Appointment, []AppointmentDetail, error)

	// Dịch vụ
	CreateService(ctx context.Context, service *Service) error
	GetServices(ctx context.Context) ([]Service, error)
	UpdateService(ctx context.Context, service *Service) error
	DeleteService(ctx context.Context, serviceID int32) error
}

// --- INTERFACE CHO APPOINTMENT SERVICE (SỬ DỤNG PROTO) ---
type AppointmentService interface {
	// Lịch hẹn
	CreateAppointment(ctx context.Context, req *pb.CreateAppointmentRequest) (*pb.CreateAppointmentResponse, error)
	GetAppointmentsByCustomer(ctx context.Context, req *pb.GetAppointmentsByCustomerRequest) (*pb.GetAppointmentsResponse, error)
	GetAppointmentsByEmployee(ctx context.Context, req *pb.GetAppointmentsByEmployeeRequest) (*pb.GetAppointmentsResponse, error)
	UpdateAppointmentStatus(ctx context.Context, req *pb.UpdateAppointmentStatusRequest) (*pb.UpdateAppointmentStatusResponse, error)
	GetAppointmentDetails(ctx context.Context, req *pb.GetAppointmentDetailsRequest) (*pb.GetAppointmentDetailsResponse, error)

	// Dịch vụ
	CreateService(ctx context.Context, req *pb.CreateServiceRequest) (*pb.CreateServiceResponse, error)
	GetServices(ctx context.Context, req *pb.GetServicesRequest) (*pb.GetServicesResponse, error)
	UpdateService(ctx context.Context, req *pb.UpdateServiceRequest) (*pb.UpdateServiceResponse, error)
	DeleteService(ctx context.Context, req *pb.DeleteServiceRequest) (*pb.DeleteServiceResponse, error)
}
