package dto

type SpecKeyCreate struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
	Uint string `json:"uint" binding:"omitempty,min=1,max=255"`
}

type SpecKeySearch struct {
	Name string `json:"name" binding:"omitempty,min=1,max=255"`
	Uint string `json:"uint" binding:"omitempty,min=1,max=255"`
	Paginate
}

type SpecKeyUpdate struct {
	Id   int    `json:"id" binding:"required,gt=0"`
	Name string `json:"name" binding:"omitempty,min=1,max=255"`
	Uint string `json:"uint" binding:"omitempty,min=1,max=255"`
}

type SpecKeyDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}
