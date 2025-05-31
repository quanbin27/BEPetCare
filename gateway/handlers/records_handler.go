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
func (h *RecordsHandler) RegisterRoutes(e *echo.Group) {
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
	e.GET("/prescriptions/examination/:examination_id", h.ListPrescriptions)
}

// --- Pet Methods ---

// CreatePet creates a new pet record
// @Summary Create a new pet
// @Description Creates a new pet record with details like name, species, age, owner ID, etc.
// @Tags Pets
// @Accept json
// @Produce json
// @Param request body object{name=string,species=string,age=integer,owner_id=string,color=string,weight=number,size=string} true "Pet details"
// @Success 200 {object} object{id=string} "Pet created successfully"
// @Failure 400 {object} object{error=string} "Invalid request or missing fields"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /pets [post]
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
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": resp.Id})
}

// GetPet retrieves a pet record by ID
// @Summary Get pet by ID
// @Description Retrieves a pet record using its unique ID
// @Tags Pets
// @Produce json
// @Param id path string true "Pet ID"
// @Success 200 {object} object{id=string,name=string,species=string,age=integer,owner_id=string,color=string,weight=number,size=string} "Pet details"
// @Failure 400 {object} object{error=string} "ID is required"
// @Failure 404 {object} object{error=string} "Pet not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /pets/{id} [get]
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

// UpdatePet updates an existing pet record
// @Summary Update a pet
// @Description Updates a pet record with details like name, species, age, owner ID, etc.
// @Tags Pets
// @Accept json
// @Produce json
// @Param request body object{id=string,name=string,species=string,age=integer,owner_id=string,color=string,weight=number,size=string} true "Updated pet details"
// @Success 200 {object} object{id=string,name=string,species=string,age=integer,owner_id=string,color=string,weight=number,size=string} "Pet updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request or ID is required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /pets [put]
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
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Pet)
}

// DeletePet deletes a pet record
// @Summary Delete a pet
// @Description Deletes a pet record by its unique ID
// @Tags Pets
// @Produce json
// @Param id path string true "Pet ID"
// @Success 200 {object} object{success=boolean} "Pet deleted successfully"
// @Failure 400 {object} object{error=string} "ID is required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /pets/{id} [delete]
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

// ListPets lists all pets for a given owner
// @Summary List pets by owner
// @Description Retrieves a list of pet records for a specific owner ID
// @Tags Pets
// @Produce json
// @Param owner_id path string true "Owner ID"
// @Success 200 {array} object{id=string,name=string,species=string,age=integer,owner_id=string,color=string,weight=number,size=string} "List of pets"
// @Failure 400 {object} object{error=string} "Owner ID is required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /pets/owner/{owner_id} [get]
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

// CreateExamination creates a new examination record
// @Summary Create a new examination
// @Description Creates a new examination record for a pet with details like pet ID, date, vet ID, diagnosis, and notes
// @Tags Examinations
// @Accept json
// @Produce json
// @Param request body object{pet_id=string,date=string,vet_id=string,diagnosis=string,notes=string} true "Examination details"
// @Success 200 {object} object{id=string} "Examination created successfully"
// @Failure 400 {object} object{error=string} "Invalid request, missing required fields, or invalid date format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /examinations [post]
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
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": resp.Id})
}

// GetExamination retrieves an examination record by ID
// @Summary Get examination by ID
// @Description Retrieves an examination record using its unique ID
// @Tags Examinations
// @Produce json
// @Param id path string true "Examination ID"
// @Success 200 {object} object{id=string,pet_id=string,date=string,vet_id=string,diagnosis=string,notes=string} "Examination details"
// @Failure 400 {object} object{error=string} "ID is required"
// @Failure 404 {object} object{error=string} "Examination not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /examinations/{id} [get]
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

// UpdateExamination updates an existing examination record
// @Summary Update an examination
// @Description Updates an examination record with details like pet ID, date, vet ID, diagnosis, and notes
// @Tags Examinations
// @Accept json
// @Produce json
// @Param request body object{id=string,pet_id=string,date=string,vet_id=string,diagnosis=string,notes=string} true "Updated examination details"
// @Success 200 {object} object{id=string,pet_id=string,date=string,vet_id=string,diagnosis=string,notes=string} "Examination updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request, ID is required, or invalid date format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /examinations [put]
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
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Examination)
}

// DeleteExamination deletes an examination record
// @Summary Delete an examination
// @Description Deletes an examination record by its unique ID
// @Tags Examinations
// @Produce json
// @Param id path string true "Examination ID"
// @Success 200 {object} object{success=boolean} "Examination deleted successfully"
// @Failure 400 {object} object{error=string} "ID is required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /examinations/{id} [delete]
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

// ListExaminations lists all examinations for a given pet
// @Summary List examinations by pet
// @Description Retrieves a list of examination records for a specific pet ID
// @Tags Examinations
// @Produce json
// @Param pet_id path string true "Pet ID"
// @Success 200 {array} object{id=string,pet_id=string,date=string,vet_id=string,diagnosis=string,notes=string} "List of examinations"
// @Failure 400 {object} object{error=string} "Pet ID is required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /examinations/pet/{pet_id} [get]
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

