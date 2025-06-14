package main

import (
	"context"
	"log"
	"time"

	pb "github.com/quanbin27/commons/genproto/appointments"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- ENUM TRẠNG THÁI LỊCH HẸN ---
type AppointmentStatus string

const (
	StatusPending    AppointmentStatus = "pending"
	StatusInProgress AppointmentStatus = "in_progress"
	StatusCompleted  AppointmentStatus = "completed"
	StatusCancelled  AppointmentStatus = "cancelled"
)

// --- BẢNG DỊCH VỤ ---
type Service struct {
	ID          int32     `gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Price       float32   `gorm:"not null"`
	ImgUrl      string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// --- BẢNG LỊCH HẸN ---
type Appointment struct {
	ID              int32             `gorm:"primaryKey"`
	CustomerID      int32             `gorm:"not null;index"`
	EmployeeID      int32             `gorm:"index"`
	CustomerAddress string            `gorm:"type:text;not null"`
	ScheduledTime   time.Time         `gorm:"not null"`
	Status          AppointmentStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	Total           float32           `gorm:"not null"`
	Note            string
	BranchID        int32               `gorm:"index"`
	CreatedAt       time.Time           `gorm:"autoCreateTime"`
	UpdatedAt       time.Time           `gorm:"autoUpdateTime"`
	Details         []AppointmentDetail `gorm:"foreignKey:AppointmentID"`
}

// --- BẢNG CHI TIẾT LỊCH HẸN ---
type AppointmentDetail struct {
	AppointmentID int32       `gorm:"primaryKey"`
	ServiceID     int32       `gorm:"primaryKey"`
	ServicePrice  float32     `gorm:"not null"`
	Quantity      int32       `gorm:"not null;default:0"`
	Appointment   Appointment `gorm:"foreignKey:AppointmentID;constraint:OnDelete:CASCADE"`
	Service       Service     `gorm:"foreignKey:ServiceID"`
}

// --- INTERFACE CHO APPOINTMENT STORE ---
type AppointmentStore interface {
	// Lịch hẹn
	CreateAppointment(ctx context.Context, customerID int32, customerAddress string, scheduledTime time.Time, services []AppointmentDetail, total float32, note string, branchID int32) (int32, error)
	GetAppointmentsByCustomer(ctx context.Context, customerID int32) ([]Appointment, error)
	GetAppointmentsByEmployee(ctx context.Context, employeeID int32) ([]Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, appointmentID int32, status AppointmentStatus) error
	GetAppointmentDetails(ctx context.Context, appointmentID int32) (*Appointment, []AppointmentDetail, error)
	GetAppointmentsByBranch(ctx context.Context, branchID int32) ([]Appointment, error)
	// Dịch vụ
	CreateService(ctx context.Context, service *Service) error
	GetServices(ctx context.Context) ([]Service, error)
	UpdateService(ctx context.Context, service *Service) error
	DeleteService(ctx context.Context, serviceID int32) error
	GetServicesByIDs(ctx context.Context, serviceIDs []int32) ([]Service, error)
	GetServiceByID(ctx context.Context, serviceID int32) (Service, error)
	UpdateAppointmentEmployee(ctx context.Context, appointmentID, employeeID int32) error
	GetAllAppointments(ctx context.Context) ([]Appointment, error)
}

// --- INTERFACE CHO APPOINTMENT SERVICE (SỬ DỤNG DỮ LIỆU NỘI BỘ) ---
type AppointmentService interface {
	// Lịch hẹn
	CreateAppointment(ctx context.Context, customerID int32, customerAddress string, scheduledTime time.Time, details []AppointmentDetail, note string, branchID int32) (int32, string, error)
	GetAppointmentsByCustomer(ctx context.Context, customerID int32) ([]Appointment, error)
	GetAppointmentsByEmployee(ctx context.Context, employeeID int32) ([]Appointment, error)
	UpdateAppointmentStatus(ctx context.Context, appointmentID int32, status AppointmentStatus) (string, error)
	GetAppointmentsByBranch(ctx context.Context, branchID int32) ([]Appointment, error)
	// Updated to include service information
	GetAppointmentDetails(ctx context.Context, appointmentID int32) (*Appointment, []AppointmentDetailWithService, error)
	UpdateEmployeeForAppointment(ctx context.Context, appointmentID, employeeID int32) (string, error)
	// Service-related methods
	CreateService(ctx context.Context, name, description string, price float32) (int32, string, error)
	GetServices(ctx context.Context) ([]Service, error)
	UpdateService(ctx context.Context, serviceID int32, name, description string, price float32) (string, error)
	DeleteService(ctx context.Context, serviceID int32) (string, error) // Trả về status
	GetAllAppointments(ctx context.Context) ([]Appointment, error)
}

// --- CHUYỂN ĐỔI ENUM PROTO <-> GO ---
func toPbAppointmentStatus(status AppointmentStatus) pb.AppointmentStatus {
	log.Printf("Status before converting: %v", status) // Kiểm tra Status trước khi chuyển đổi

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
		return StatusPending // Mặc định PENDING nếu không xác định
	}
}

// --- HÀM CHUYỂN ĐỔI GIỮA DỮ LIỆU NỘI BỘ VÀ PROTOBUF ---
func toProtoAppointment(a *Appointment) *pb.Appointment {
	// Log ra giá trị của Appointment trước khi chuyển đổi
	log.Printf("Converting Appointment: %+v", a)

	// Thực hiện chuyển đổi sang pb.Appointment
	pbApp := &pb.Appointment{
		Id:              a.ID,
		CustomerId:      a.CustomerID,
		EmployeeId:      a.EmployeeID,
		CustomerAddress: a.CustomerAddress,
		ScheduledTime:   timestamppb.New(a.ScheduledTime),
		Status:          toPbAppointmentStatus(a.Status),
		BranchId:        a.BranchID,
		Note:            a.Note,
		Total:           a.Total,
	}
	log.Printf("Converted pb.Appointment: id=%d, status=%v", pbApp.Id, pbApp.Status)
	// Log ra giá trị sau khi chuyển đổi

	// Trả về kết quả
	return pbApp
}

func toProtoService(s *Service) *pb.Service {
	return &pb.Service{
		Id:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Price:       s.Price,
		Imgurl:      s.ImgUrl,
	}
}

func toProtoAppointmentDetail(ad *AppointmentDetail) *pb.AppointmentDetail {
	return &pb.AppointmentDetail{
		AppointmentId: ad.AppointmentID,
		ServiceId:     ad.ServiceID,
		ServicePrice:  ad.ServicePrice,
	}
}

// Helper function to convert domain AppointmentDetail to proto AppointmentDetailWithService
func toProtoAppointmentDetailWithService(detail *AppointmentDetail, serviceName string) *pb.AppointmentDetailWithService {
	return &pb.AppointmentDetailWithService{
		AppointmentId: detail.AppointmentID,
		ServiceId:     detail.ServiceID,
		ServiceName:   serviceName,
		ServicePrice:  detail.ServicePrice,
		Quantity:      detail.Quantity,
	}
}

type AppointmentDetailWithService struct {
	AppointmentDetail
	ServiceName string
}
