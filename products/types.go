package main

import (
	"context"
	"time"

	pb "github.com/quanbin27/commons/genproto/products"
)

// Food - Bảng thực phẩm cho thú cưng
type Food struct {
	ID          int32     `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"size:500"`
	Price       float32   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// Accessory - Bảng phụ kiện thú cưng
type Accessory struct {
	ID          int32     `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"size:500"`
	Price       float32   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// Medicine - Bảng thuốc thú cưng
type Medicine struct {
	ID          int32     `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"size:500"`
	Price       float32   `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// Branch - Bảng chi nhánh
type Branch struct {
	ID       int32  `gorm:"primaryKey"`
	Name     string `gorm:"size:255;not null"`
	Location string `gorm:"size:500"`
}

// BranchProduct - Quản lý tồn kho theo từng chi nhánh
type BranchProduct struct {
	BranchID      int32  `gorm:"primaryKey"`
	ProductID     int32  `gorm:"primaryKey"`
	ProductType   string `gorm:"size:50;not null"` // "food", "accessory", "medicine"
	StockQuantity int32  `gorm:"not null"`
}

// ProductStore Interface - Cung cấp các thao tác với database
type ProductStore interface {
	// Thực phẩm
	GetFoodByID(ctx context.Context, id int32) (*Food, error)
	ListFoods(ctx context.Context) ([]Food, error)
	CreateFood(ctx context.Context, food *Food) error
	UpdateFood(ctx context.Context, food *Food) error
	DeleteFood(ctx context.Context, id int32) error

	// Phụ kiện
	GetAccessoryByID(ctx context.Context, id int32) (*Accessory, error)
	ListAccessories(ctx context.Context) ([]Accessory, error)
	CreateAccessory(ctx context.Context, accessory *Accessory) error
	UpdateAccessory(ctx context.Context, accessory *Accessory) error
	DeleteAccessory(ctx context.Context, id int32) error

	// Thuốc
	GetMedicineByID(ctx context.Context, id int32) (*Medicine, error)
	ListMedicines(ctx context.Context) ([]Medicine, error)
	CreateMedicine(ctx context.Context, medicine *Medicine) error
	UpdateMedicine(ctx context.Context, medicine *Medicine) error
	DeleteMedicine(ctx context.Context, id int32) error

	// Chi nhánh
	GetBranchByID(ctx context.Context, id int32) (*Branch, error)
	ListBranches(ctx context.Context) ([]Branch, error)

	// Tồn kho
	GetBranchInventory(ctx context.Context, branchID int32) ([]BranchProduct, error)
	UpdateBranchInventory(ctx context.Context, branchID int32, productID int32, productType string, stockQuantity int32) error
}

// ProductService Interface - Implement gRPC server logic
type ProductService interface {
	// Thực phẩm
	GetFoodByID(ctx context.Context, req *pb.GetFoodRequest) (*pb.Food, error)
	ListFoods(ctx context.Context, req *pb.ListFoodRequest) (*pb.ListFoodResponse, error)
	CreateFood(ctx context.Context, req *pb.CreateFoodRequest) (*pb.CreateFoodResponse, error)
	UpdateFood(ctx context.Context, req *pb.UpdateFoodRequest) (*pb.UpdateFoodResponse, error)
	DeleteFood(ctx context.Context, req *pb.DeleteFoodRequest) (*pb.DeleteFoodResponse, error)

	// Phụ kiện
	GetAccessoryByID(ctx context.Context, req *pb.GetAccessoryRequest) (*pb.Accessory, error)
	ListAccessories(ctx context.Context, req *pb.ListAccessoryRequest) (*pb.ListAccessoryResponse, error)
	CreateAccessory(ctx context.Context, req *pb.CreateAccessoryRequest) (*pb.CreateAccessoryResponse, error)
	UpdateAccessory(ctx context.Context, req *pb.UpdateAccessoryRequest) (*pb.UpdateAccessoryResponse, error)
	DeleteAccessory(ctx context.Context, req *pb.DeleteAccessoryRequest) (*pb.DeleteAccessoryResponse, error)

	// Thuốc
	GetMedicineByID(ctx context.Context, req *pb.GetMedicineRequest) (*pb.Medicine, error)
	ListMedicines(ctx context.Context, req *pb.ListMedicineRequest) (*pb.ListMedicineResponse, error)
	CreateMedicine(ctx context.Context, req *pb.CreateMedicineRequest) (*pb.CreateMedicineResponse, error)
	UpdateMedicine(ctx context.Context, req *pb.UpdateMedicineRequest) (*pb.UpdateMedicineResponse, error)
	DeleteMedicine(ctx context.Context, req *pb.DeleteMedicineRequest) (*pb.DeleteMedicineResponse, error)

	// Chi nhánh
	GetBranchByID(ctx context.Context, req *pb.GetBranchRequest) (*pb.Branch, error)
	ListBranches(ctx context.Context, req *pb.ListBranchRequest) (*pb.ListBranchResponse, error)

	// Tồn kho
	GetBranchInventory(ctx context.Context, req *pb.GetBranchInventoryRequest) (*pb.GetBranchInventoryResponse, error)
	UpdateBranchInventory(ctx context.Context, req *pb.UpdateBranchInventoryRequest) (*pb.UpdateBranchInventoryResponse, error)
}
