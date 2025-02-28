package main

import (
	"context"

	pb "github.com/quanbin27/commons/genproto/products"
)

type Service struct {
	store ProductStore
}

func NewService(store ProductStore) *Service {
	return &Service{store: store}
}

// ------------------ Foods ------------------
func (s *Service) GetFoodByID(ctx context.Context, req *pb.GetFoodRequest) (*pb.Food, error) {
	food, err := s.store.GetFoodByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Food{Id: food.ID, Name: food.Name, Description: food.Description, Price: food.Price}, nil
}

func (s *Service) ListFoods(ctx context.Context, req *pb.ListFoodRequest) (*pb.ListFoodResponse, error) {
	foods, err := s.store.ListFoods(ctx)
	if err != nil {
		return nil, err
	}
	var pbFoods []*pb.Food
	for _, food := range foods {
		pbFoods = append(pbFoods, &pb.Food{Id: food.ID, Name: food.Name, Description: food.Description, Price: food.Price})
	}
	return &pb.ListFoodResponse{Foods: pbFoods}, nil
}

func (s *Service) CreateFood(ctx context.Context, req *pb.CreateFoodRequest) (*pb.CreateFoodResponse, error) {
	food := &Food{Name: req.Name, Description: req.Description, Price: req.Price}
	if err := s.store.CreateFood(ctx, food); err != nil {
		return &pb.CreateFoodResponse{Status: "Failed"}, err
	}
	return &pb.CreateFoodResponse{Status: "Success"}, nil
}

func (s *Service) UpdateFood(ctx context.Context, req *pb.UpdateFoodRequest) (*pb.UpdateFoodResponse, error) {
	food := &Food{ID: req.Id, Name: req.Name, Description: req.Description, Price: req.Price}
	if err := s.store.UpdateFood(ctx, food); err != nil {
		return &pb.UpdateFoodResponse{Status: "Failed"}, err
	}
	return &pb.UpdateFoodResponse{Status: "Success"}, nil
}

func (s *Service) DeleteFood(ctx context.Context, req *pb.DeleteFoodRequest) (*pb.DeleteFoodResponse, error) {
	if err := s.store.DeleteFood(ctx, req.Id); err != nil {
		return &pb.DeleteFoodResponse{Status: "Failed"}, err
	}
	return &pb.DeleteFoodResponse{Status: "Success"}, nil
}

// Accessories
func (s *Service) GetAccessoryByID(ctx context.Context, req *pb.GetAccessoryRequest) (*pb.Accessory, error) {
	accessory, err := s.store.GetAccessoryByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Accessory{Id: accessory.ID, Name: accessory.Name, Description: accessory.Description, Price: accessory.Price}, nil
}

func (s *Service) ListAccessories(ctx context.Context, req *pb.ListAccessoryRequest) (*pb.ListAccessoryResponse, error) {
	accessories, err := s.store.ListAccessories(ctx)
	if err != nil {
		return nil, err
	}
	var pbAccessories []*pb.Accessory
	for _, accessory := range accessories {
		pbAccessories = append(pbAccessories, &pb.Accessory{Id: accessory.ID, Name: accessory.Name, Description: accessory.Description, Price: accessory.Price})
	}
	return &pb.ListAccessoryResponse{Accessories: pbAccessories}, nil
}

func (s *Service) CreateAccessory(ctx context.Context, req *pb.CreateAccessoryRequest) (*pb.CreateAccessoryResponse, error) {
	accessory := &Accessory{Name: req.Name, Description: req.Description, Price: req.Price}
	if err := s.store.CreateAccessory(ctx, accessory); err != nil {
		return &pb.CreateAccessoryResponse{Status: "Failed"}, err
	}
	return &pb.CreateAccessoryResponse{Status: "Success"}, nil
}

func (s *Service) UpdateAccessory(ctx context.Context, req *pb.UpdateAccessoryRequest) (*pb.UpdateAccessoryResponse, error) {
	accessory := &Accessory{ID: req.Id, Name: req.Name, Description: req.Description, Price: req.Price}
	if err := s.store.UpdateAccessory(ctx, accessory); err != nil {
		return &pb.UpdateAccessoryResponse{Status: "Failed"}, err
	}
	return &pb.UpdateAccessoryResponse{Status: "Success"}, nil
}

func (s *Service) DeleteAccessory(ctx context.Context, req *pb.DeleteAccessoryRequest) (*pb.DeleteAccessoryResponse, error) {
	if err := s.store.DeleteAccessory(ctx, req.Id); err != nil {
		return &pb.DeleteAccessoryResponse{Status: "Failed"}, err
	}
	return &pb.DeleteAccessoryResponse{Status: "Success"}, nil
}

