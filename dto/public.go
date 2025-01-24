package dto

type Paginate struct {
	Page int `json:"page" binding:"omitempty,min=1"`
	Size int `json:"size" binding:"omitempty,min=1"`
}

type PaginateCount struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Count int         `json:"count"`
}
