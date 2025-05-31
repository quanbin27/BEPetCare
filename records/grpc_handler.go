package main

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/quanbin27/commons/genproto/records"
)

type GRPCHandler struct {
	pb.UnimplementedPetRecordServiceServer
	service RecordsService
}

func NewGrpcHandler(grpcServer *grpc.Server, service RecordsService) {
	grpcHandler := &GRPCHandler{
		service: service,
	}
	pb.RegisterPetRecordServiceServer(grpcServer, grpcHandler)
}

// --- Pet Methods ---
func (h *GRPCHandler) CreatePet(ctx context.Context, req *pb.CreatePetRequest) (*pb.CreatePetResponse, error) {
	id, err := h.service.CreatePet(ctx, req.Name, req.Species, req.Age, req.OwnerId, req.Color, req.Weight, req.Size)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create pet: %v", err)
	}
	return &pb.CreatePetResponse{Id: id}, nil
}

func (h *GRPCHandler) GetPet(ctx context.Context, req *pb.GetPetRequest) (*pb.GetPetResponse, error) {
	pet, err := h.service.GetPet(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "pet not found: %v", err)
	}
	return &pb.GetPetResponse{
		Pet: &pb.Pet{
			Id:      pet.ID.Hex(),
			Name:    pet.Name,
			Species: pet.Species,
			Age:     pet.Age,
			OwnerId: pet.OwnerID,
			Color:   pet.Color,
			Weight:  pet.Weight,
			Size:    pet.Size,
		},
	}, nil
}

func (h *GRPCHandler) UpdatePet(ctx context.Context, req *pb.UpdatePetRequest) (*pb.UpdatePetResponse, error) {
	pet, err := h.service.UpdatePet(ctx, req.Id, req.Name, req.Species, req.Age, req.OwnerId, req.Color, req.Weight, req.Size)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update pet: %v", err)
	}
	return &pb.UpdatePetResponse{
		Pet: &pb.Pet{
			Id:      pet.ID.Hex(),
			Name:    pet.Name,
			Species: pet.Species,
			Age:     pet.Age,
			OwnerId: pet.OwnerID,
			Color:   pet.Color,
			Weight:  pet.Weight,
			Size:    pet.Size,
		},
	}, nil
}

func (h *GRPCHandler) DeletePet(ctx context.Context, req *pb.DeletePetRequest) (*pb.DeletePetResponse, error) {
	if err := h.service.DeletePet(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete pet: %v", err)
	}
	return &pb.DeletePetResponse{Success: true}, nil
}

func (h *GRPCHandler) ListPets(ctx context.Context, req *pb.ListPetsRequest) (*pb.ListPetsResponse, error) {
	pets, err := h.service.ListPets(ctx, req.OwnerId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list pets: %v", err)
	}
	resp := &pb.ListPetsResponse{}
	for _, pet := range pets {
		resp.Pets = append(resp.Pets, &pb.Pet{
			Id:      pet.ID.Hex(),
			Name:    pet.Name,
			Species: pet.Species,
			Age:     pet.Age,
			OwnerId: pet.OwnerID,
			Color:   pet.Color,
			Weight:  pet.Weight,
			Size:    pet.Size,
		})
	}
	return resp, nil
}

// --- Examination Methods ---
func (h *GRPCHandler) CreateExamination(ctx context.Context, req *pb.CreateExaminationRequest) (*pb.CreateExaminationResponse, error) {
	id, err := h.service.CreateExamination(ctx, req.PetId, req.Date, req.VetId, req.Diagnosis, req.Notes)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create examination: %v", err)
	}
	return &pb.CreateExaminationResponse{Id: id}, nil
}

func (h *GRPCHandler) GetExamination(ctx context.Context, req *pb.GetExaminationRequest) (*pb.GetExaminationResponse, error) {
	exam, err := h.service.GetExamination(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "examination not found: %v", err)
	}
	return &pb.GetExaminationResponse{
		Examination: &pb.Examination{
			Id:        exam.ID.Hex(),
			PetId:     exam.PetID,
			Date:      exam.Date.Format("2006-01-02"),
			VetId:     exam.VetID,
			Diagnosis: exam.Diagnosis,
			Notes:     exam.Notes,
		},
	}, nil
}

func (h *GRPCHandler) UpdateExamination(ctx context.Context, req *pb.UpdateExaminationRequest) (*pb.UpdateExaminationResponse, error) {
	exam, err := h.service.UpdateExamination(ctx, req.Id, req.PetId, req.Date, req.VetId, req.Diagnosis, req.Notes)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update examination: %v", err)
	}
	return &pb.UpdateExaminationResponse{
		Examination: &pb.Examination{
			Id:        exam.ID.Hex(),
			PetId:     exam.PetID,
			Date:      exam.Date.Format("2006-01-02"),
			VetId:     exam.VetID,
			Diagnosis: exam.Diagnosis,
			Notes:     exam.Notes,
		},
	}, nil
}

