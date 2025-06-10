package main

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
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
func (s *Store) ListAttachableProducts(ctx context.Context) ([]GeneralProduct, error) {
	var products []GeneralProduct

	// Query Food
	var foods []Food
	if err := s.db.Where("is_attachable = ?", true).Find(&foods).Error; err != nil {
		return nil, err
	}
	for _, food := range foods {
		products = append(products, GeneralProduct{
			Name:        food.Name,
			Description: food.Description,
			Price:       food.Price,
			ImgUrl:      food.ImgUrl,
			ProductID:   food.ID,
			ProductType: "food",
		})
	}

	// Query Accessory
	var accessories []Accessory
	if err := s.db.Where("is_attachable = ?", true).Find(&accessories).Error; err != nil {
		return nil, err
	}
	for _, accessory := range accessories {
		products = append(products, GeneralProduct{
			Name:        accessory.Name,
			Description: accessory.Description,
			Price:       accessory.Price,
			ImgUrl:      accessory.ImgUrl,
			ProductID:   accessory.ID,
			ProductType: "accessory",
		})
	}

	// Query Medicine
	var medicines []Medicine
	if err := s.db.Where("is_attachable = ?", true).Find(&medicines).Error; err != nil {
		return nil, err
	}
	for _, medicine := range medicines {
		products = append(products, GeneralProduct{
			Name:        medicine.Name,
			Description: medicine.Description,
			Price:       medicine.Price,
			ImgUrl:      medicine.ImgUrl,
			ProductID:   medicine.ID,
			ProductType: "medicine",
		})
	}

	return products, nil
}
func (s *Store) ListAllProducts(ctx context.Context) ([]GeneralProduct, error) {
	var products []GeneralProduct

	// Query Food
	var foods []Food
	if err := s.db.Find(&foods).Error; err != nil {
		return nil, err
	}
	for _, food := range foods {
		products = append(products, GeneralProduct{
			Name:         food.Name,
			Description:  food.Description,
			Price:        food.Price,
			ImgUrl:       food.ImgUrl,
			ProductID:    food.ID,
			ProductType:  "food",
			IsAttachable: food.IsAttachable,
		})
	}

	// Query Accessory
	var accessories []Accessory
	if err := s.db.Find(&accessories).Error; err != nil {
		return nil, err
	}
	for _, accessory := range accessories {
		products = append(products, GeneralProduct{
			Name:         accessory.Name,
			Description:  accessory.Description,
			Price:        accessory.Price,
			ImgUrl:       accessory.ImgUrl,
			ProductID:    accessory.ID,
			ProductType:  "accessory",
			IsAttachable: accessory.IsAttachable,
		})
	}

	// Query Medicine
	var medicines []Medicine
	if err := s.db.Find(&medicines).Error; err != nil {
		return nil, err
	}
	for _, medicine := range medicines {
		products = append(products, GeneralProduct{
			Name:         medicine.Name,
			Description:  medicine.Description,
			Price:        medicine.Price,
			ImgUrl:       medicine.ImgUrl,
			ProductID:    medicine.ID,
			ProductType:  "medicine",
			IsAttachable: medicine.IsAttachable,
		})
	}

	return products, nil
}
func (s *Store) ListAvailableProductsByBranch(ctx context.Context, branchID int32, productType string) ([]GeneralProduct, error) {
	var results []GeneralProduct
	log.Printf("Listing available products for branch %d with type %s", branchID, productType)
	// Helper xử lý 1 loại sản phẩm
	handle := func(tableName string, productTypeVal string, dest interface{}, getProductInfo func(interface{}) GeneralProduct) error {
		// Query tồn kho còn hàng
		var inventories []BranchProduct
		log.Printf("Querying inventory for branch %d and product type %s", branchID, productTypeVal)
		if err := s.db.
			Where("branch_id = ? AND product_type = ? AND stock_quantity > 0", branchID, productTypeVal).
			Find(&inventories).Error; err != nil {
			return err
		}

		// Lấy danh sách product_id
		productIDs := make([]int32, 0, len(inventories))
		log.Println("Found inventories:", len(inventories))
		log.Println("PRODUCTS:", inventories)
		stockMap := make(map[int32]int32)
		for _, inv := range inventories {
			productIDs = append(productIDs, inv.ProductID)
			stockMap[inv.ProductID] = inv.StockQuantity - inv.ReservedQuantity
		}

		if len(productIDs) == 0 {
			return nil // không có sản phẩm nào
		}

		// Lấy chi tiết sản phẩm
		if err := s.db.
			Where("id IN (?)", productIDs).
			Find(dest).Error; err != nil {
			return err
		}

		// Chuyển đổi kết quả
		switch list := dest.(type) {
		case *[]Food:
			for _, p := range *list {
				res := getProductInfo(p)
				res.AvailableQuantity = stockMap[p.ID]
				results = append(results, res)
			}
		case *[]Accessory:
			for _, p := range *list {
				res := getProductInfo(p)
				res.AvailableQuantity = stockMap[p.ID]
				results = append(results, res)
			}
		case *[]Medicine:
			for _, p := range *list {
				res := getProductInfo(p)
				res.AvailableQuantity = stockMap[p.ID]
				results = append(results, res)
			}
		}

		return nil
	}

	// Tùy theo productType lọc
	switch productType {
	case "food", "FOOD":
		if err := handle("foods", "food", &[]Food{}, func(p interface{}) GeneralProduct {
			f := p.(Food)
			return GeneralProduct{
				ProductID:   f.ID,
				ProductType: "food",
				Name:        f.Name,
				Description: f.Description,
				Price:       f.Price,
				ImgUrl:      f.ImgUrl,
			}
		}); err != nil {
			return nil, err
		}
	case "accessory", "ACCESSORY":
		if err := handle("accessories", "accessory", &[]Accessory{}, func(p interface{}) GeneralProduct {
			a := p.(Accessory)
			return GeneralProduct{
				ProductID:   a.ID,
				ProductType: "accessory",
				Name:        a.Name,
				Description: a.Description,
				Price:       a.Price,
				ImgUrl:      a.ImgUrl,
			}
		}); err != nil {
			return nil, err
		}
	case "medicine", "MEDICINE":
		if err := handle("medicines", "medicine", &[]Medicine{}, func(p interface{}) GeneralProduct {
			m := p.(Medicine)
			return GeneralProduct{
				ProductID:   m.ID,
				ProductType: "medicine",
				Name:        m.Name,
				Description: m.Description,
				Price:       m.Price,
				ImgUrl:      m.ImgUrl,
			}
		}); err != nil {
			return nil, err
		}
	}

	return results, nil
}
func (s *Store) ListAvailableAllProductsByBranch(ctx context.Context, branchID int32) ([]GeneralProduct, error) {
	var products []GeneralProduct

	// 1. Query Food
	var foods []struct {
		Food
		StockQty    int32
		ReservedQty int32
	}
	if err := s.db.
		Table("foods").
		Select("foods.*, bp.stock_quantity AS stock_qty, bp.reserved_quantity AS reserved_qty").
		Joins("JOIN branch_products bp ON bp.product_id = foods.id AND bp.product_type = ?", "food").
		Where("bp.branch_id = ? AND bp.stock_quantity > 0", branchID).
		Scan(&foods).Error; err != nil {
		return nil, err
	}
	for _, food := range foods {
		products = append(products, GeneralProduct{
			Name:              food.Name,
			Description:       food.Description,
			Price:             food.Price,
			ImgUrl:            food.ImgUrl,
			ProductID:         food.ID,
			ProductType:       "food",
			IsAttachable:      food.IsAttachable,
			AvailableQuantity: food.StockQty, // Tính số lượng có sẵn
		})
	}

	// 2. Query Accessory
	var accessories []struct {
		Accessory
		StockQty    int32
		ReservedQty int32
	}
	if err := s.db.
		Table("accessories").
		Select("accessories.*, bp.stock_quantity AS stock_qty, bp.reserved_quantity AS reserved_qty").
		Joins("JOIN branch_products bp ON bp.product_id = accessories.id AND bp.product_type = ?", "accessory").
		Where("bp.branch_id = ? AND bp.stock_quantity - bp.reserved_quantity > 0", branchID).
		Scan(&accessories).Error; err != nil {
		return nil, err
	}
	for _, acc := range accessories {
		products = append(products, GeneralProduct{
			Name:              acc.Name,
			Description:       acc.Description,
			Price:             acc.Price,
			ImgUrl:            acc.ImgUrl,
			ProductID:         acc.ID,
			ProductType:       "accessory",
			IsAttachable:      acc.IsAttachable,
			AvailableQuantity: acc.StockQty,
		})
	}

	// 3. Query Medicine
	var medicines []struct {
		Medicine
		StockQty    int32
		ReservedQty int32
	}
	if err := s.db.
		Table("medicines").
		Select("medicines.*, bp.stock_quantity AS stock_qty, bp.reserved_quantity AS reserved_qty").
		Joins("JOIN branch_products bp ON bp.product_id = medicines.id AND bp.product_type = ?", "medicine").
		Where("bp.branch_id = ? AND bp.stock_quantity - bp.reserved_quantity > 0", branchID).
		Scan(&medicines).Error; err != nil {
		return nil, err
	}
	for _, med := range medicines {
		products = append(products, GeneralProduct{
			Name:              med.Name,
			Description:       med.Description,
			Price:             med.Price,
			ImgUrl:            med.ImgUrl,
			ProductID:         med.ID,
			ProductType:       "medicine",
			IsAttachable:      med.IsAttachable,
			AvailableQuantity: med.StockQty, // Tính số lượng có sẵn
		})
	}

	return products, nil
}
func (s *Store) ReserveProduct(ctx context.Context, branchID, productID int32, productType string, quantity int32) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stock BranchProduct
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("branch_id = ? AND product_id = ? AND product_type = ?", branchID, productID, productType).
			First(&stock).Error; err != nil {
			return err
		}

		if stock.StockQuantity < quantity {
			return fmt.Errorf("not enough stock available")
		}

		stock.StockQuantity -= quantity
		stock.ReservedQuantity += quantity

		return tx.Save(&stock).Error
	})
}

func (s *Store) ConfirmPickup(ctx context.Context, branchID, productID int32, productType string, quantity int32) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stock BranchProduct
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("branch_id = ? AND product_id = ? AND product_type = ?", branchID, productID, productType).
			First(&stock).Error; err != nil {
			return err
		}

		if stock.ReservedQuantity < quantity {
			return fmt.Errorf("not enough reserved stock to confirm")
		}

		stock.ReservedQuantity -= quantity
		// Không cộng lại vào available vì đã lấy ra.

		return tx.Save(&stock).Error
	})
}

func (s *Store) ReleaseReservation(ctx context.Context, branchID, productID int32, productType string, quantity int32) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stock BranchProduct
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("branch_id = ? AND product_id = ? AND product_type = ?", branchID, productID, productType).
			First(&stock).Error; err != nil {
			return err
		}

		if stock.ReservedQuantity < quantity {
			return fmt.Errorf("not enough reserved stock to release")
		}

		stock.ReservedQuantity -= quantity
		stock.StockQuantity += quantity

		return tx.Save(&stock).Error
	})
}
