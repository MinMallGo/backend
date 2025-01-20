package dto

type CouponCreate struct {
	Name      string `json:"name" binding:"required,min=1,max=32"`
	Type      int    `json:"type" binding:"required,oneof=1 2"`
	MinSpend  int    `json:"min_spend" binding:"required,gte=0"`
	Price     int    `json:"price" binding:"required,gte=0"`
	StartTime string `json:"start_time" binding:"required,datetime=2006-01-02 15:04:05"`
	EndTime   string `json:"end_time" binding:"required,datetime=2006-01-02  15:04:05"`
}

type CouponDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}

type CouponUpdate struct {
	Id int `json:"id" binding:"required,gt=0"`
	CouponUpdateOrSearch
}

type CouponSearch struct {
	Id int `json:"id" binding:"omitempty,gt=0"`
	CouponUpdateOrSearch
	Page  int `json:"page" binding:"gt=0"`
	Limit int `json:"limit"`
}

type CouponUpdateOrSearch struct {
	Name      string `json:"name" binding:"omitempty,min=1,max=32"`
	Type      int    `json:"type" binding:"omitempty,oneof=1 2"`
	MinSpend  int    `json:"min_spend" binding:"omitempty,gte=0"`
	Price     int    `json:"price" binding:"omitempty,gte=0"`
	StartTime string `json:"start_time" binding:"omitempty,datetime=2006-01-02 15:04:05"`
	EndTime   string `json:"end_time" binding:"omitempty,datetime=2006-01-02 15:04:05"`
}
