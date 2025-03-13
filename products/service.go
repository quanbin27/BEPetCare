package main

import (
	"context"
	"errors"
)

type ProductServiceImpl struct {
	store ProductStore
}

func NewProductService(store ProductStore) ProductService {
	return &ProductServiceImpl{store: store}
}

// Thực phẩm
func (s *ProductServiceImpl) GetFoodByID(ctx context.Context, id int32) (*Food, error) {
	return s.store.GetFoodByID(ctx, id)
}

func (s *ProductServiceImpl) ListFoods(ctx context.Context) ([]Food, error) {
	return s.store.ListFoods(ctx)
}

func (s *ProductServiceImpl) CreateFood(ctx context.Context, name, description string, price float32) (string, error) {
	food := &Food{
		Name:        name,
		Description: description,
		Price:       price,
	}
	if err := s.store.CreateFood(ctx, food); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *ProductServiceImpl) UpdateFood(ctx context.Context, id int32, name, description string, price float32) (string, error) {
	food, err := s.store.GetFoodByID(ctx, id)
	if err != nil {
		return "Failed", errors.New("food not found")
	}
	food.Name = name
	food.Description = description
	food.Price = price
	if err := s.store.UpdateFood(ctx, food); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *ProductServiceImpl) DeleteFood(ctx context.Context, id int32) (string, error) {
	if err := s.store.DeleteFood(ctx, id); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

// Phụ kiện
func (s *ProductServiceImpl) GetAccessoryByID(ctx context.Context, id int32) (*Accessory, error) {
	return s.store.GetAccessoryByID(ctx, id)
}

func (s *ProductServiceImpl) ListAccessories(ctx context.Context) ([]Accessory, error) {
	return s.store.ListAccessories(ctx)
}

func (s *ProductServiceImpl) CreateAccessory(ctx context.Context, name, description string, price float32) (string, error) {
	accessory := &Accessory{
		Name:        name,
		Description: description,
		Price:       price,
	}
	if err := s.store.CreateAccessory(ctx, accessory); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *ProductServiceImpl) UpdateAccessory(ctx context.Context, id int32, name, description string, price float32) (string, error) {
	accessory, err := s.store.GetAccessoryByID(ctx, id)
	if err != nil {
		return "Failed", errors.New("accessory not found")
	}
	accessory.Name = name
	accessory.Description = description
	accessory.Price = price
	if err := s.store.UpdateAccessory(ctx, accessory); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *ProductServiceImpl) DeleteAccessory(ctx context.Context, id int32) (string, error) {
	if err := s.store.DeleteAccessory(ctx, id); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

// Thuốc
func (s *ProductServiceImpl) GetMedicineByID(ctx context.Context, id int32) (*Medicine, error) {
	return s.store.GetMedicineByID(ctx, id)
}

func (s *ProductServiceImpl) ListMedicines(ctx context.Context) ([]Medicine, error) {
	return s.store.ListMedicines(ctx)
}

func (s *ProductServiceImpl) CreateMedicine(ctx context.Context, name, description string, price float32) (string, error) {
	medicine := &Medicine{
		Name:        name,
		Description: description,
		Price:       price,
	}
	if err := s.store.CreateMedicine(ctx, medicine); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *ProductServiceImpl) UpdateMedicine(ctx context.Context, id int32, name, description string, price float32) (string, error) {
	medicine, err := s.store.GetMedicineByID(ctx, id)
	if err != nil {
		return "Failed", errors.New("medicine not found")
	}
	medicine.Name = name
	medicine.Description = description
	medicine.Price = price
	if err := s.store.UpdateMedicine(ctx, medicine); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

func (s *ProductServiceImpl) DeleteMedicine(ctx context.Context, id int32) (string, error) {
	if err := s.store.DeleteMedicine(ctx, id); err != nil {
		return "Failed", err
	}
	return "Success", nil
}

// Chi nhánh
func (s *ProductServiceImpl) GetBranchByID(ctx context.Context, id int32) (*Branch, error) {
	return s.store.GetBranchByID(ctx, id)
}

func (s *ProductServiceImpl) ListBranches(ctx context.Context) ([]Branch, error) {
	return s.store.ListBranches(ctx)
}
func (s *ProductServiceImpl) ListAttachableProducts(ctx context.Context) ([]GeneralProduct, error) {
	return s.store.ListAttachableProducts(ctx)
}

// Tồn kho
func (s *ProductServiceImpl) GetBranchInventory(ctx context.Context, branchID int32) ([]BranchProduct, error) {
	return s.store.GetBranchInventory(ctx, branchID)
}

func (s *ProductServiceImpl) UpdateBranchInventory(ctx context.Context, branchID, productID int32, productType string, stockQuantity int32) (string, error) {
	if err := s.store.UpdateBranchInventory(ctx, branchID, productID, productType, stockQuantity); err != nil {
		return "Failed", err
	}
	return "Success", nil
}
