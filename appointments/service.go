package main

import (
	"context"
	"fmt"
	"time"
)

type AppService struct {
	store AppointmentStore
}

func NewAppointmentService(store AppointmentStore) AppointmentService {
	return &AppService{store: store}
}

// --- LỊCH HẸN ---
// Tạo lịch hẹn
func (s *AppService) CreateAppointment(ctx context.Context, customerID int32, customerAddress string, scheduledTime time.Time, services []AppointmentDetail, note string) (int32, string, error) {
	// Lấy danh sách service IDs từ request
	var serviceIDs []int32
	for _, item := range services {
		serviceIDs = append(serviceIDs, item.ServiceID)
	}

	// Truy vấn database để lấy giá dịch vụ
	serviceList, err := s.store.GetServicesByIDs(ctx, serviceIDs)
	if err != nil {
		return 0, "Failed", err
	}

	// Tạo map giá dịch vụ từ database
	servicePriceMap := make(map[int32]float32)
	for _, svc := range serviceList {
		servicePriceMap[svc.ID] = svc.Price
	}

	// Tính lại total dựa trên giá từ database
	var total float32 = 0
	for i, item := range services {
		price, exists := servicePriceMap[item.ServiceID]
		if !exists {
			return 0, "Failed", fmt.Errorf("service ID %d not found", item.ServiceID)
		}
		services[i].ServicePrice = price // Gán giá từ database vào struct
		total += float32(item.Quantity) * price
	}

	// Gọi Store để tạo lịch hẹn
	id, err := s.store.CreateAppointment(ctx, customerID, customerAddress, scheduledTime, services, total, note)
	if err != nil {
		return 0, "Failed", err
	}

	return id, "Success", nil
}

// Lấy lịch hẹn theo khách hàng
func (s *AppService) GetAppointmentsByCustomer(ctx context.Context, customerID int32) ([]Appointment, error) {
	return s.store.GetAppointmentsByCustomer(ctx, customerID)
}

// Lấy lịch hẹn theo nhân viên
func (s *AppService) GetAppointmentsByEmployee(ctx context.Context, employeeID int32) ([]Appointment, error) {
	return s.store.GetAppointmentsByEmployee(ctx, employeeID)
}

// Cập nhật trạng thái lịch hẹn
func (s *AppService) UpdateAppointmentStatus(ctx context.Context, appointmentID int32, status AppointmentStatus) (string, error) {
	if err := s.store.UpdateAppointmentStatus(ctx, appointmentID, status); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

// Lấy chi tiết lịch hẹn
func (s *AppService) GetAppointmentDetails(ctx context.Context, appointmentID int32) (*Appointment, []AppointmentDetail, error) {
	return s.store.GetAppointmentDetails(ctx, appointmentID)
}

// --- DỊCH VỤ ---
// Tạo dịch vụ
func (s *AppService) CreateService(ctx context.Context, name, description string, price float32) (int32, string, error) {
	service := &Service{
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateService(ctx, service); err != nil {
		return 0, "Failed", err
	}
	return service.ID, "Success", nil
}

// Lấy danh sách dịch vụ
func (s *AppService) GetServices(ctx context.Context) ([]Service, error) {
	return s.store.GetServices(ctx)
}

// Cập nhật dịch vụ
func (s *AppService) UpdateService(ctx context.Context, serviceID int32, name, description string, price float32) (string, error) {
	service := &Service{
		ID:          serviceID,
		Name:        name,
		Description: description,
		Price:       price,
		UpdatedAt:   time.Now(),
	}
	if err := s.store.UpdateService(ctx, service); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

// Xóa dịch vụ
func (s *AppService) DeleteService(ctx context.Context, serviceID int32) (string, error) {
	if err := s.store.DeleteService(ctx, serviceID); err != nil {
		return "Failed", err
	}
	return "Success", nil
}
