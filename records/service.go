package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type PetRecordService struct {
	store RecordsStore
}

func NewPetRecordService(store RecordsStore) RecordsService {
	return &PetRecordService{store: store}
}

// --- Pet Methods ---
func (s *PetRecordService) CreatePet(ctx context.Context, name, species string, dob string, ownerID, color string, weight float32, identityMark string) (string, error) {
	if name == "" || species == "" || ownerID == "" {
		return "", errors.New("name, species, and ownerID are required")
	}
	pet := &Pet{
		Name:         name,
		Species:      species,
		Dob:          dob,
		OwnerID:      ownerID,
		Color:        color,
		Weight:       weight,
		identityMark: identityMark,
	}
	return s.store.CreatePet(ctx, pet)
}

func (s *PetRecordService) GetPet(ctx context.Context, id string) (*Pet, error) {
	return s.store.GetPet(ctx, id)
}

func (s *PetRecordService) UpdatePet(ctx context.Context, id, name, species string, dob string, ownerID, color string, weight float32, identityMark string) (*Pet, error) {
	if name == "" || species == "" || ownerID == "" {
		return nil, errors.New("name, species, and ownerID are required")
	}
	pet, err := s.store.GetPet(ctx, id)
	if err != nil {
		return nil, err
	}
	pet.Name = name
	pet.Species = species
	pet.Dob = dob
	pet.OwnerID = ownerID
	pet.Color = color
	pet.Weight = weight
	pet.identityMark = identityMark
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
func (s *PetRecordService) CreateExamination(ctx context.Context, petID, dateStr, vetID, diagnosis, notes, vetName string) (string, error) {
	if petID == "" || vetID == "" {
		return "", errors.New("petID and vetID are required")
	}
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
		VetName:   vetName,
	}
	return s.store.CreateExamination(ctx, exam)
}

func (s *PetRecordService) GetExamination(ctx context.Context, id string) (*Examination, error) {
	return s.store.GetExamination(ctx, id)
}

func (s *PetRecordService) UpdateExamination(ctx context.Context, id, petID, dateStr, vetID, diagnosis, notes, vetName string) (*Examination, error) {
	if petID == "" || vetID == "" {
		return nil, errors.New("petID and vetID are required")
	}
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
	exam.VetName = vetName
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
func (s *PetRecordService) CreateVaccination(ctx context.Context, petID, vaccineName, dateStr, nextDoseStr, vetID, vetName string) (string, error) {
	if petID == "" || vaccineName == "" || vetID == "" {
		return "", errors.New("petID, vaccineName, and vetID are required")
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", errors.New("invalid date format")
	}
	var nextDose time.Time
	if nextDoseStr != "" {
		nextDose, err = time.Parse("2006-01-02", nextDoseStr)
		if err != nil {
			return "", errors.New("invalid next dose date format")
		}
	}
	vacc := &Vaccination{
		PetID:       petID,
		VaccineName: vaccineName,
		Date:        date,
		NextDose:    nextDose,
		VetID:       vetID,
		VetName:     vetName,
	}
	return s.store.CreateVaccination(ctx, vacc)
}

func (s *PetRecordService) GetVaccination(ctx context.Context, id string) (*Vaccination, error) {
	return s.store.GetVaccination(ctx, id)
}

func (s *PetRecordService) UpdateVaccination(ctx context.Context, id, petID, vaccineName, dateStr, nextDoseStr, vetID, vetName string) (*Vaccination, error) {
	if petID == "" || vaccineName == "" || vetID == "" {
		return nil, errors.New("petID, vaccineName, and vetID are required")
	}
	vacc, err := s.store.GetVaccination(ctx, id)
	if err != nil {
		return nil, err
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format")
	}
	var nextDose time.Time
	if nextDoseStr != "" {
		nextDose, err = time.Parse("2006-01-02", nextDoseStr)
		if err != nil {
			return nil, errors.New("invalid next dose date format")
		}
	}
	vacc.PetID = petID
	vacc.VaccineName = vaccineName
	vacc.Date = date
	vacc.NextDose = nextDose
	vacc.VetID = vetID
	vacc.VetName = vetName
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
func (s *PetRecordService) CreatePrescription(ctx context.Context, examinationID string, medications []Medication) (string, error) {
	if examinationID == "" {
		return "", errors.New("examinationID is required")
	}
	if len(medications) == 0 {
		return "", errors.New("medications list cannot be empty")
	}
	for i, med := range medications {
		if med.Name == "" || med.Dosage == "" {
			return "", fmt.Errorf("medication at index %d has empty name or dosage", i)
		}
		if med.StartDate.IsZero() || med.EndDate.IsZero() {
			return "", fmt.Errorf("medication at index %d has invalid start or end date", i)
		}
		if med.EndDate.Before(med.StartDate) {
			return "", fmt.Errorf("medication at index %d has end date before start date", i)
		}
	}
	presc := &Prescription{
		ExaminationID: examinationID,
		Medications:   medications,
	}
	return s.store.CreatePrescription(ctx, presc)
}

func (s *PetRecordService) GetPrescription(ctx context.Context, id string) (*Prescription, error) {
	return s.store.GetPrescription(ctx, id)
}

func (s *PetRecordService) UpdatePrescription(ctx context.Context, id, examinationID string, medications []Medication) (*Prescription, error) {
	if examinationID == "" {
		return nil, errors.New("examinationID is required")
	}
	if len(medications) == 0 {
		return nil, errors.New("medications list cannot be empty")
	}
	for i, med := range medications {
		if med.Name == "" || med.Dosage == "" {
			return nil, fmt.Errorf("medication at index %d has empty name or dosage", i)
		}
		if med.StartDate.IsZero() || med.EndDate.IsZero() {
			return nil, fmt.Errorf("medication at index %d has invalid start or end date", i)
		}
		if med.EndDate.Before(med.StartDate) {
			return nil, fmt.Errorf("medication at index %d has end date before start date", i)
		}
	}
	presc, err := s.store.GetPrescription(ctx, id)
	if err != nil {
		return nil, err
	}
	presc.ExaminationID = examinationID
	presc.Medications = medications
	if err := s.store.UpdatePrescription(ctx, presc); err != nil {
		return nil, err
	}
	return presc, nil
}

func (s *PetRecordService) DeletePrescription(ctx context.Context, id string) error {
	return s.store.DeletePrescription(ctx, id)
}

func (s *PetRecordService) GetPrescriptionByExaminationID(ctx context.Context, examinationID string) (*Prescription, error) {
	return s.store.GetPrescriptionByExaminationID(ctx, examinationID)
}