func (h *GRPCHandler) DeleteExamination(ctx context.Context, req *pb.DeleteExaminationRequest) (*pb.DeleteExaminationResponse, error) {
	if err := h.service.DeleteExamination(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete examination: %v", err)
	}
	return &pb.DeleteExaminationResponse{Success: true}, nil
}

func (h *GRPCHandler) ListExaminations(ctx context.Context, req *pb.ListExaminationsRequest) (*pb.ListExaminationsResponse, error) {
	exams, err := h.service.ListExaminations(ctx, req.PetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list examinations: %v", err)
	}
	resp := &pb.ListExaminationsResponse{}
	for _, exam := range exams {
		resp.Examinations = append(resp.Examinations, &pb.Examination{
			Id:        exam.ID.Hex(),
			PetId:     exam.PetID,
			Date:      exam.Date.Format("2006-01-02"),
			VetId:     exam.VetID,
			Diagnosis: exam.Diagnosis,
			Notes:     exam.Notes,
		})
	}
	return resp, nil
}

// --- Vaccination Methods ---
func (h *GRPCHandler) CreateVaccination(ctx context.Context, req *pb.CreateVaccinationRequest) (*pb.CreateVaccinationResponse, error) {
	id, err := h.service.CreateVaccination(ctx, req.PetId, req.VaccineName, req.Date, req.NextDose, req.VetId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create vaccination: %v", err)
	}
	return &pb.CreateVaccinationResponse{Id: id}, nil
}

func (h *GRPCHandler) GetVaccination(ctx context.Context, req *pb.GetVaccinationRequest) (*pb.GetVaccinationResponse, error) {
	vacc, err := h.service.GetVaccination(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "vaccination not found: %v", err)
	}
	var nextDose string
	if !vacc.NextDose.IsZero() {
		nextDose = vacc.NextDose.Format("2006-01-02")
	}
	return &pb.GetVaccinationResponse{
		Vaccination: &pb.Vaccination{
			Id:          vacc.ID.Hex(),
			PetId:       vacc.PetID,
			VaccineName: vacc.VaccineName,
			Date:        vacc.Date.Format("2006-01-02"),
			NextDose:    nextDose,
			VetId:       vacc.VetID,
		},
	}, nil
}

func (h *GRPCHandler) UpdateVaccination(ctx context.Context, req *pb.UpdateVaccinationRequest) (*pb.UpdateVaccinationResponse, error) {
	vacc, err := h.service.UpdateVaccination(ctx, req.Id, req.PetId, req.VaccineName, req.Date, req.NextDose, req.VetId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update vaccination: %v", err)
	}
	var nextDose string
	if !vacc.NextDose.IsZero() {
		nextDose = vacc.NextDose.Format("2006-01-02")
	}
	return &pb.UpdateVaccinationResponse{
		Vaccination: &pb.Vaccination{
			Id:          vacc.ID.Hex(),
			PetId:       vacc.PetID,
			VaccineName: vacc.VaccineName,
			Date:        vacc.Date.Format("2006-01-02"),
			NextDose:    nextDose,
			VetId:       vacc.VetID,
		},
	}, nil
}

func (h *GRPCHandler) DeleteVaccination(ctx context.Context, req *pb.DeleteVaccinationRequest) (*pb.DeleteVaccinationResponse, error) {
	if err := h.service.DeleteVaccination(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete vaccination: %v", err)
	}
	return &pb.DeleteVaccinationResponse{Success: true}, nil
}

func (h *GRPCHandler) ListVaccinations(ctx context.Context, req *pb.ListVaccinationsRequest) (*pb.ListVaccinationsResponse, error) {
	vaccs, err := h.service.ListVaccinations(ctx, req.PetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list vaccinations: %v", err)
	}
	resp := &pb.ListVaccinationsResponse{}
	for _, vacc := range vaccs {
		var nextDose string
		if !vacc.NextDose.IsZero() {
			nextDose = vacc.NextDose.Format("2006-01-02")
		}
		resp.Vaccinations = append(resp.Vaccinations, &pb.Vaccination{
			Id:          vacc.ID.Hex(),
			PetId:       vacc.PetID,
			VaccineName: vacc.VaccineName,
			Date:        vacc.Date.Format("2006-01-02"),
			NextDose:    nextDose,
			VetId:       vacc.VetID,
		})
	}
	return resp, nil
}

