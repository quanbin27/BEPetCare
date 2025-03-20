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
