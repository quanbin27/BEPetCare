package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	pb "github.com/quanbin27/commons/genproto/products"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	client pb.ProductServiceClient
}

func NewProductHandler(client pb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{client: client}
}

// RegisterRoutes đăng ký các route cho Products service
func (h *ProductHandler) RegisterRoutes(e *echo.Group) {
	// Thực phẩm
	e.GET("/products/foods/:id", h.GetFoodByID)
	e.GET("/products/foods", h.ListFoods)
	e.POST("/products/foods", h.CreateFood)
	e.PUT("/products/foods", h.UpdateFood)
	e.DELETE("/products/foods/:id", h.DeleteFood)

	// Phụ kiện
	e.GET("/products/accessories/:id", h.GetAccessoryByID)
	e.GET("/products/accessories", h.ListAccessories)
	e.POST("/products/accessories", h.CreateAccessory)
	e.PUT("/products/accessories", h.UpdateAccessory)
	e.DELETE("/products/accessories/:id", h.DeleteAccessory)

	// Thuốc
	e.GET("/products/medicines/:id", h.GetMedicineByID)
	e.GET("/products/medicines", h.ListMedicines)
	e.POST("/products/medicines", h.CreateMedicine)
	e.PUT("/products/medicines", h.UpdateMedicine)
	e.DELETE("/products/medicines/:id", h.DeleteMedicine)

	// Chi nhánh
	e.GET("/branches/:id", h.GetBranchByID)
	e.GET("/branches", h.ListBranches)

	// Tồn kho
	e.GET("/branches/:branch_id/inventory", h.GetBranchInventory)
	e.PUT("/branches/inventory", h.UpdateBranchInventory)

	e.GET("/products/is_attachable", h.ListAttachableProduct)
	e.GET("/products", h.ListAllProduct)
}

// --- Thực phẩm ---

