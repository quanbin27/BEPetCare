package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Pet, Examination, and Vaccination structs remain unchanged
type Pet struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Species      string             `bson:"species"`
	Dob          string             `bson:"dob"`
	OwnerID      string             `bson:"owner_id"`
	Color        string             `bson:"color"`
	Weight       float32            `bson:"weight"`
	identityMark string             `bson:"identity_mark"`
}

type Examination struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	PetID     string             `bson:"pet_id"`
	Date      time.Time          `bson:"date"`
	VetID     string             `bson:"vet_id"`
	Diagnosis string             `bson:"diagnosis"`
	Notes     string             `bson:"notes"`
	VetName   string             `bson:"vet_name"`
}

type Vaccination struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PetID       string             `bson:"pet_id"`
	VaccineName string             `bson:"vaccine_name"`
	Date        time.Time          `bson:"date"`
	NextDose    time.Time          `bson:"next_dose,omitempty"`
	VetID       string             `bson:"vet_id"`
	VetName     string             `bson:"vet_name"`
}

// Medication represents a single medication in a prescription
type Medication struct {
	MedicineID string    `bson:"medicine_id"`
	Name       string    `bson:"name"`
	Dosage     string    `bson:"dosage"`
	StartDate  time.Time `bson:"start_date"` // Moved to Medication
	EndDate    time.Time `bson:"end_date"`   // Moved to Medication
}

// Prescription represents a prescription record, linked to an examination
type Prescription struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	ExaminationID string             `bson:"examination_id"`
	Medications   []Medication       `bson:"medications"`
}

// RecordsStore defines the interface for data storage operations
type RecordsStore interface {
	// Pet, Examination, and Vaccination methods remain unchanged
	CreatePet(ctx context.Context, pet *Pet) (string, error)
	GetPet(ctx context.Context, id string) (*Pet, error)
	UpdatePet(ctx context.Context, pet *Pet) error
	DeletePet(ctx context.Context, id string) error
	ListPets(ctx context.Context, ownerID string) ([]*Pet, error)

	CreateExamination(ctx context.Context, exam *Examination) (string, error)
	GetExamination(ctx context.Context, id string) (*Examination, error)
	UpdateExamination(ctx context.Context, exam *Examination) error
	DeleteExamination(ctx context.Context, id string) error
	ListExaminations(ctx context.Context, petID string) ([]*Examination, error)

	CreateVaccination(ctx context.Context, vacc *Vaccination) (string, error)
	GetVaccination(ctx context.Context, id string) (*Vaccination, error)
	UpdateVaccination(ctx context.Context, vacc *Vaccination) error
	DeleteVaccination(ctx context.Context, id string) error
	ListVaccinations(ctx context.Context, petID string) ([]*Vaccination, error)

	// Prescription methods (updated to reflect Medications containing StartDate and EndDate)
	CreatePrescription(ctx context.Context, presc *Prescription) (string, error)
	GetPrescription(ctx context.Context, id string) (*Prescription, error)
	UpdatePrescription(ctx context.Context, presc *Prescription) error
	DeletePrescription(ctx context.Context, id string) error
	GetPrescriptionByExaminationID(ctx context.Context, examinationID string) (*Prescription, error)
}

// RecordsService defines the interface for business logic operations
type RecordsService interface {
	// Pet, Examination, and Vaccination methods remain unchanged
	CreatePet(ctx context.Context, name, species string, dob string, ownerID, color string, weight float32, identityMark string) (string, error)
	GetPet(ctx context.Context, id string) (*Pet, error)
	UpdatePet(ctx context.Context, id, name, species string, dob string, ownerID, color string, weight float32, identityMark string) (*Pet, error)
	DeletePet(ctx context.Context, id string) error
	ListPets(ctx context.Context, ownerID string) ([]*Pet, error)

	CreateExamination(ctx context.Context, petID, dateStr, vetID, diagnosis, notes, vetName string) (string, error)
	GetExamination(ctx context.Context, id string) (*Examination, error)
	UpdateExamination(ctx context.Context, id, petID, dateStr, vetID, diagnosis, notes, vetName string) (*Examination, error)
	DeleteExamination(ctx context.Context, id string) error
	ListExaminations(ctx context.Context, petID string) ([]*Examination, error)

	CreateVaccination(ctx context.Context, petID, vaccineName, dateStr, nextDoseStr, vetID, vetName string) (string, error)
	GetVaccination(ctx context.Context, id string) (*Vaccination, error)
	UpdateVaccination(ctx context.Context, id, petID, vaccineName, dateStr, nextDoseStr, vetID, vetName string) (*Vaccination, error)
	DeleteVaccination(ctx context.Context, id string) error
	ListVaccinations(ctx context.Context, petID string) ([]*Vaccination, error)

	// Prescription methods (updated to handle Medications with StartDate and EndDate)
	CreatePrescription(ctx context.Context, examinationID string, medications []Medication) (string, error)
	GetPrescription(ctx context.Context, id string) (*Prescription, error)
	UpdatePrescription(ctx context.Context, id, examinationID string, medications []Medication) (*Prescription, error)
	DeletePrescription(ctx context.Context, id string) error
	GetPrescriptionByExaminationID(ctx context.Context, examinationID string) (*Prescription, error)
}
