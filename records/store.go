package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoStore struct {
	client     *mongo.Client
	database   *mongo.Database
	pets       *mongo.Collection
	exams      *mongo.Collection
	vaccs      *mongo.Collection
	prescripts *mongo.Collection
}

func NewMongoStore(dsn string) (RecordsStore, error) {
	// Kết nối tới MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	// Kiểm tra kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	// Chọn database
	db := client.Database("petrecord_db")

	// Khởi tạo collections
	store := &MongoStore{
		client:     client,
		database:   db,
		pets:       db.Collection("pets"),
		exams:      db.Collection("examinations"),
		vaccs:      db.Collection("vaccinations"),
		prescripts: db.Collection("prescriptions"),
	}

	// Tạo index cho các collection (tùy chọn)
	if err := store.initIndexes(ctx); err != nil {
		return nil, err
	}

	return store, nil
}

// initIndexes tạo các index để tối ưu hóa truy vấn
func (s *MongoStore) initIndexes(ctx context.Context) error {
	// Index cho pets.owner_id
	_, err := s.pets.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "owner_id", Value: 1}},
	})
	if err != nil {
		log.Printf("Failed to create index for pets.owner_id: %v", err)
		return err
	}

	// Index cho examinations.pet_id
	_, err = s.exams.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "pet_id", Value: 1}},
	})
	if err != nil {
		log.Printf("Failed to create index for examinations.pet_id: %v", err)
		return err
	}

	// Index cho vaccinations.pet_id
	_, err = s.vaccs.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "pet_id", Value: 1}},
	})
	if err != nil {
		log.Printf("Failed to create index for vaccinations.pet_id: %v", err)
		return err
	}

	// Index cho prescriptions.pet_id
	_, err = s.prescripts.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "pet_id", Value: 1}},
	})
	if err != nil {
		log.Printf("Failed to create index for prescriptions.pet_id: %v", err)
		return err
	}

	return nil
}

// --- Pet Methods ---
func (s *MongoStore) CreatePet(ctx context.Context, pet *Pet) (string, error) {
	pet.ID = primitive.NewObjectID()
	_, err := s.pets.InsertOne(ctx, pet)
	if err != nil {
		return "", err
	}
	return pet.ID.Hex(), nil
}

func (s *MongoStore) GetPet(ctx context.Context, id string) (*Pet, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var pet Pet
	err = s.pets.FindOne(ctx, bson.M{"_id": objID}).Decode(&pet)
	if err != nil {
		return nil, err
	}
	return &pet, nil
}

func (s *MongoStore) UpdatePet(ctx context.Context, pet *Pet) error {
	_, err := s.pets.ReplaceOne(ctx, bson.M{"_id": pet.ID}, pet)
	return err
}

func (s *MongoStore) DeletePet(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.pets.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (s *MongoStore) ListPets(ctx context.Context, ownerID string) ([]*Pet, error) {
	filter := bson.M{}
	if ownerID != "" {
		filter["owner_id"] = ownerID
	}
	cursor, err := s.pets.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pets []*Pet
	for cursor.Next(ctx) {
		var pet Pet
		if err := cursor.Decode(&pet); err != nil {
			return nil, err
		}
		pets = append(pets, &pet)
	}
	return pets, nil
}

// --- Examination Methods ---
func (s *MongoStore) CreateExamination(ctx context.Context, exam *Examination) (string, error) {
	exam.ID = primitive.NewObjectID()
	_, err := s.exams.InsertOne(ctx, exam)
	if err != nil {
		return "", err
	}
	return exam.ID.Hex(), nil
}

func (s *MongoStore) GetExamination(ctx context.Context, id string) (*Examination, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var exam Examination
	err = s.exams.FindOne(ctx, bson.M{"_id": objID}).Decode(&exam)
	if err != nil {
		return nil, err
	}
	return &exam, nil
}

func (s *MongoStore) UpdateExamination(ctx context.Context, exam *Examination) error {
	_, err := s.exams.ReplaceOne(ctx, bson.M{"_id": exam.ID}, exam)
	return err
}

func (s *MongoStore) DeleteExamination(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.exams.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (s *MongoStore) ListExaminations(ctx context.Context, petID string) ([]*Examination, error) {
	cursor, err := s.exams.Find(ctx, bson.M{"pet_id": petID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var exams []*Examination
	for cursor.Next(ctx) {
		var exam Examination
		if err := cursor.Decode(&exam); err != nil {
			return nil, err
		}
		exams = append(exams, &exam)
	}
	return exams, nil
}

// --- Vaccination Methods ---
func (s *MongoStore) CreateVaccination(ctx context.Context, vacc *Vaccination) (string, error) {
	vacc.ID = primitive.NewObjectID()
	_, err := s.vaccs.InsertOne(ctx, vacc)
	if err != nil {
		return "", err
	}
	return vacc.ID.Hex(), nil
}

func (s *MongoStore) GetVaccination(ctx context.Context, id string) (*Vaccination, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var vacc Vaccination
	err = s.vaccs.FindOne(ctx, bson.M{"_id": objID}).Decode(&vacc)
	if err != nil {
		return nil, err
	}
	return &vacc, nil
}

func (s *MongoStore) UpdateVaccination(ctx context.Context, vacc *Vaccination) error {
	_, err := s.vaccs.ReplaceOne(ctx, bson.M{"_id": vacc.ID}, vacc)
	return err
}

func (s *MongoStore) DeleteVaccination(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.vaccs.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (s *MongoStore) ListVaccinations(ctx context.Context, petID string) ([]*Vaccination, error) {
	cursor, err := s.vaccs.Find(ctx, bson.M{"pet_id": petID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var vaccs []*Vaccination
	for cursor.Next(ctx) {
		var vacc Vaccination
		if err := cursor.Decode(&vacc); err != nil {
			return nil, err
		}
		vaccs = append(vaccs, &vacc)
	}
	return vaccs, nil
}

// --- Prescription Methods ---
func (s *MongoStore) CreatePrescription(ctx context.Context, presc *Prescription) (string, error) {
	presc.ID = primitive.NewObjectID()
	_, err := s.prescripts.InsertOne(ctx, presc)
	if err != nil {
		return "", err
	}
	return presc.ID.Hex(), nil
}

func (s *MongoStore) GetPrescription(ctx context.Context, id string) (*Prescription, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var presc Prescription
	err = s.prescripts.FindOne(ctx, bson.M{"_id": objID}).Decode(&presc)
	if err != nil {
		return nil, err
	}
	return &presc, nil
}

func (s *MongoStore) UpdatePrescription(ctx context.Context, presc *Prescription) error {
	_, err := s.prescripts.ReplaceOne(ctx, bson.M{"_id": presc.ID}, presc)
	return err
}

func (s *MongoStore) DeletePrescription(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.prescripts.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (s *MongoStore) ListPrescriptions(ctx context.Context, petID string) ([]*Prescription, error) {
	cursor, err := s.prescripts.Find(ctx, bson.M{"pet_id": petID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var prescs []*Prescription
	for cursor.Next(ctx) {
		var presc Prescription
		if err := cursor.Decode(&presc); err != nil {
			return nil, err
		}
		prescs = append(prescs, &presc)
	}
	return prescs, nil
}
