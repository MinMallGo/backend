package dto

type Paginate struct {
	Page int `json:"page" binding:"omitempty,min=1"`
	Size int `json:"Size" binding:"omitempty,min=1"`
}
