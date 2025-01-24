package dto

type SpecValueCreate struct {
	SpecKeyId int    `json:"spec_key_id" binding:"required,gt=0"`
	Name      string `json:"name" binding:"required,min=0,max=255"`
}

type SpecValueSearch struct {
	SpecKeyId int    `json:"spec_key_id" binding:"omitempty,gt=0"`
	Name      string `json:"name" binding:"omitempty,min=0,max=255"`
	Paginate
}

type SpecValueUpdate struct {
	Id        int    `json:"id" binding:"required,gt=0"`
	SpecKeyId int    `json:"spec_key_id" binding:"omitempty,gt=0"`
	Name      string `json:"name" binding:"omitempty,min=0,max=255"`
}

type SpecValueDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}
