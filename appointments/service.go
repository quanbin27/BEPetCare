package main

import (
	"context"
	"fmt"
	"time"
)

type AppService struct {
	store         AppointmentStore
	priceStrategy PriceCalculationStrategy
}

func NewAppointmentService(store AppointmentStore, strategy PriceCalculationStrategy) AppointmentService {
	return &AppService{store: store, priceStrategy: strategy}
}

type PriceCalculationStrategy interface {
	CalculateTotal(services []AppointmentDetail, servicePriceMap map[int32]float32) (float32, error)
}
type StandardPriceStrategy struct{}

func (s *StandardPriceStrategy) CalculateTotal(services []AppointmentDetail, servicePriceMap map[int32]float32) (float32, error) {
	var total float32
	for _, item := range services {
		price, exists := servicePriceMap[item.ServiceID]
		if !exists {
			return 0, fmt.Errorf("service ID %d not found", item.ServiceID)
		}
		total += float32(item.Quantity) * price
	}
	return total, nil
}

type DiscountPriceStrategy struct {
	discount float32
}

func (s *DiscountPriceStrategy) CalculateTotal(services []AppointmentDetail, servicePriceMap map[int32]float32) (float32, error) {
	var total float32
	for _, item := range services {
		price, exists := servicePriceMap[item.ServiceID]
		if !exists {
			return 0, fmt.Errorf("service ID %d not found", item.ServiceID)
		}
		total += float32(item.Quantity) * price
	}
	return total * (1 - s.discount), nil
}

// --- LỊCH HẸN ---
// Tạo lịch hẹn
func (s *AppService) CreateAppointment(ctx context.Context, customerID int32, customerAddress string, scheduledTime time.Time, services []AppointmentDetail, note string, branchID int32) (int32, string, error) {
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
	for i, item := range services {
		price, exists := servicePriceMap[item.ServiceID]
		if !exists {
			return 0, "Failed", fmt.Errorf("service ID %d not found", item.ServiceID)
		}
		services[i].ServicePrice = price
	}
	total, err := s.priceStrategy.CalculateTotal(services, servicePriceMap)

	// Gọi Store để tạo lịch hẹn
	id, err := s.store.CreateAppointment(ctx, customerID, customerAddress, scheduledTime, services, total, note, branchID)
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

// Lấy lịch hẹn theo nhân viên
func (s *AppService) GetAppointmentsByBranch(ctx context.Context, branchID int32) ([]Appointment, error) {
	return s.store.GetAppointmentsByBranch(ctx, branchID)
}

// Cập nhật trạng thái lịch hẹn
func (s *AppService) UpdateAppointmentStatus(ctx context.Context, appointmentID int32, status AppointmentStatus) (string, error) {
	if err := s.store.UpdateAppointmentStatus(ctx, appointmentID, status); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *AppService) GetAppointmentDetails(ctx context.Context, appointmentID int32) (*Appointment, []AppointmentDetailWithService, error) {
	// First, get the appointment
	appointment, details, err := s.store.GetAppointmentDetails(ctx, appointmentID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	// Enhance details with service names
	detailsWithService := make([]AppointmentDetailWithService, len(details))
	for i, detail := range details {
		// Get service name from service ID
		service, err := s.store.GetServiceByID(ctx, detail.ServiceID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get service information: %w", err)
		}

		detailsWithService[i] = AppointmentDetailWithService{
			AppointmentDetail: detail,
			ServiceName:       service.Name,
		}
	}

	return appointment, detailsWithService, nil
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
func (s *AppService) UpdateEmployeeForAppointment(ctx context.Context, appointmentID, employeeID int32) (string, error) {
	if err := s.store.UpdateAppointmentEmployee(ctx, appointmentID, employeeID); err != nil {
		return "Failed", err
	}
	return "Success", nil
}
func (s *AppService) GetAllAppointments(ctx context.Context) ([]Appointment, error) {
	return s.store.GetAllAppointments(ctx)
}
