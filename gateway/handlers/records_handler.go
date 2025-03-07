package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/records"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecordsHandler struct {
	client pb.PetRecordServiceClient
}

func NewRecordsHandler(client pb.PetRecordServiceClient) *RecordsHandler {
	return &RecordsHandler{client: client}
}

// RegisterRoutes đăng ký các route cho PetRecord service
func (h *RecordsHandler) RegisterRoutes(e *echo.Echo) {
	// Pet routes
	e.POST("/pets", h.CreatePet)
	e.GET("/pets/:id", h.GetPet)
	e.PUT("/pets", h.UpdatePet)
	e.DELETE("/pets/:id", h.DeletePet)
	e.GET("/pets/owner/:owner_id", h.ListPets)

	// Examination routes
	e.POST("/examinations", h.CreateExamination)
	e.GET("/examinations/:id", h.GetExamination)
	e.PUT("/examinations", h.UpdateExamination)
	e.DELETE("/examinations/:id", h.DeleteExamination)
	e.GET("/examinations/pet/:pet_id", h.ListExaminations)

	// Vaccination routes
	e.POST("/vaccinations", h.CreateVaccination)
	e.GET("/vaccinations/:id", h.GetVaccination)
	e.PUT("/vaccinations", h.UpdateVaccination)
	e.DELETE("/vaccinations/:id", h.DeleteVaccination)
	e.GET("/vaccinations/pet/:pet_id", h.ListVaccinations)

	// Prescription routes
	e.POST("/prescriptions", h.CreatePrescription)
	e.GET("/prescriptions/:id", h.GetPrescription)
	e.PUT("/prescriptions", h.UpdatePrescription)
	e.DELETE("/prescriptions/:id", h.DeletePrescription)
	e.GET("/prescriptions/pet/:pet_id", h.ListPrescriptions)
}

// --- Pet Methods ---

func (h *RecordsHandler) CreatePet(c echo.Context) error {
	var req struct {
		Name    string  `json:"name"`
		Species string  `json:"species"`
		Age     int32   `json:"age"`
		OwnerID string  `json:"owner_id"`
		Color   string  `json:"color"`
		Weight  float32 `json:"weight"`
		Size    string  `json:"size"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validation cơ bản
	if req.Name == "" || req.Species == "" || req.OwnerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name, species, and owner_id are required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.CreatePet(ctx, &pb.CreatePetRequest{
		Name:    req.Name,
		Species: req.Species,
		Age:     req.Age,
		OwnerId: req.OwnerID,
		Color:   req.Color,
		Weight:  req.Weight,
		Size:    req.Size,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": resp.Id})
}

func (h *RecordsHandler) GetPet(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.GetPet(ctx, &pb.GetPetRequest{Id: id})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Pet not found"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Pet)
}

func (h *RecordsHandler) UpdatePet(c echo.Context) error {
	var req struct {
		ID      string  `json:"id"`
		Name    string  `json:"name"`
		Species string  `json:"species"`
		Age     int32   `json:"age"`
		OwnerID string  `json:"owner_id"`
		Color   string  `json:"color"`
		Weight  float32 `json:"weight"`
		Size    string  `json:"size"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdatePet(ctx, &pb.UpdatePetRequest{
		Id:      req.ID,
		Name:    req.Name,
		Species: req.Species,
		Age:     req.Age,
		OwnerId: req.OwnerID,
		Color:   req.Color,
		Weight:  req.Weight,
		Size:    req.Size,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Pet)
}

func (h *RecordsHandler) DeletePet(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.DeletePet(ctx, &pb.DeletePetRequest{Id: id})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]bool{"success": resp.Success})
}

func (h *RecordsHandler) ListPets(c echo.Context) error {
	ownerID := c.Param("owner_id")
	if ownerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Owner ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.ListPets(ctx, &pb.ListPetsRequest{OwnerId: ownerID})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Pets)
}

// --- Examination Methods ---

func (h *RecordsHandler) CreateExamination(c echo.Context) error {
	var req struct {
		PetID     string `json:"pet_id"`
		Date      string `json:"date"` // Format: "2006-01-02"
		VetID     string `json:"vet_id"`
		Diagnosis string `json:"diagnosis"`
		Notes     string `json:"notes"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.PetID == "" || req.Date == "" || req.VetID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pet ID, date, and vet ID are required"})
	}

	// Parse date
	_, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format, must be YYYY-MM-DD"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.CreateExamination(ctx, &pb.CreateExaminationRequest{
		PetId:     req.PetID,
		Date:      req.Date,
		VetId:     req.VetID,
		Diagnosis: req.Diagnosis,
		Notes:     req.Notes,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid arguments"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": resp.Id})
}

func (h *RecordsHandler) GetExamination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.GetExamination(ctx, &pb.GetExaminationRequest{Id: id})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Examination not found"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Examination)
}

func (h *RecordsHandler) UpdateExamination(c echo.Context) error {
	var req struct {
		ID        string `json:"id"`
		PetID     string `json:"pet_id"`
		Date      string `json:"date"` // Format: "2006-01-02"
		VetID     string `json:"vet_id"`
		Diagnosis string `json:"diagnosis"`
		Notes     string `json:"notes"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	// Parse date
	if req.Date != "" {
		_, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format, must be YYYY-MM-DD"})
		}
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdateExamination(ctx, &pb.UpdateExaminationRequest{
		Id:        req.ID,
		PetId:     req.PetID,
		Date:      req.Date,
		VetId:     req.VetID,
		Diagnosis: req.Diagnosis,
		Notes:     req.Notes,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Examination)
}

func (h *RecordsHandler) DeleteExamination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.DeleteExamination(ctx, &pb.DeleteExaminationRequest{Id: id})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]bool{"success": resp.Success})
}

