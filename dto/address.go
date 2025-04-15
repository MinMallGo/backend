package dto

type AddressCreate struct {
	Name        string `json:"name" binding:"required,min=2,max=20"`
	Phone       string `json:"phone" binding:"required,len=11"`
	Address     string `json:"address" binding:"required,min=1,max=255"`
	HouseNumber string `json:"house_number" binding:"required,min=1,max=255"`
	IsDefault   int    `json:"is_default" binding:"required,oneof=0 1"`
	Tag         string `json:"tag" binding:"omitempty,min=1,max=20"`
}
