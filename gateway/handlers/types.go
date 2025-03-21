package handlers

type ServiceResponse struct {
	ServiceID   int32   `json:"serviceId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	ImgURL      string  `json:"imgUrl"`
}
type UserResponse struct {
	UserID      int32  `json:"userId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}
type AllProductResponse struct {
	ID           int32   `json:"id"`
	Name         string  `json:"name"`
	Price        float32 `json:"price"`
	Description  string  `json:"description"`
	ImgURL       string  `json:"imgUrl"`
	ProductType  string  `json:"productType"`
	IsAttachable bool    `json:"isAttachable"`
}
type ProductResponse struct {
	ID           int32   `json:"id"`
	Name         string  `json:"name"`
	Price        float32 `json:"price"`
	Description  string  `json:"description"`
	ImgURL       string  `json:"imgUrl"`
	IsAttachable bool    `json:"isAttachable"`
}
