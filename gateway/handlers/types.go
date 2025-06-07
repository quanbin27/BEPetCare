package handlers

import (
	pb "github.com/quanbin27/commons/genproto/appointments"
	"time"
)

// ServiceResponse represents a service for appointment
// @Description Service information
type ServiceResponse struct {
	ServiceID   int32   `json:"serviceId" example:"1" description:"Unique identifier of the service"`
	Name        string  `json:"name" example:"Pet Grooming" description:"Name of the service"`
	Description string  `json:"description" example:"Complete pet grooming service" description:"Detailed description of the service"`
	Price       float32 `json:"price" example:"150.50" description:"Price of the service"`
	ImgURL      string  `json:"imgUrl" example:"https://example.com/images/grooming.jpg" description:"URL to the service image"`
}

// UserResponse represents a user in the system
// @Description User information
type UserResponse struct {
	UserID      int32  `json:"userId" example:"1001" description:"Unique identifier of the user"`
	Name        string `json:"name" example:"John Doe" description:"Full name of the user"`
	Email       string `json:"email" example:"johndoe@example.com" description:"Email address of the user"`
	PhoneNumber string `json:"phoneNumber" example:"+84912345678" description:"Phone number of the user"`
	Address     string `json:"address" example:"123 Main St, City" description:"Address of the user"`
	BranchID    int32  `json:"branchId" example:"301" description:"ID of the branch associated with the user"`
}

// AllProductResponse represents a product with type information
// @Description Complete product information including type
type AllProductResponse struct {
	ID           int32   `json:"id" example:"101" description:"Unique identifier of the product"`
	Name         string  `json:"name" example:"Premium Dog Food" description:"Name of the product"`
	Price        float32 `json:"price" example:"45.99" description:"Price of the product"`
	Description  string  `json:"description" example:"High-quality nutrition for dogs" description:"Detailed description of the product"`
	ImgURL       string  `json:"imgUrl" example:"https://example.com/images/dog-food.jpg" description:"URL to the product image"`
	ProductType  string  `json:"productType" example:"Food" description:"Type of the product"`
	IsAttachable bool    `json:"isAttachable" example:"false" description:"Whether the product can be attached to other products"`
}

// ProductResponse represents a basic product
// @Description Basic product information
type ProductResponse struct {
	ID           int32   `json:"id" example:"102" description:"Unique identifier of the product"`
	Name         string  `json:"name" example:"Pet Collar" description:"Name of the product"`
	Price        float32 `json:"price" example:"15.99" description:"Price of the product"`
	Description  string  `json:"description" example:"Comfortable collar for pets" description:"Detailed description of the product"`
	ImgURL       string  `json:"imgUrl" example:"https://example.com/images/collar.jpg" description:"URL to the product image"`
	IsAttachable bool    `json:"isAttachable" example:"true" description:"Whether the product can be attached to other products"`
}

// AppointmentResponse represents an appointment
// @Description Appointment information
type AppointmentResponse struct {
	ID              int32   `json:"id" example:"5001" description:"Unique identifier of the appointment"`
	CustomerID      int32   `json:"customer_id" example:"1001" description:"ID of the customer who made the appointment"`
	EmployeeID      int32   `json:"employee_id" example:"2001" description:"ID of the employee assigned to the appointment"`
	CustomerAddress string  `json:"customer_address" example:"123 Main St, City" description:"Address where the service will be provided"`
	ScheduledTime   string  `json:"scheduled_time" example:"2025-03-07T10:00:00Z" description:"Scheduled time for the appointment (ISO 8601 format)"`
	Status          string  `json:"status" example:"PENDING" description:"Current status of the appointment"`
	Note            string  `json:"note" example:"Please bring all necessary equipment" description:"Additional notes for the appointment"`
	Total           float32 `json:"total" example:"180.50" description:"Total cost of the appointment"`
	BranchID        int32   `json:"branch_id" example:"301" description:"ID of the branch handling the appointment"`
}

