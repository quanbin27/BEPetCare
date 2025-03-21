package main

import (
	"context"
	"time"

	pb "github.com/quanbin27/commons/genproto/products"
)

// Food - Bảng thực phẩm cho thú cưng
type Food struct {
	ID           int32     `gorm:"primaryKey"`
	Name         string    `gorm:"size:255;not null"`
	Description  string    `gorm:"size:500"`
	Price        float32   `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	IsAttachable bool      `gorm:"default:false"`
	ImgUrl       string
}

// Accessory - Bảng phụ kiện thú cưng
type Accessory struct {
	ID           int32     `gorm:"primaryKey"`
	Name         string    `gorm:"size:255;not null"`
	Description  string    `gorm:"size:500"`
	Price        float32   `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	IsAttachable bool      `gorm:"default:false"`
	ImgUrl       string
}

// Medicine - Bảng thuốc thú cưng
type Medicine struct {
	ID           int32     `gorm:"primaryKey"`
	Name         string    `gorm:"size:255;not null"`
	Description  string    `gorm:"size:500"`
	Price        float32   `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	IsAttachable bool      `gorm:"default:false"`
	ImgUrl       string
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
type GeneralProduct struct {
	Name         string
	Description  string
	Price        float32
	ImgUrl       string
	ProductID    int32
	ProductType  string
	IsAttachable bool
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
	ListAttachableProducts(ctx context.Context) ([]GeneralProduct, error)
	ListAllProducts(ctx context.Context) ([]GeneralProduct, error)
}

// ProductService Interface - Implement business logic with internal types
type ProductService interface {
	// Thực phẩm
	GetFoodByID(ctx context.Context, id int32) (*Food, error)
	ListFoods(ctx context.Context) ([]Food, error)
	CreateFood(ctx context.Context, name, description string, price float32) (string, error)
	UpdateFood(ctx context.Context, id int32, name, description string, price float32) (string, error)
	DeleteFood(ctx context.Context, id int32) (string, error)

	// Phụ kiện
	GetAccessoryByID(ctx context.Context, id int32) (*Accessory, error)
	ListAccessories(ctx context.Context) ([]Accessory, error)
	CreateAccessory(ctx context.Context, name, description string, price float32) (string, error)
	UpdateAccessory(ctx context.Context, id int32, name, description string, price float32) (string, error)
	DeleteAccessory(ctx context.Context, id int32) (string, error)

	// Thuốc
	GetMedicineByID(ctx context.Context, id int32) (*Medicine, error)
	ListMedicines(ctx context.Context) ([]Medicine, error)
	CreateMedicine(ctx context.Context, name, description string, price float32) (string, error)
	UpdateMedicine(ctx context.Context, id int32, name, description string, price float32) (string, error)
	DeleteMedicine(ctx context.Context, id int32) (string, error)

	// Chi nhánh
	GetBranchByID(ctx context.Context, id int32) (*Branch, error)
	ListBranches(ctx context.Context) ([]Branch, error)

	// Tồn kho
	GetBranchInventory(ctx context.Context, branchID int32) ([]BranchProduct, error)
	UpdateBranchInventory(ctx context.Context, branchID, productID int32, productType string, stockQuantity int32) (string, error)
	ListAttachableProducts(ctx context.Context) ([]GeneralProduct, error)
	ListAllProducts(ctx context.Context) ([]GeneralProduct, error)
}

// Helper functions to convert between internal types and protobuf types
func toProtoFood(f *Food) *pb.Food {
	return &pb.Food{
		Id:           f.ID,
		Name:         f.Name,
		Description:  f.Description,
		Price:        f.Price,
		Imgurl:       f.ImgUrl,
		IsAttachable: f.IsAttachable,
	}
}
func toProtoGeneralProduct(f *GeneralProduct) *pb.GeneralProduct {
	return &pb.GeneralProduct{
		Name:         f.Name,
		Description:  f.Description,
		Price:        f.Price,
		Imgurl:       f.ImgUrl,
		ProductType:  f.ProductType,
		ProductId:    f.ProductID,
		IsAttachable: f.IsAttachable,
	}
}
func toProtoAccessory(a *Accessory) *pb.Accessory {
	return &pb.Accessory{
		Id:           a.ID,
		Name:         a.Name,
		Description:  a.Description,
		Price:        a.Price,
		Imgurl:       a.ImgUrl,
		IsAttachable: a.IsAttachable,
	}
}

func toProtoMedicine(m *Medicine) *pb.Medicine {
	return &pb.Medicine{
		Id:           m.ID,
		Name:         m.Name,
		Description:  m.Description,
		Price:        m.Price,
		Imgurl:       m.ImgUrl,
		IsAttachable: m.IsAttachable,
	}
}

func toProtoBranch(b *Branch) *pb.Branch {
	return &pb.Branch{
		Id:       b.ID,
		Name:     b.Name,
		Location: b.Location,
	}
}

func toProtoBranchProduct(bp *BranchProduct) *pb.BranchProduct {
	return &pb.BranchProduct{
		BranchId:      bp.BranchID,
		ProductId:     bp.ProductID,
		ProductType:   bp.ProductType,
		StockQuantity: bp.StockQuantity,
	}
}