func (h *ProductHandler) GetFoodByID(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.GetFoodByID(ctx, &pb.GetFoodRequest{Id: int32(id)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Food not found"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) ListFoods(c echo.Context) error {
	ctx := c.Request().Context()
	resp, err := h.client.ListFoods(ctx, &pb.ListFoodRequest{})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var foods []ProductResponse
	for _, food := range resp.Foods {
		foods = append(foods, ProductResponse{
			ID:           food.Id,
			Name:         food.Name,
			Price:        food.Price,
			Description:  food.Description,
			ImgURL:       food.Imgurl,
			IsAttachable: food.IsAttachable,
		})
	}

	return c.JSON(http.StatusOK, foods)
}

func (h *ProductHandler) CreateFood(c echo.Context) error {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validation cơ bản
	if req.Name == "" || req.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required and price must be positive"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.CreateFood(ctx, &pb.CreateFoodRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Giả sử service trả về ID trong status hoặc cần thêm field trong proto
	return c.JSON(http.StatusOK, map[string]string{
		"status": resp.Status,
		// "id":     "newly_created_id", // Nếu proto trả về ID, thêm vào đây
	})
}

func (h *ProductHandler) UpdateFood(c echo.Context) error {
	var req struct {
		ID          int32   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ID <= 0 || req.Name == "" || req.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID, name are required and price must be positive"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdateFood(ctx, &pb.UpdateFoodRequest{
		Id:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid arguments"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

func (h *ProductHandler) DeleteFood(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.DeleteFood(ctx, &pb.DeleteFoodRequest{Id: int32(id)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

// --- Phụ kiện ---

func (h *ProductHandler) GetAccessoryByID(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.GetAccessoryByID(ctx, &pb.GetAccessoryRequest{Id: int32(id)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Accessory not found"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) ListAccessories(c echo.Context) error {
	ctx := c.Request().Context()
	resp, err := h.client.ListAccessories(ctx, &pb.ListAccessoryRequest{})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var accessories []ProductResponse
	for _, acc := range resp.Accessories {
		accessories = append(accessories, ProductResponse{
			ID:           acc.Id,
			Name:         acc.Name,
			Price:        acc.Price,
			Description:  acc.Description,
			ImgURL:       acc.Imgurl,
			IsAttachable: acc.IsAttachable,
		})
	}

	return c.JSON(http.StatusOK, accessories)
}

func (h *ProductHandler) CreateAccessory(c echo.Context) error {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Name == "" || req.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required and price must be positive"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.CreateAccessory(ctx, &pb.CreateAccessoryRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

func (h *ProductHandler) UpdateAccessory(c echo.Context) error {
	var req struct {
		ID          int32   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ID <= 0 || req.Name == "" || req.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID, name are required and price must be positive"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdateAccessory(ctx, &pb.UpdateAccessoryRequest{
		Id:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid arguments"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

func (h *ProductHandler) DeleteAccessory(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.DeleteAccessory(ctx, &pb.DeleteAccessoryRequest{Id: int32(id)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

// --- Thuốc ---

func (h *ProductHandler) GetMedicineByID(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.GetMedicineByID(ctx, &pb.GetMedicineRequest{Id: int32(id)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Medicine not found"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) ListMedicines(c echo.Context) error {
	ctx := c.Request().Context()
	resp, err := h.client.ListMedicines(ctx, &pb.ListMedicineRequest{})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var medicines []ProductResponse
	for _, med := range resp.Medicines {
		medicines = append(medicines, ProductResponse{
			ID:           med.Id,
			Name:         med.Name,
			Price:        med.Price,
			Description:  med.Description,
			ImgURL:       med.Imgurl,
			IsAttachable: med.IsAttachable,
		})
	}

	return c.JSON(http.StatusOK, medicines)
}

func (h *ProductHandler) CreateMedicine(c echo.Context) error {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.Name == "" || req.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required and price must be positive"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.CreateMedicine(ctx, &pb.CreateMedicineRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

func (h *ProductHandler) UpdateMedicine(c echo.Context) error {
	var req struct {
		ID          int32   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if req.ID <= 0 || req.Name == "" || req.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID, name are required and price must be positive"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdateMedicine(ctx, &pb.UpdateMedicineRequest{
		Id:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid arguments"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

func (h *ProductHandler) DeleteMedicine(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.DeleteMedicine(ctx, &pb.DeleteMedicineRequest{Id: int32(id)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}

// --- Chi nhánh ---

func (h *ProductHandler) GetBranchByID(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format, must be an integer"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.GetBranchByID(ctx, &pb.GetBranchRequest{Id: int32(id)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.NotFound:
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Branch not found"})
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) ListBranches(c echo.Context) error {
	ctx := c.Request().Context()
	resp, err := h.client.ListBranches(ctx, &pb.ListBranchRequest{})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Branches)
}

// --- Tồn kho ---

func (h *ProductHandler) GetBranchInventory(c echo.Context) error {
	branchIDStr := c.Param("branch_id")
	if branchIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Branch ID is required"})
	}

	branchID, err := strconv.ParseInt(branchIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid branch_id format, must be an integer"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.GetBranchInventory(ctx, &pb.GetBranchInventoryRequest{BranchId: int32(branchID)})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Inventory)
}

func (h *ProductHandler) UpdateBranchInventory(c echo.Context) error {
	var req struct {
		BranchID      int32  `json:"branch_id"`
		ProductID     int32  `json:"product_id"`
		ProductType   string `json:"product_type"` // Giữ nguyên là string
		StockQuantity int32  `json:"stock_quantity"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validation cơ bản
	if req.BranchID <= 0 || req.ProductID <= 0 || req.StockQuantity < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Branch ID and Product ID must be positive, stock quantity must be non-negative"})
	}

	// Kiểm tra ProductType có hợp lệ không (tuỳ chọn)
	validProductTypes := map[string]bool{
		"FOOD":      true,
		"ACCESSORY": true,
		"MEDICINE":  true,
	}
	if req.ProductType == "" || !validProductTypes[req.ProductType] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid or missing product_type (must be FOOD, ACCESSORY, or MEDICINE)"})
	}

	ctx := c.Request().Context()
	resp, err := h.client.UpdateBranchInventory(ctx, &pb.UpdateBranchInventoryRequest{
		BranchId:      req.BranchID,
		ProductId:     req.ProductID,
		ProductType:   req.ProductType, // Truyền trực tiếp string
		StockQuantity: req.StockQuantity,
	})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": resp.Status})
}
func (h *ProductHandler) ListAttachableProduct(c echo.Context) error {
	ctx := c.Request().Context()
	resp, err := h.client.ListAttachableProducts(ctx, &pb.ListAttachableProductsRequest{})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp.Products)
}
func (h *ProductHandler) ListAllProduct(c echo.Context) error {
	ctx := c.Request().Context()
	resp, err := h.client.ListAllProducts(ctx, &pb.ListAllProductsRequest{})
	if err != nil {
		if grpcErr, ok := status.FromError(err); ok {
			switch grpcErr.Code() {
			case codes.Internal:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var products []AllProductResponse
	for _, prod := range resp.Products {
		products = append(products, AllProductResponse{
			ID:           prod.ProductId,
			Name:         prod.Name,
			Price:        prod.Price,
			Description:  prod.Description,
			ImgURL:       prod.Imgurl,
			ProductType:  prod.ProductType,
			IsAttachable: prod.IsAttachable, // Đảm bảo trường này được gán
		})
	}

	return c.JSON(http.StatusOK, products)
}