// ServiceItemRequest represents a service item for appointment creation
// @Description Service item for creating an appointment
type ServiceItemRequest struct {
	ServiceID int32 `json:"service_id" example:"1" description:"ID of the service"`
	Quantity  int32 `json:"quantity" example:"1" description:"Quantity of the service"`
}

// CreateAppointmentRequest represents a request to create a new appointment
// @Description Request to create a new appointment
type CreateAppointmentRequest struct {
	CustomerID      int32                `json:"customer_id" example:"1001" description:"ID of the customer making the appointment"`
	CustomerAddress string               `json:"customer_address" example:"123 Main St, City" description:"Address where the service will be provided"`
	ScheduledTime   string               `json:"scheduled_time" example:"2025-03-07T10:00:00Z" description:"Scheduled time for the appointment (ISO 8601 format)"`
	Services        []ServiceItemRequest `json:"services" description:"List of services requested"`
	Note            string               `json:"note" example:"Please arrive on time" description:"Additional notes for the appointment"`
	BranchID        int32                `json:"branch_id" example:"301" description:"ID of the branch handling the appointment"`
}

// CreateAppointmentResponse represents a response after creating an appointment
// @Description Response after creating an appointment
type CreateAppointmentResponse struct {
	AppointmentID int32  `json:"appointment_id" example:"5001" description:"ID of the created appointment"`
	Status        string `json:"status" example:"success" description:"Status of the operation"`
}

// UpdateAppointmentStatusRequest represents a request to update appointment status
// @Description Request to update appointment status
type UpdateAppointmentStatusRequest struct {
	AppointmentID string `json:"appointment_id" example:"5001" description:"ID of the appointment to update"`
	Status        string `json:"status" example:"COMPLETED" description:"New status for the appointment"`
}

// UpdateAppointmentStatusResponse represents a response after updating appointment status
// @Description Response after updating appointment status
type UpdateAppointmentStatusResponse struct {
	Status string `json:"status" example:"success" description:"Status of the operation"`
}

// AppointmentDetailsResponse represents detailed appointment information
// @Description Detailed appointment information including services
type AppointmentDetailsResponse struct {
	Appointment AppointmentResponse               `json:"appointment" description:"Basic appointment information"`
	Details     []pb.AppointmentDetailWithService `json:"details" description:"Details of services included in the appointment"`
	Order       interface{}                       `json:"order" description:"Associated order information if available"`
}

// VerifyEmailRequest represents a request to verify email
// @Description Request to verify email
type VerifyEmailRequest struct {
	Token string `json:"token" example:"abc123def456" description:"Verification token sent to the user's email"`
}

// MedicationResponse represents a medication item in the prescription.
type MedicationResponse struct {
	Name       string `json:"name"`
	Dosage     string `json:"dosage"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	MedicineID string `json:"medicine_id"`
}

// PrescriptionResponse represents the structure of the prescription returned in the API.
type PrescriptionResponse struct {
	ID            string               `json:"id"`
	ExaminationID string               `json:"examination_id"`
	Medications   []MedicationResponse `json:"medications"`
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

// Additional request and response types for Swagger documentation
type CreateServiceRequest struct {
	Name        string  `json:"name" example:"Pet Bathing" description:"Name of the service"`
	Description string  `json:"description" example:"Full bathing service for pets" description:"Detailed description of the service"`
	Price       float32 `json:"price" example:"75.50" description:"Price of the service"`
}

type CreateServiceResponse struct {
	ServiceID int32  `json:"service_id" example:"10" description:"ID of the created service"`
	Status    string `json:"status" example:"success" description:"Status of the operation"`
}

type UpdateServiceRequest struct {
	ServiceID   string  `json:"service_id" example:"10" description:"ID of the service to update"`
	Name        string  `json:"name" example:"Premium Pet Bathing" description:"Updated name of the service"`
	Description string  `json:"description" example:"Enhanced bathing service for pets" description:"Updated description of the service"`
	Price       float32 `json:"price" example:"85.50" description:"Updated price of the service"`
}

type UpdateServiceResponse struct {
	Status string `json:"status" example:"success" description:"Status of the operation"`
}

type DeleteServiceResponse struct {
	Status string `json:"status" example:"success" description:"Status of the operation"`
}