func (h *RecordsHandler) ListExaminations(c echo.Context) error {
	petID := c.Param("pet_id")
	if petID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pet ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.ListExaminations(ctx, &pb.ListExaminationsRequest{PetId: petID})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Examinations)
}

// --- Vaccination Methods ---

func (h *RecordsHandler) CreateVaccination(c echo.Context) error {
	var req struct {
		PetID       string `json:"pet_id"`
		VaccineName string `json:"vaccine_name"`
		Date        string `json:"date"` // Format: "2006-01-02"
		VetID       string `json:"vet_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.PetID == "" || req.VaccineName == "" || req.Date == "" || req.VetID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pet ID, vaccine name, date, and vet ID are required"})
	}

	// Parse date
	_, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format, must be YYYY-MM-DD"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.CreateVaccination(ctx, &pb.CreateVaccinationRequest{
		PetId:       req.PetID,
		VaccineName: req.VaccineName,
		Date:        req.Date,
		VetId:       req.VetID,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid arguments"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": resp.Id})
}

func (h *RecordsHandler) GetVaccination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.GetVaccination(ctx, &pb.GetVaccinationRequest{Id: id})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Vaccination not found"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Vaccination)
}

func (h *RecordsHandler) UpdateVaccination(c echo.Context) error {
	var req struct {
		ID          string `json:"id"`
		PetID       string `json:"pet_id"`
		VaccineName string `json:"vaccine_name"`
		Date        string `json:"date"` // Format: "2006-01-02"
		VetID       string `json:"vet_id"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	// Parse date
	if req.Date != "" {
		_, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format, must be YYYY-MM-DD"})
		}
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdateVaccination(ctx, &pb.UpdateVaccinationRequest{
		Id:          req.ID,
		PetId:       req.PetID,
		VaccineName: req.VaccineName,
		Date:        req.Date,
		VetId:       req.VetID,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Vaccination)
}

func (h *RecordsHandler) DeleteVaccination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.DeleteVaccination(ctx, &pb.DeleteVaccinationRequest{Id: id})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]bool{"success": resp.Success})
}

func (h *RecordsHandler) ListVaccinations(c echo.Context) error {
	petID := c.Param("pet_id")
	if petID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pet ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.ListVaccinations(ctx, &pb.ListVaccinationsRequest{PetId: petID})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Vaccinations)
}

// --- Prescription Methods ---

func (h *RecordsHandler) CreatePrescription(c echo.Context) error {
	var req struct {
		PetID      string `json:"pet_id"`
		Medication string `json:"medication"`
		Dosage     string `json:"dosage"`
		StartDate  string `json:"start_date"` // Format: "2006-01-02"
		EndDate    string `json:"end_date"`   // Format: "2006-01-02"
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.PetID == "" || req.Medication == "" || req.Dosage == "" || req.StartDate == "" || req.EndDate == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pet ID, medication, dosage, start_date, and end_date are required"})
	}

	// Parse dates
	_, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date format, must be YYYY-MM-DD"})
	}
	_, err = time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date format, must be YYYY-MM-DD"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.CreatePrescription(ctx, &pb.CreatePrescriptionRequest{
		PetId:      req.PetID,
		Medication: req.Medication,
		Dosage:     req.Dosage,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid arguments"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": resp.Id})
}

func (h *RecordsHandler) GetPrescription(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.GetPrescription(ctx, &pb.GetPrescriptionRequest{Id: id})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Prescription not found"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Prescription)
}

func (h *RecordsHandler) UpdatePrescription(c echo.Context) error {
	var req struct {
		ID         string `json:"id"`
		PetID      string `json:"pet_id"`
		Medication string `json:"medication"`
		Dosage     string `json:"dosage"`
		StartDate  string `json:"start_date"` // Format: "2006-01-02"
		EndDate    string `json:"end_date"`   // Format: "2006-01-02"
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	// Parse dates
	if req.StartDate != "" {
		_, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date format, must be YYYY-MM-DD"})
		}
	}
	if req.EndDate != "" {
		_, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date format, must be YYYY-MM-DD"})
		}
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdatePrescription(ctx, &pb.UpdatePrescriptionRequest{
		Id:         req.ID,
		PetId:      req.PetID,
		Medication: req.Medication,
		Dosage:     req.Dosage,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Prescription)
}

func (h *RecordsHandler) DeletePrescription(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.DeletePrescription(ctx, &pb.DeletePrescriptionRequest{Id: id})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]bool{"success": resp.Success})
}

func (h *RecordsHandler) ListPrescriptions(c echo.Context) error {
	petID := c.Param("pet_id")
	if petID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pet ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.ListPrescriptions(ctx, &pb.ListPrescriptionsRequest{PetId: petID})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Prescriptions)
}
