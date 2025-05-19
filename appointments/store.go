package main

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// --- APPOINTMENT STORE ---
type Store struct {
	db *gorm.DB
}

// NewStore khởi tạo Store
func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// Tạo lịch hẹn + chi tiết dịch vụ
func (s *Store) CreateAppointment(ctx context.Context, customerID int32, customerAddress string, scheduledTime time.Time, services []AppointmentDetail, total float32, note string, branchID int32) (int32, error) {
	var appointmentID int32

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Tạo lịch hẹn
		appointment := Appointment{
			CustomerID:      customerID,
			CustomerAddress: customerAddress,
			ScheduledTime:   scheduledTime,
			Status:          StatusPending,
			Total:           total,
			Note:            note,
			BranchID:        branchID,
		}

		if err := tx.Create(&appointment).Error; err != nil {
			return err
		}

		// Chuẩn bị danh sách chi tiết lịch hẹn
		var details []AppointmentDetail
		for _, svc := range services {
			detail := AppointmentDetail{
				AppointmentID: appointment.ID,
				ServiceID:     svc.ServiceID,
				ServicePrice:  svc.ServicePrice, // Giá đã được lấy từ DB ở tầng AppService
				Quantity:      svc.Quantity,
			}
			details = append(details, detail)
		}

		// Lưu chi tiết lịch hẹn
		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		appointmentID = appointment.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return appointmentID, nil
}
func (s *Store) GetServicesByIDs(ctx context.Context, serviceIDs []int32) ([]Service, error) {
	var services []Service
	if err := s.db.WithContext(ctx).Where("id IN ?", serviceIDs).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}
func (s *Store) GetServiceByID(ctx context.Context, serviceID int32) (Service, error) {
	var service Service
	if err := s.db.WithContext(ctx).First(&service, serviceID).Error; err != nil {
		return Service{}, err
	}
	return service, nil
}

// Lấy lịch hẹn theo khách hàng
func (s *Store) GetAppointmentsByCustomer(ctx context.Context, customerID int32) ([]Appointment, error) {
	var appointments []Appointment
	err := s.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&appointments).Error

	return appointments, err
}

// Lấy lịch hẹn theo nhân viên
func (s *Store) GetAppointmentsByEmployee(ctx context.Context, employeeID int32) ([]Appointment, error) {
	var appointments []Appointment
	err := s.db.WithContext(ctx).Where("employee_id = ?", employeeID).Find(&appointments).Error
	return appointments, err
}

// Cập nhật trạng thái lịch hẹn
func (s *Store) UpdateAppointmentStatus(ctx context.Context, appointmentID int32, status AppointmentStatus) error {
	return s.db.WithContext(ctx).Model(&Appointment{}).
		Where("id = ?", appointmentID).
		Update("status", status).Error
}

// Lấy chi tiết lịch hẹn + dịch vụ
func (s *Store) GetAppointmentDetails(ctx context.Context, appointmentID int32) (*Appointment, []AppointmentDetail, error) {
	var appointment Appointment
	var details []AppointmentDetail

	err := s.db.WithContext(ctx).
		Preload("Details.Service").
		Where("id = ?", appointmentID).
		First(&appointment).Error
	if err != nil {
		return nil, nil, err
	}

	err = s.db.WithContext(ctx).
		Where("appointment_id = ?", appointmentID).
		Find(&details).Error
	if err != nil {
		return &appointment, nil, err
	}

	return &appointment, details, nil
}

// --- SERVICE STORE ---
// Tạo dịch vụ mới
func (s *Store) CreateService(ctx context.Context, service *Service) error {
	return s.db.WithContext(ctx).Create(service).Error
}

// Lấy danh sách dịch vụ
func (s *Store) GetServices(ctx context.Context) ([]Service, error) {
	var services []Service
	err := s.db.WithContext(ctx).Find(&services).Error
	return services, err
}

// Cập nhật dịch vụ
func (s *Store) UpdateService(ctx context.Context, service *Service) error {
	return s.db.WithContext(ctx).
		Model(&Service{}).
		Where("id = ?", service.ID).
		Updates(service).Error
}

// Xóa dịch vụ
func (s *Store) DeleteService(ctx context.Context, serviceID int32) error {
	return s.db.WithContext(ctx).
		Where("id = ?", serviceID).
		Delete(&Service{}).Error
}
