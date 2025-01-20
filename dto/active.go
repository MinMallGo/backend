package dto

type ActiveCreate struct {
	Name      string `json:"name" binding:"required,min=1,max=32"`
	Type      int    `json:"type" binding:"required,oneof=1 2"`
	Desc      string `json:"desc" binding:"omitempty,max=255"`
	StartTime string `json:"start_time" binding:"required,datetime=2006-01-02 15:04:05"`
	EndTime   string `json:"end_time" binding:"required,datetime=2006-01-02  15:04:05"`
	Coupons   string `json:"coupons" binding:"omitempty,dive"`
}

type ActiveDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}

type ActiveUpdate struct {
	Id int `json:"id" binding:"required,gt=0"`
	ActiveUpdateOrSearch
}

type ActiveSearch struct {
	Id int `json:"id" binding:"omitempty,gt=0"`
	ActiveUpdateOrSearch
	Page  int `json:"page" binding:"gt=0"`
	Limit int `json:"limit"`
}

type ActiveUpdateOrSearch struct {
	Name      string `json:"name" binding:"omitempty,min=1,max=32"`
	Type      int    `json:"type" binding:"omitempty,oneof=1 2"`
	Desc      string `json:"desc" binding:"omitempty,max=255"`
	StartTime string `json:"start_time" binding:"omitempty,datetime=2006-01-02 15:04:05"`
	EndTime   string `json:"end_time" binding:"omitempty,datetime=2006-01-02  15:04:05"`
	Coupons   string `json:"coupons" binding:"omitempty,dive"`
}
