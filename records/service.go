package main

import (
	"context"
	"errors"
	"time"
)

type PetRecordService struct {
	store RecordsStore
}

func NewPetRecordService(store RecordsStore) RecordsService {
	return &PetRecordService{store: store}
}

// --- Pet Methods ---
func (s *PetRecordService) CreatePet(ctx context.Context, name, species string, age int32, ownerID, color string, weight float32, size string) (string, error) {
	pet := &Pet{
		Name:    name,
		Species: species,
		Age:     age,
		OwnerID: ownerID,
		Color:   color,
		Weight:  weight,
		Size:    size,
	}
	return s.store.CreatePet(ctx, pet)
}

func (s *PetRecordService) GetPet(ctx context.Context, id string) (*Pet, error) {
	return s.store.GetPet(ctx, id)
}

func (s *PetRecordService) UpdatePet(ctx context.Context, id, name, species string, age int32, ownerID, color string, weight float32, size string) (*Pet, error) {
	pet, err := s.store.GetPet(ctx, id)
	if err != nil {
		return nil, err
	}
	pet.Name = name
	pet.Species = species
	pet.Age = age
	pet.OwnerID = ownerID
	pet.Color = color
	pet.Weight = weight
	pet.Size = size
	if err := s.store.UpdatePet(ctx, pet); err != nil {
		return nil, err
	}
	return pet, nil
}

func (s *PetRecordService) DeletePet(ctx context.Context, id string) error {
	return s.store.DeletePet(ctx, id)
}

func (s *PetRecordService) ListPets(ctx context.Context, ownerID string) ([]*Pet, error) {
	return s.store.ListPets(ctx, ownerID)
}

// --- Examination Methods ---
func (s *PetRecordService) CreateExamination(ctx context.Context, petID, dateStr, vetID, diagnosis, notes string) (string, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", errors.New("invalid date format")
	}
	exam := &Examination{
		PetID:     petID,
		Date:      date,
		VetID:     vetID,
		Diagnosis: diagnosis,
		Notes:     notes,
	}
	return s.store.CreateExamination(ctx, exam)
}

func (s *PetRecordService) GetExamination(ctx context.Context, id string) (*Examination, error) {
	return s.store.GetExamination(ctx, id)
}

func (s *PetRecordService) UpdateExamination(ctx context.Context, id, petID, dateStr, vetID, diagnosis, notes string) (*Examination, error) {
	exam, err := s.store.GetExamination(ctx, id)
	if err != nil {
		return nil, err
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format")
	}
	exam.PetID = petID
	exam.Date = date
	exam.VetID = vetID
	exam.Diagnosis = diagnosis
	exam.Notes = notes
	if err := s.store.UpdateExamination(ctx, exam); err != nil {
		return nil, err
	}
	return exam, nil
}

func (s *PetRecordService) DeleteExamination(ctx context.Context, id string) error {
	return s.store.DeleteExamination(ctx, id)
}

func (s *PetRecordService) ListExaminations(ctx context.Context, petID string) ([]*Examination, error) {
	return s.store.ListExaminations(ctx, petID)
}

// --- Vaccination Methods ---
func (s *PetRecordService) CreateVaccination(ctx context.Context, petID, vaccineName, dateStr, vetID string) (string, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", errors.New("invalid date format")
	}
	vacc := &Vaccination{
		PetID:       petID,
		VaccineName: vaccineName,
		Date:        date,
		VetID:       vetID,
	}
	return s.store.CreateVaccination(ctx, vacc)
}

func (s *PetRecordService) GetVaccination(ctx context.Context, id string) (*Vaccination, error) {
	return s.store.GetVaccination(ctx, id)
}

func (s *PetRecordService) UpdateVaccination(ctx context.Context, id, petID, vaccineName, dateStr, vetID string) (*Vaccination, error) {
	vacc, err := s.store.GetVaccination(ctx, id)
	if err != nil {
		return nil, err
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format")
	}
	vacc.PetID = petID
	vacc.VaccineName = vaccineName
	vacc.Date = date
	vacc.VetID = vetID
	if err := s.store.UpdateVaccination(ctx, vacc); err != nil {
		return nil, err
	}
	return vacc, nil
}

func (s *PetRecordService) DeleteVaccination(ctx context.Context, id string) error {
	return s.store.DeleteVaccination(ctx, id)
}

func (s *PetRecordService) ListVaccinations(ctx context.Context, petID string) ([]*Vaccination, error) {
	return s.store.ListVaccinations(ctx, petID)
}

// --- Prescription Methods ---
func (s *PetRecordService) CreatePrescription(ctx context.Context, petID, medication, dosage, startDateStr, endDateStr string) (string, error) {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return "", errors.New("invalid start date format")
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return "", errors.New("invalid end date format")
	}
	presc := &Prescription{
		PetID:      petID,
		Medication: medication,
		Dosage:     dosage,
		StartDate:  startDate,
		EndDate:    endDate,
	}
	return s.store.CreatePrescription(ctx, presc)
}

func (s *PetRecordService) GetPrescription(ctx context.Context, id string) (*Prescription, error) {
	return s.store.GetPrescription(ctx, id)
}

func (s *PetRecordService) UpdatePrescription(ctx context.Context, id, petID, medication, dosage, startDateStr, endDateStr string) (*Prescription, error) {
	presc, err := s.store.GetPrescription(ctx, id)
	if err != nil {
		return nil, err
	}
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}
	presc.PetID = petID
	presc.Medication = medication
	presc.Dosage = dosage
	presc.StartDate = startDate
	presc.EndDate = endDate
	if err := s.store.UpdatePrescription(ctx, presc); err != nil {
		return nil, err
	}
	return presc, nil
}

func (s *PetRecordService) DeletePrescription(ctx context.Context, id string) error {
	return s.store.DeletePrescription(ctx, id)
}

func (s *PetRecordService) ListPrescriptions(ctx context.Context, petID string) ([]*Prescription, error) {
	return s.store.ListPrescriptions(ctx, petID)
}
