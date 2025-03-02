package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Pet represents a pet record
type Pet struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Species string             `bson:"species"`
	Age     int32              `bson:"age"`
	OwnerID string             `bson:"owner_id"`
	Color   string             `bson:"color"`
	Weight  float32            `bson:"weight"`
	Size    string             `bson:"size"`
}

// Examination represents an examination record
type Examination struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	PetID     string             `bson:"pet_id"`
	Date      time.Time          `bson:"date"`
	VetID     string             `bson:"vet_id"`
	Diagnosis string             `bson:"diagnosis"`
	Notes     string             `bson:"notes"`
}

// Vaccination represents a vaccination record
type Vaccination struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PetID       string             `bson:"pet_id"`
	VaccineName string             `bson:"vaccine_name"`
	Date        time.Time          `bson:"date"`
	VetID       string             `bson:"vet_id"`
}

// Prescription represents a prescription record
type Prescription struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	PetID      string             `bson:"pet_id"`
	Medication string             `bson:"medication"`
	Dosage     string             `bson:"dosage"`
	StartDate  time.Time          `bson:"start_date"`
	EndDate    time.Time          `bson:"end_date"`
}

// Store defines the interface for data storage operations
type RecordsStore interface {
	// Pet methods
	CreatePet(ctx context.Context, pet *Pet) (string, error)
	GetPet(ctx context.Context, id string) (*Pet, error)
	UpdatePet(ctx context.Context, pet *Pet) error
	DeletePet(ctx context.Context, id string) error
	ListPets(ctx context.Context, ownerID string) ([]*Pet, error)

	// Examination methods
	CreateExamination(ctx context.Context, exam *Examination) (string, error)
	GetExamination(ctx context.Context, id string) (*Examination, error)
	UpdateExamination(ctx context.Context, exam *Examination) error
	DeleteExamination(ctx context.Context, id string) error
	ListExaminations(ctx context.Context, petID string) ([]*Examination, error)

	// Vaccination methods
	CreateVaccination(ctx context.Context, vacc *Vaccination) (string, error)
	GetVaccination(ctx context.Context, id string) (*Vaccination, error)
	UpdateVaccination(ctx context.Context, vacc *Vaccination) error
	DeleteVaccination(ctx context.Context, id string) error
	ListVaccinations(ctx context.Context, petID string) ([]*Vaccination, error)

	// Prescription methods
	CreatePrescription(ctx context.Context, presc *Prescription) (string, error)
	GetPrescription(ctx context.Context, id string) (*Prescription, error)
	UpdatePrescription(ctx context.Context, presc *Prescription) error
	DeletePrescription(ctx context.Context, id string) error
	ListPrescriptions(ctx context.Context, petID string) ([]*Prescription, error)
}

// Service defines the interface for business logic operations
type RecordsService interface {
	// Pet methods
	CreatePet(ctx context.Context, name, species string, age int32, ownerID, color string, weight float32, size string) (string, error)
	GetPet(ctx context.Context, id string) (*Pet, error)
	UpdatePet(ctx context.Context, id, name, species string, age int32, ownerID, color string, weight float32, size string) (*Pet, error)
	DeletePet(ctx context.Context, id string) error
	ListPets(ctx context.Context, ownerID string) ([]*Pet, error)

	// Examination methods
	CreateExamination(ctx context.Context, petID, dateStr, vetID, diagnosis, notes string) (string, error)
	GetExamination(ctx context.Context, id string) (*Examination, error)
	UpdateExamination(ctx context.Context, id, petID, dateStr, vetID, diagnosis, notes string) (*Examination, error)
	DeleteExamination(ctx context.Context, id string) error
	ListExaminations(ctx context.Context, petID string) ([]*Examination, error)

	// Vaccination methods
	CreateVaccination(ctx context.Context, petID, vaccineName, dateStr, vetID string) (string, error)
	GetVaccination(ctx context.Context, id string) (*Vaccination, error)
	UpdateVaccination(ctx context.Context, id, petID, vaccineName, dateStr, vetID string) (*Vaccination, error)
	DeleteVaccination(ctx context.Context, id string) error
	ListVaccinations(ctx context.Context, petID string) ([]*Vaccination, error)

	// Prescription methods
	CreatePrescription(ctx context.Context, petID, medication, dosage, startDateStr, endDateStr string) (string, error)
	GetPrescription(ctx context.Context, id string) (*Prescription, error)
	UpdatePrescription(ctx context.Context, id, petID, medication, dosage, startDateStr, endDateStr string) (*Prescription, error)
	DeletePrescription(ctx context.Context, id string) error
	ListPrescriptions(ctx context.Context, petID string) ([]*Prescription, error)
}