// Medicines
func (s *Service) GetMedicineByID(ctx context.Context, req *pb.GetMedicineRequest) (*pb.Medicine, error) {
	medicine, err := s.store.GetMedicineByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Medicine{Id: medicine.ID, Name: medicine.Name, Description: medicine.Description, Price: medicine.Price}, nil
}

func (s *Service) ListMedicines(ctx context.Context, req *pb.ListMedicineRequest) (*pb.ListMedicineResponse, error) {
	medicines, err := s.store.ListMedicines(ctx)
	if err != nil {
		return nil, err
	}
	var pbMedicines []*pb.Medicine
	for _, medicine := range medicines {
		pbMedicines = append(pbMedicines, &pb.Medicine{Id: medicine.ID, Name: medicine.Name, Description: medicine.Description, Price: medicine.Price})
	}
	return &pb.ListMedicineResponse{Medicines: pbMedicines}, nil
}

func (s *Service) CreateMedicine(ctx context.Context, req *pb.CreateMedicineRequest) (*pb.CreateMedicineResponse, error) {
	medicine := &Medicine{Name: req.Name, Description: req.Description, Price: req.Price}
	if err := s.store.CreateMedicine(ctx, medicine); err != nil {
		return &pb.CreateMedicineResponse{Status: "Failed"}, err
	}
	return &pb.CreateMedicineResponse{Status: "Success"}, nil
}

func (s *Service) UpdateMedicine(ctx context.Context, req *pb.UpdateMedicineRequest) (*pb.UpdateMedicineResponse, error) {
	medicine := &Medicine{ID: req.Id, Name: req.Name, Description: req.Description, Price: req.Price}
	if err := s.store.UpdateMedicine(ctx, medicine); err != nil {
		return &pb.UpdateMedicineResponse{Status: "Failed"}, err
	}
	return &pb.UpdateMedicineResponse{Status: "Success"}, nil
}

func (s *Service) DeleteMedicine(ctx context.Context, req *pb.DeleteMedicineRequest) (*pb.DeleteMedicineResponse, error) {
	if err := s.store.DeleteMedicine(ctx, req.Id); err != nil {
		return &pb.DeleteMedicineResponse{Status: "Failed"}, err
	}
	return &pb.DeleteMedicineResponse{Status: "Success"}, nil
}

// ------------------ Branches ------------------
func (s *Service) GetBranchByID(ctx context.Context, req *pb.GetBranchRequest) (*pb.Branch, error) {
	branch, err := s.store.GetBranchByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Branch{Id: branch.ID, Name: branch.Name, Location: branch.Location}, nil
}

func (s *Service) ListBranches(ctx context.Context, req *pb.ListBranchRequest) (*pb.ListBranchResponse, error) {
	branches, err := s.store.ListBranches(ctx)
	if err != nil {
		return nil, err
	}
	var pbBranches []*pb.Branch
	for _, branch := range branches {
		pbBranches = append(pbBranches, &pb.Branch{Id: branch.ID, Name: branch.Name, Location: branch.Location})
	}
	return &pb.ListBranchResponse{Branches: pbBranches}, nil
}

// ------------------ Branch Inventory ------------------
func (s *Service) GetBranchInventory(ctx context.Context, req *pb.GetBranchInventoryRequest) (*pb.GetBranchInventoryResponse, error) {
	inventory, err := s.store.GetBranchInventory(ctx, req.BranchId)
	if err != nil {
		return nil, err
	}
	var pbInventory []*pb.BranchProduct
	for _, item := range inventory {
		pbInventory = append(pbInventory, &pb.BranchProduct{
			BranchId:      item.BranchID,
			ProductId:     item.ProductID,
			ProductType:   item.ProductType,
			StockQuantity: item.StockQuantity,
		})
	}
	return &pb.GetBranchInventoryResponse{Inventory: pbInventory}, nil
}

func (s *Service) UpdateBranchInventory(ctx context.Context, req *pb.UpdateBranchInventoryRequest) (*pb.UpdateBranchInventoryResponse, error) {
	err := s.store.UpdateBranchInventory(ctx, req.BranchId, req.ProductId, req.ProductType, req.StockQuantity)
	if err != nil {
		return &pb.UpdateBranchInventoryResponse{Status: "Failed"}, err
	}
	return &pb.UpdateBranchInventoryResponse{Status: "Success"}, nil
}
