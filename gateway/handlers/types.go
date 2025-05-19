package handlers

import (
	pb "github.com/quanbin27/commons/genproto/appointments"
	"time"
)

type ServiceResponse struct {
	ServiceID   int32   `json:"serviceId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	ImgURL      string  `json:"imgUrl"`
}
type UserResponse struct {
	UserID      int32  `json:"userId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}
type AllProductResponse struct {
	ID           int32   `json:"id"`
	Name         string  `json:"name"`
	Price        float32 `json:"price"`
	Description  string  `json:"description"`
	ImgURL       string  `json:"imgUrl"`
	ProductType  string  `json:"productType"`
	IsAttachable bool    `json:"isAttachable"`
}
type ProductResponse struct {
	ID           int32   `json:"id"`
	Name         string  `json:"name"`
	Price        float32 `json:"price"`
	Description  string  `json:"description"`
	ImgURL       string  `json:"imgUrl"`
	IsAttachable bool    `json:"isAttachable"`
}
type AppointmentResponse struct {
	ID              int32   `json:"id"`
	CustomerID      int32   `json:"customer_id"`
	EmployeeID      int32   `json:"employee_id"`
	CustomerAddress string  `json:"customer_address"`
	ScheduledTime   string  `json:"scheduled_time"` // ISO 8601
	Status          string  `json:"status"`
	Note            string  `json:"note"`
	Total           float32 `json:"total"`
	BranchID        int32   `json:"branch_id"`
}

func toAppointmentResponse(pbApp *pb.Appointment) AppointmentResponse {
	return AppointmentResponse{
		ID:              pbApp.Id,
		CustomerID:      pbApp.CustomerId,
		EmployeeID:      pbApp.EmployeeId,
		CustomerAddress: pbApp.CustomerAddress,
		ScheduledTime:   pbApp.ScheduledTime.AsTime().Format(time.RFC3339),
		Status:          pb.AppointmentStatus_name[int32(pbApp.Status)],
		Note:            pbApp.Note,
		Total:           pbApp.Total,
		BranchID:        pbApp.BranchId,
	}
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}
