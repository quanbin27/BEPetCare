package main

import (
	"context"
	"google.golang.org/grpc"

	pb "github.com/quanbin27/commons/genproto/products"
)

type ProductsGrpcHandler struct {
	pb.UnimplementedProductServiceServer
	service ProductService
}

func NewGrpcProductHandler(grpc *grpc.Server, productService ProductService) {
	grpcHandler := &ProductsGrpcHandler{
		service: productService,
	}
	pb.RegisterProductServiceServer(grpc, grpcHandler)
}

// ------------------ Foods ------------------
func (h *ProductsGrpcHandler) GetFoodByID(ctx context.Context, req *pb.GetFoodRequest) (*pb.Food, error) {
	return h.service.GetFoodByID(ctx, req)
}

func (h *ProductsGrpcHandler) ListFoods(ctx context.Context, req *pb.ListFoodRequest) (*pb.ListFoodResponse, error) {
	return h.service.ListFoods(ctx, req)
}

func (h *ProductsGrpcHandler) CreateFood(ctx context.Context, req *pb.CreateFoodRequest) (*pb.CreateFoodResponse, error) {
	return h.service.CreateFood(ctx, req)
}

func (h *ProductsGrpcHandler) UpdateFood(ctx context.Context, req *pb.UpdateFoodRequest) (*pb.UpdateFoodResponse, error) {
	return h.service.UpdateFood(ctx, req)
}

func (h *ProductsGrpcHandler) DeleteFood(ctx context.Context, req *pb.DeleteFoodRequest) (*pb.DeleteFoodResponse, error) {
	return h.service.DeleteFood(ctx, req)
}

// ------------------ Accessories ------------------
func (h *ProductsGrpcHandler) GetAccessoryByID(ctx context.Context, req *pb.GetAccessoryRequest) (*pb.Accessory, error) {
	return h.service.GetAccessoryByID(ctx, req)
}

func (h *ProductsGrpcHandler) ListAccessories(ctx context.Context, req *pb.ListAccessoryRequest) (*pb.ListAccessoryResponse, error) {
	return h.service.ListAccessories(ctx, req)
}

func (h *ProductsGrpcHandler) CreateAccessory(ctx context.Context, req *pb.CreateAccessoryRequest) (*pb.CreateAccessoryResponse, error) {
	return h.service.CreateAccessory(ctx, req)
}

func (h *ProductsGrpcHandler) UpdateAccessory(ctx context.Context, req *pb.UpdateAccessoryRequest) (*pb.UpdateAccessoryResponse, error) {
	return h.service.UpdateAccessory(ctx, req)
}

func (h *ProductsGrpcHandler) DeleteAccessory(ctx context.Context, req *pb.DeleteAccessoryRequest) (*pb.DeleteAccessoryResponse, error) {
	return h.service.DeleteAccessory(ctx, req)
}

// ------------------ Medicines ------------------
func (h *ProductsGrpcHandler) GetMedicineByID(ctx context.Context, req *pb.GetMedicineRequest) (*pb.Medicine, error) {
	return h.service.GetMedicineByID(ctx, req)
}

func (h *ProductsGrpcHandler) ListMedicines(ctx context.Context, req *pb.ListMedicineRequest) (*pb.ListMedicineResponse, error) {
	return h.service.ListMedicines(ctx, req)
}

func (h *ProductsGrpcHandler) CreateMedicine(ctx context.Context, req *pb.CreateMedicineRequest) (*pb.CreateMedicineResponse, error) {
	return h.service.CreateMedicine(ctx, req)
}

func (h *ProductsGrpcHandler) UpdateMedicine(ctx context.Context, req *pb.UpdateMedicineRequest) (*pb.UpdateMedicineResponse, error) {
	return h.service.UpdateMedicine(ctx, req)
}

func (h *ProductsGrpcHandler) DeleteMedicine(ctx context.Context, req *pb.DeleteMedicineRequest) (*pb.DeleteMedicineResponse, error) {
	return h.service.DeleteMedicine(ctx, req)
}

// ------------------ Branches ------------------
func (h *ProductsGrpcHandler) GetBranchByID(ctx context.Context, req *pb.GetBranchRequest) (*pb.Branch, error) {
	return h.service.GetBranchByID(ctx, req)
}

func (h *ProductsGrpcHandler) ListBranches(ctx context.Context, req *pb.ListBranchRequest) (*pb.ListBranchResponse, error) {
	return h.service.ListBranches(ctx, req)
}

// ------------------ Branch Inventory ------------------
func (h *ProductsGrpcHandler) GetBranchInventory(ctx context.Context, req *pb.GetBranchInventoryRequest) (*pb.GetBranchInventoryResponse, error) {
	return h.service.GetBranchInventory(ctx, req)
}

func (h *ProductsGrpcHandler) UpdateBranchInventory(ctx context.Context, req *pb.UpdateBranchInventoryRequest) (*pb.UpdateBranchInventoryResponse, error) {
	return h.service.UpdateBranchInventory(ctx, req)
}
