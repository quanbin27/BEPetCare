package main

import (
	"context"
	"errors"

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
func (s *Store) CreateAppointment(ctx context.Context, appointment *Appointment, serviceIDs []int32) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Tạo lịch hẹn
		if err := tx.Create(appointment).Error; err != nil {
			return err
		}

		// Lấy giá dịch vụ từ bảng `services`
		var services []Service
		if err := tx.Where("id IN ?", serviceIDs).Find(&services).Error; err != nil {
			return err
		}

		// Lưu chi tiết dịch vụ vào `appointment_details`
		var details []AppointmentDetail
		servicePriceMap := make(map[int32]float32)
		for _, svc := range services {
			servicePriceMap[svc.ID] = svc.Price
		}

		for _, serviceID := range serviceIDs {
			price, exists := servicePriceMap[serviceID]
			if !exists {
				return errors.New("service ID not found")
			}

			details = append(details, AppointmentDetail{
				AppointmentID: appointment.ID,
				ServiceID:     serviceID,
				ServicePrice:  price,
			})
		}

		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		return nil
	})
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
