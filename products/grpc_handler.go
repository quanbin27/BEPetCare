package main

import (
	"context"
	pb "github.com/quanbin27/commons/genproto/products"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGrpcHandler struct {
	productService ProductService
	pb.UnimplementedProductServiceServer
}

func NewProductGrpcHandler(grpc *grpc.Server, productService ProductService) {
	grpcHandler := &ProductGrpcHandler{
		productService: productService,
	}
	pb.RegisterProductServiceServer(grpc, grpcHandler)
}

// Thực phẩm
func (h *ProductGrpcHandler) GetFoodByID(ctx context.Context, req *pb.GetFoodRequest) (*pb.Food, error) {
	food, err := h.productService.GetFoodByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return toProtoFood(food), nil
}

func (h *ProductGrpcHandler) ListFoods(ctx context.Context, req *pb.ListFoodRequest) (*pb.ListFoodResponse, error) {
	foods, err := h.productService.ListFoods(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.ListFoodResponse{}
	for _, f := range foods {
		resp.Foods = append(resp.Foods, toProtoFood(&f))
	}
	return resp, nil
}

func (h *ProductGrpcHandler) CreateFood(ctx context.Context, req *pb.CreateFoodRequest) (*pb.CreateFoodResponse, error) {
	stt, err := h.productService.CreateFood(ctx, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateFoodResponse{Status: stt}, nil
}

func (h *ProductGrpcHandler) UpdateFood(ctx context.Context, req *pb.UpdateFoodRequest) (*pb.UpdateFoodResponse, error) {
	stt, err := h.productService.UpdateFood(ctx, req.Id, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &pb.UpdateFoodResponse{Status: stt}, nil
}

func (h *ProductGrpcHandler) DeleteFood(ctx context.Context, req *pb.DeleteFoodRequest) (*pb.DeleteFoodResponse, error) {
	stt, err := h.productService.DeleteFood(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteFoodResponse{Status: stt}, nil
}

// Phụ kiện
func (h *ProductGrpcHandler) GetAccessoryByID(ctx context.Context, req *pb.GetAccessoryRequest) (*pb.Accessory, error) {
	accessory, err := h.productService.GetAccessoryByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return toProtoAccessory(accessory), nil
}

func (h *ProductGrpcHandler) ListAccessories(ctx context.Context, req *pb.ListAccessoryRequest) (*pb.ListAccessoryResponse, error) {
	accessories, err := h.productService.ListAccessories(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.ListAccessoryResponse{}
	for _, a := range accessories {
		resp.Accessories = append(resp.Accessories, toProtoAccessory(&a))
	}
	return resp, nil
}

func (h *ProductGrpcHandler) CreateAccessory(ctx context.Context, req *pb.CreateAccessoryRequest) (*pb.CreateAccessoryResponse, error) {
	stt, err := h.productService.CreateAccessory(ctx, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateAccessoryResponse{Status: stt}, nil
}

func (h *ProductGrpcHandler) UpdateAccessory(ctx context.Context, req *pb.UpdateAccessoryRequest) (*pb.UpdateAccessoryResponse, error) {
	stt, err := h.productService.UpdateAccessory(ctx, req.Id, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &pb.UpdateAccessoryResponse{Status: stt}, nil
}

func (h *ProductGrpcHandler) DeleteAccessory(ctx context.Context, req *pb.DeleteAccessoryRequest) (*pb.DeleteAccessoryResponse, error) {
	stt, err := h.productService.DeleteAccessory(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteAccessoryResponse{Status: stt}, nil
}

// Thuốc
func (h *ProductGrpcHandler) GetMedicineByID(ctx context.Context, req *pb.GetMedicineRequest) (*pb.Medicine, error) {
	medicine, err := h.productService.GetMedicineByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return toProtoMedicine(medicine), nil
}

func (h *ProductGrpcHandler) ListMedicines(ctx context.Context, req *pb.ListMedicineRequest) (*pb.ListMedicineResponse, error) {
	medicines, err := h.productService.ListMedicines(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.ListMedicineResponse{}
	for _, m := range medicines {
		resp.Medicines = append(resp.Medicines, toProtoMedicine(&m))
	}
	return resp, nil
}

func (h *ProductGrpcHandler) CreateMedicine(ctx context.Context, req *pb.CreateMedicineRequest) (*pb.CreateMedicineResponse, error) {
	stt, err := h.productService.CreateMedicine(ctx, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateMedicineResponse{Status: stt}, nil
}

func (h *ProductGrpcHandler) UpdateMedicine(ctx context.Context, req *pb.UpdateMedicineRequest) (*pb.UpdateMedicineResponse, error) {
	stt, err := h.productService.UpdateMedicine(ctx, req.Id, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &pb.UpdateMedicineResponse{Status: stt}, nil
}

func (h *ProductGrpcHandler) DeleteMedicine(ctx context.Context, req *pb.DeleteMedicineRequest) (*pb.DeleteMedicineResponse, error) {
	stt, err := h.productService.DeleteMedicine(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.DeleteMedicineResponse{Status: stt}, nil
}

// Chi nhánh
func (h *ProductGrpcHandler) GetBranchByID(ctx context.Context, req *pb.GetBranchRequest) (*pb.Branch, error) {
	branch, err := h.productService.GetBranchByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return toProtoBranch(branch), nil
}

func (h *ProductGrpcHandler) ListBranches(ctx context.Context, req *pb.ListBranchRequest) (*pb.ListBranchResponse, error) {
	branches, err := h.productService.ListBranches(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.ListBranchResponse{}
	for _, b := range branches {
		resp.Branches = append(resp.Branches, toProtoBranch(&b))
	}
	return resp, nil
}

// Tồn kho
func (h *ProductGrpcHandler) GetBranchInventory(ctx context.Context, req *pb.GetBranchInventoryRequest) (*pb.GetBranchInventoryResponse, error) {
	inventory, err := h.productService.GetBranchInventory(ctx, req.BranchId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &pb.GetBranchInventoryResponse{}
	for _, bp := range inventory {
		resp.Inventory = append(resp.Inventory, toProtoBranchProduct(&bp))
	}
	return resp, nil
}

func (h *ProductGrpcHandler) UpdateBranchInventory(ctx context.Context, req *pb.UpdateBranchInventoryRequest) (*pb.UpdateBranchInventoryResponse, error) {
	stt, err := h.productService.UpdateBranchInventory(ctx, req.BranchId, req.ProductId, req.ProductType, req.StockQuantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.UpdateBranchInventoryResponse{Status: stt}, nil
}
func (h *ProductGrpcHandler) ListAttachableProducts(ctx context.Context, req *pb.ListAttachableProductsRequest) (*pb.ListAttachableProductsResponse, error) {
	products, err := h.productService.ListAttachableProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.ListAttachableProductsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, toProtoGeneralProduct(&p))
	}
	return resp, nil
}
func (h *ProductGrpcHandler) ListAllProducts(ctx context.Context, req *pb.ListAllProductsRequest) (*pb.ListAllProductsResponse, error) {
	products, err := h.productService.ListAllProducts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.ListAllProductsResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, toProtoGeneralProduct(&p))
	}
	return resp, nil
}
func (h *ProductGrpcHandler) ListAllProductsWithStock(ctx context.Context, req *pb.ListAllProductsRequest) (*pb.ListAllProductsWithStockResponse, error) {
	products, err := h.productService.ListAllProductsWithStock(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &pb.ListAllProductsWithStockResponse{}

	for _, p := range products {
		resp.Products = append(resp.Products, toProtoProductWithStock(&p))
	}
	return resp, nil
}

func (h *ProductGrpcHandler) ListAvailableProductsByBranch(ctx context.Context, req *pb.ListAvailableProductsByBranchRequest) (*pb.ListAvailableProductsByBranchResponse, error) {
	products, err := h.productService.ListAvailableProductsByBranch(ctx, req.BranchId, req.ProductType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.ListAvailableProductsByBranchResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, toProtoGeneralProduct(&p))
	}
	return resp, nil
}
func (h *ProductGrpcHandler) ListAvailableAllProductsByBranch(ctx context.Context, req *pb.ListAvailableAllProductsByBranchRequest) (*pb.ListAvailableAllProductsByBranchResponse, error) {
	products, err := h.productService.ListAvailableAllProductsByBranch(ctx, req.BranchId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.ListAvailableAllProductsByBranchResponse{}
	for _, p := range products {
		resp.Products = append(resp.Products, toProtoGeneralProduct(&p))
	}
	return resp, nil
}
func (h *ProductGrpcHandler) ReserveProduct(ctx context.Context, req *pb.ReserveProductRequest) (*pb.ReserveProductResponse, error) {
	err := h.productService.ReserveProduct(ctx, req.BranchId, req.ProductId, req.ProductType, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to reserve product: %v", err)
	}
	return &pb.ReserveProductResponse{}, nil
}

func (h *ProductGrpcHandler) ConfirmPickup(ctx context.Context, req *pb.ConfirmPickupRequest) (*pb.ConfirmPickupResponse, error) {
	err := h.productService.ConfirmPickup(ctx, req.BranchId, req.ProductId, req.ProductType, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to confirm pickup: %v", err)
	}
	return &pb.ConfirmPickupResponse{}, nil
}

func (h *ProductGrpcHandler) ReleaseReservation(ctx context.Context, req *pb.ReleaseReservationRequest) (*pb.ReleaseReservationResponse, error) {
	err := h.productService.ReleaseReservation(ctx, req.BranchId, req.ProductId, req.ProductType, req.Quantity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to release reservation: %v", err)
	}
	return &pb.ReleaseReservationResponse{}, nil
}
