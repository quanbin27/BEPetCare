package handlers

type ServiceResponse struct {
	ServiceID int32  `json:"serviceId"`
	Name      string `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	ImgURL      string  `json:"imgUrl"`
}