// --- Prescription Methods ---
func (h *GRPCHandler) CreatePrescription(ctx context.Context, req *pb.CreatePrescriptionRequest) (*pb.CreatePrescriptionResponse, error) {
	medications := make([]Medication, len(req.Medications))
	for i, med := range req.Medications {
		startDate, err := time.Parse("2006-01-02", med.StartDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid start date for medication at index %d: %v", i, err)
		}
		endDate, err := time.Parse("2006-01-02", med.EndDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid end date for medication at index %d: %v", i, err)
		}
		medications[i] = Medication{
			Name:      med.Name,
			Dosage:    med.Dosage,
			StartDate: startDate,
			EndDate:   endDate,
		}
	}
	id, err := h.service.CreatePrescription(ctx, req.ExaminationId, medications)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create prescription: %v", err)
	}
	return &pb.CreatePrescriptionResponse{Id: id}, nil
}

func (h *GRPCHandler) GetPrescription(ctx context.Context, req *pb.GetPrescriptionRequest) (*pb.GetPrescriptionResponse, error) {
	presc, err := h.service.GetPrescription(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "prescription not found: %v", err)
	}
	pbMedications := make([]*pb.Medication, len(presc.Medications))
	for i, med := range presc.Medications {
		pbMedications[i] = &pb.Medication{
			Name:      med.Name,
			Dosage:    med.Dosage,
			StartDate: med.StartDate.Format("2006-01-02"),
			EndDate:   med.EndDate.Format("2006-01-02"),
		}
	}
	return &pb.GetPrescriptionResponse{
		Prescription: &pb.Prescription{
			Id:            presc.ID.Hex(),
			ExaminationId: presc.ExaminationID,
			Medications:   pbMedications,
		},
	}, nil
}

func (h *GRPCHandler) UpdatePrescription(ctx context.Context, req *pb.UpdatePrescriptionRequest) (*pb.UpdatePrescriptionResponse, error) {
	medications := make([]Medication, len(req.Medications))
	for i, med := range req.Medications {
		startDate, err := time.Parse("2006-01-02", med.StartDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid start date for medication at index %d: %v", i, err)
		}
		endDate, err := time.Parse("2006-01-02", med.EndDate)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid end date for medication at index %d: %v", i, err)
		}
		medications[i] = Medication{
			Name:      med.Name,
			Dosage:    med.Dosage,
			StartDate: startDate,
			EndDate:   endDate,
		}
	}
	presc, err := h.service.UpdatePrescription(ctx, req.Id, req.ExaminationId, medications)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to update prescription: %v", err)
	}
	pbMedications := make([]*pb.Medication, len(presc.Medications))
	for i, med := range presc.Medications {
		pbMedications[i] = &pb.Medication{
			Name:      med.Name,
			Dosage:    med.Dosage,
			StartDate: med.StartDate.Format("2006-01-02"),
			EndDate:   med.EndDate.Format("2006-01-02"),
		}
	}
	return &pb.UpdatePrescriptionResponse{
		Prescription: &pb.Prescription{
			Id:            presc.ID.Hex(),
			ExaminationId: presc.ExaminationID,
			Medications:   pbMedications,
		},
	}, nil
}

func (h *GRPCHandler) DeletePrescription(ctx context.Context, req *pb.DeletePrescriptionRequest) (*pb.DeletePrescriptionResponse, error) {
	if err := h.service.DeletePrescription(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete prescription: %v", err)
	}
	return &pb.DeletePrescriptionResponse{Success: true}, nil
}

func (h *GRPCHandler) ListPrescriptions(ctx context.Context, req *pb.ListPrescriptionsRequest) (*pb.ListPrescriptionsResponse, error) {
	prescs, err := h.service.ListPrescriptions(ctx, req.ExaminationId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list prescriptions: %v", err)
	}
	resp := &pb.ListPrescriptionsResponse{}
	for _, presc := range prescs {
		pbMedications := make([]*pb.Medication, len(presc.Medications))
		for i, med := range presc.Medications {
			pbMedications[i] = &pb.Medication{
				Name:      med.Name,
				Dosage:    med.Dosage,
				StartDate: med.StartDate.Format("2006-01-02"),
				EndDate:   med.EndDate.Format("2006-01-02"),
			}
		}
		resp.Prescriptions = append(resp.Prescriptions, &pb.Prescription{
			Id:            presc.ID.Hex(),
			ExaminationId: presc.ExaminationID,
			Medications:   pbMedications,
		})
	}
	return resp, nil
}
