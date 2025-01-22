package dto

type SpecValCreate struct {
	SpecId string `json:"spec_id" binding:"required,gt=0"`
	Name   string `json:"name" binding:"required,min=0,max=255"`
}

type SpecValSearch struct {
	SpecId string `json:"spec_id" binding:"omitempty,gt=0"`
	Name   string `json:"name" binding:"omitempty,min=0,max=255"`
	Paginate
}

type SpecValUpdate struct {
	Id     int    `json:"id" binding:"required,gt=0"`
	SpecId string `json:"spec_id" binding:"omitempty,gt=0"`
	Name   string `json:"name" binding:"omitempty,min=0,max=255"`
}

type SpecValDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}
