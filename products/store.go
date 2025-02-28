package main

import (
	"context"

	"gorm.io/gorm"
)

// Store - Triển khai ProductStore
type Store struct {
	db *gorm.DB
}

// NewStore - Hàm khởi tạo Store với database
func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// ------------------ Foods ------------------
func (s *Store) GetFoodByID(ctx context.Context, id int32) (*Food, error) {
	var food Food
	if err := s.db.WithContext(ctx).First(&food, id).Error; err != nil {
		return nil, err
	}
	return &food, nil
}

func (s *Store) ListFoods(ctx context.Context) ([]Food, error) {
	var foods []Food
	if err := s.db.WithContext(ctx).Find(&foods).Error; err != nil {
		return nil, err
	}
	return foods, nil
}

func (s *Store) CreateFood(ctx context.Context, food *Food) error {
	return s.db.WithContext(ctx).Create(food).Error
}

func (s *Store) UpdateFood(ctx context.Context, food *Food) error {
	return s.db.WithContext(ctx).Save(food).Error
}

func (s *Store) DeleteFood(ctx context.Context, id int32) error {
	return s.db.WithContext(ctx).Delete(&Food{}, id).Error
}

// ------------------ Accessories ------------------
func (s *Store) GetAccessoryByID(ctx context.Context, id int32) (*Accessory, error) {
	var accessory Accessory
	if err := s.db.WithContext(ctx).First(&accessory, id).Error; err != nil {
		return nil, err
	}
	return &accessory, nil
}

func (s *Store) ListAccessories(ctx context.Context) ([]Accessory, error) {
	var accessories []Accessory
	if err := s.db.WithContext(ctx).Find(&accessories).Error; err != nil {
		return nil, err
	}
	return accessories, nil
}

func (s *Store) CreateAccessory(ctx context.Context, accessory *Accessory) error {
	return s.db.WithContext(ctx).Create(accessory).Error
}

func (s *Store) UpdateAccessory(ctx context.Context, accessory *Accessory) error {
	return s.db.WithContext(ctx).Save(accessory).Error
}

func (s *Store) DeleteAccessory(ctx context.Context, id int32) error {
	return s.db.WithContext(ctx).Delete(&Accessory{}, id).Error
}

// ------------------ Medicines ------------------
func (s *Store) GetMedicineByID(ctx context.Context, id int32) (*Medicine, error) {
	var medicine Medicine
	if err := s.db.WithContext(ctx).First(&medicine, id).Error; err != nil {
		return nil, err
	}
	return &medicine, nil
}

func (s *Store) ListMedicines(ctx context.Context) ([]Medicine, error) {
	var medicines []Medicine
	if err := s.db.WithContext(ctx).Find(&medicines).Error; err != nil {
		return nil, err
	}
	return medicines, nil
}

func (s *Store) CreateMedicine(ctx context.Context, medicine *Medicine) error {
	return s.db.WithContext(ctx).Create(medicine).Error
}

func (s *Store) UpdateMedicine(ctx context.Context, medicine *Medicine) error {
	return s.db.WithContext(ctx).Save(medicine).Error
}

func (s *Store) DeleteMedicine(ctx context.Context, id int32) error {
	return s.db.WithContext(ctx).Delete(&Medicine{}, id).Error
}

// ------------------ Branches ------------------
func (s *Store) GetBranchByID(ctx context.Context, id int32) (*Branch, error) {
	var branch Branch
	if err := s.db.WithContext(ctx).First(&branch, id).Error; err != nil {
		return nil, err
	}
	return &branch, nil
}

func (s *Store) ListBranches(ctx context.Context) ([]Branch, error) {
	var branches []Branch
	if err := s.db.WithContext(ctx).Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}

// ------------------ Branch Inventory ------------------
func (s *Store) GetBranchInventory(ctx context.Context, branchID int32) ([]BranchProduct, error) {
	var inventory []BranchProduct
	if err := s.db.WithContext(ctx).Where("branch_id = ?", branchID).Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}

func (s *Store) UpdateBranchInventory(ctx context.Context, branchID int32, productID int32, productType string, stockQuantity int32) error {
	return s.db.WithContext(ctx).Model(&BranchProduct{}).
		Where("branch_id = ? AND product_id = ? AND product_type = ?", branchID, productID, productType).
		Update("stock_quantity", stockQuantity).Error
}
