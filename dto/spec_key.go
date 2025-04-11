package dto

type SpecKeyCreate struct {
	Name  string   `json:"name" binding:"required,min=1,max=255"`
	Uint  string   `json:"uint" binding:"omitempty,min=1,max=255"`
	Value []string `json:"value" binding:"required,dive,min=1"`
}

type SpecKeyCreateResp struct {
	Id     int                      `json:"id"`
	Values []SpecKeyCreateRespValue `json:"values"`
}
type SpecKeyCreateRespValue struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
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