// CreateVaccination creates a new vaccination record
// @Summary Create a new vaccination
// @Description Creates a new vaccination record for a pet with details like pet ID, vaccine name, date, next dose, and vet ID
// @Tags Vaccinations
// @Accept json
// @Produce json
// @Param request body object{pet_id=string,vaccine_name=string,date=string,next_dose=string,vet_id=string} true "Vaccination details"
// @Success 200 {object} object{id=string} "Vaccination created successfully"
// @Failure 400 {object} object{error=string} "Invalid request, missing required fields, or invalid date format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /vaccinations [post]
func (h *RecordsHandler) CreateVaccination(c echo.Context) error {
	var req struct {
		PetID       string `json:"pet_id"`
		VaccineName string `json:"vaccine_name"`
		Date        string `json:"date"`      // Format: "2006-01-02"
		NextDose    string `json:"next_dose"` // Format: "2006-01-02"
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
	// Parse next_dose if provided
	if req.NextDose != "" {
		_, err := time.Parse("2006-01-02", req.NextDose)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid next_dose format, must be YYYY-MM-DD"})
		}
	}

	ctx := c.Request().Context()
	resp, err := h.client.CreateVaccination(ctx, &pb.CreateVaccinationRequest{
		PetId:       req.PetID,
		VaccineName: req.VaccineName,
		Date:        req.Date,
		NextDose:    req.NextDose,
		VetId:       req.VetID,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": resp.Id})
}

// GetVaccination retrieves a vaccination record by ID
// @Summary Get vaccination by ID
// @Description Retrieves a vaccination record using its unique ID
// @Tags Vaccinations
// @Produce json
// @Param id path string true "Vaccination ID"
// @Success 200 {object} object{id=string,pet_id=string,vaccine_name=string,date=string,next_dose=string,vet_id=string} "Vaccination details"
// @Failure 400 {object} object{error=string} "ID is required"
// @Failure 404 {object} object{error=string} "Vaccination not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /vaccinations/{id} [get]
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

// UpdateVaccination updates an existing Série record
// @Summary Update a vaccination
// @Description Updates a vaccination record with details like pet ID, vaccine name, date, next dose date, and vet ID
// @Tags Vaccinations
// @Accept json
// @Produce json
// @Param request body object{id=string,pet_id=string,vaccine_name=string,date=string,next_dose=string,vet_id=string} true "Updated vaccination details"
// @Success 200 {object} object{id=string,pet_id=string,vaccine_name=string,date=string,next_dose=string,vet_id=string} "Vaccination updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request, ID is required, or invalid date format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /vaccinations [put]
func (h *RecordsHandler) UpdateVaccination(c echo.Context) error {
	var req struct {
		ID          string `json:"id"`
		PetID       string `json:"pet_id"`
		VaccineName string `json:"vaccine_name"`
		Date        string `json:"date"`      // Format: "2006-01-02"
		NextDose    string `json:"next_dose"` // Format: "2006-01-02"
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
	// Parse date
	if req.NextDose != "" {
		_, err := time.Parse("2006-01-02", req.NextDose)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid next_dose format, must be YYYY-MM-DD"})
		}
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdateVaccination(ctx, &pb.UpdateVaccinationRequest{
		Id:          req.ID,
		PetId:       req.PetID,
		VaccineName: req.VaccineName,
		Date:        req.Date,
		NextDose:    req.NextDose,
		VetId:       req.VetID,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Vaccination)
}

// DeleteVaccination deletes a vaccination record
// @Summary Delete a vaccination
// @Description Deletes a vaccination record by its unique ID
// @Tags Vaccinations
// @Produce json
// @Param id path string true "Vaccination ID"
// @Success 200 {object} object{success=boolean} "Vaccination deleted successfully"
// @Failure 400 {object} object{error=string} "ID is required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /vaccinations/{id} [delete]
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

// ListVaccinations lists all vaccinations for a given pet
// @Summary List vaccinations by pet
// @Description Retrieves a list of vaccination records for a specific pet ID
// @Tags Vaccinations
// @Produce json
// @Param pet_id path string true "Pet ID"
// @Success 200 {array} object{id=string,pet_id=string,vaccine_name=string,date=string,next_dose=string,vet_id=string} "List of vaccinations"
// @Failure 400 {object} object{error=string} "Pet ID is required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /vaccinations/pet/{pet_id} [get]
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

// CreatePrescription creates a new prescription record
// @Summary Create a new prescription
// @Description Creates a new prescription record associated with an examination, including a list of medications
// @Tags Prescriptions
// @Accept json
// @Produce json
// @Param request body object{examination_id=string,medications=[]object{name=string,dosage=string,start_date=string,end_date=string}} true "Prescription details"
// @Success 200 {object} object{id=string} "Prescription created successfully"
// @Failure 400 {object} object{error=string} "Invalid request, missing required fields, or invalid date format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /prescriptions [post]
func (h *RecordsHandler) CreatePrescription(c echo.Context) error {
	var req struct {
		ExaminationID string `json:"examination_id"`
		Medications   []struct {
			Name      string `json:"name"`
			Dosage    string `json:"dosage"`
			StartDate string `json:"start_date"` // Format: "2006-01-02"
			EndDate   string `json:"end_date"`   // Format: "2006-01-02"
		} `json:"medications"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ExaminationID == "" || len(req.Medications) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Examination ID and at least one medication are required"})
	}

	// Validate medications
	medications := make([]*pb.Medication, len(req.Medications))
	for i, med := range req.Medications {
		if med.Name == "" || med.Dosage == "" || med.StartDate == "" || med.EndDate == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Medication name, dosage, start_date, and end_date are required"})
		}
		_, err := time.Parse("2006-01-02", med.StartDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date format, must be YYYY-MM-DD"})
		}
		_, err = time.Parse("2006-01-02", med.EndDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date format, must be YYYY-MM-DD"})
		}
		medications[i] = &pb.Medication{
			Name:      med.Name,
			Dosage:    med.Dosage,
			StartDate: med.StartDate,
			EndDate:   med.EndDate,
		}
	}

	ctx := c.Request().Context()
	resp, err := h.client.CreatePrescription(ctx, &pb.CreatePrescriptionRequest{
		ExaminationId: req.ExaminationID,
		Medications:   medications,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": resp.Id})
}

// GetPrescription retrieves a prescription record by ID
// @Summary Get prescription by ID
// @Description Retrieves a prescription record using its unique ID
// @Tags Prescriptions
// @Produce json
// @Param id path string true "Prescription ID"
// @Success 200 {object} object{id=string,examination_id=string,medications=[]object{name=string,dosage=string,start_date=string,end_date=string}} "Prescription details"
// @Failure 400 {object} object{error=string} "ID is required"
// @Failure 404 {object} object{error=string} "Prescription not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /prescriptions/{id} [get]
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

// UpdatePrescription updates an existing prescription record
// @Summary Update a prescription
// @Description Updates a prescription record with details like examination ID and a list of medications
// @Tags Prescriptions
// @Accept json
// @Produce json
// @Param request body object{id=string,examination_id=string,medications=[]object{name=string,dosage=string,start_date=string,end_date=string}} true "Updated prescription details"
// @Success 200 {object} object{id=string,examination_id=string,medications=[]object{name=string,dosage=string,start_date=string,end_date=string}} "Prescription updated successfully"
// @Failure 400 {object} object{error=string} "Invalid request, ID is required, or invalid date format"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /prescriptions [put]
func (h *RecordsHandler) UpdatePrescription(c echo.Context) error {
	var req struct {
		ID            string `json:"id"`
		ExaminationID string `json:"examination_id"`
		Medications   []struct {
			Name      string `json:"name"`
			Dosage    string `json:"dosage"`
			StartDate string `json:"start_date"` // Format: "2006-01-02"
			EndDate   string `json:"end_date"`   // Format: "2006-01-02"
		} `json:"medications"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	// Validate medications
	medications := make([]*pb.Medication, len(req.Medications))
	for i, med := range req.Medications {
		if med.Name == "" || med.Dosage == "" || med.StartDate == "" || med.EndDate == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Medication name, dosage, start_date, and end_date are required"})
		}
		_, err := time.Parse("2006-01-02", med.StartDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start_date format, must be YYYY-MM-DD"})
		}
		_, err = time.Parse("2006-01-02", med.EndDate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end_date format, must be YYYY-MM-DD"})
		}
		medications[i] = &pb.Medication{
			Name:      med.Name,
			Dosage:    med.Dosage,
			StartDate: med.StartDate,
			EndDate:   med.EndDate,
		}
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdatePrescription(ctx, &pb.UpdatePrescriptionRequest{
		Id:            req.ID,
		ExaminationId: req.ExaminationID,
		Medications:   medications,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": grpcErr.Message()})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Prescription)
}

// DeletePrescription deletes a prescription record
// @Summary Delete a prescription
// @Description Deletes a prescription record by its unique ID
// @Tags Prescriptions
// @Produce json
// @Param id path string true "Prescription ID"
// @Success 200 {object} object{success=boolean} "Prescription deleted successfully"
// @Failure 400 {object} object{error=string} "ID is required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /prescriptions/{id} [delete]
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

// ListPrescriptions lists all prescriptions for a given examination
// @Summary List prescriptions by examination
// @Description Retrieves a list of prescription records for a specific examination ID
// @Tags Prescriptions
// @Produce json
// @Param examination_id path string true "Examination ID"
// @Success 200 {array} object{id=string,examination_id=string,medications=[]object{name=string,dosage=string,start_date=string,end_date=string}} "List of prescriptions"
// @Failure 400 {object} object{error=string} "Examination ID is required"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /prescriptions/examination/{examination_id} [get]
func (h *RecordsHandler) ListPrescriptions(c echo.Context) error {
	examinationID := c.Param("examination_id")
	if examinationID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Examination ID is required"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.ListPrescriptions(ctx, &pb.ListPrescriptionsRequest{ExaminationId: examinationID})
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
