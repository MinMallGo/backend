package dto

import "time"

type CouponCreate struct {
	Name             string    `json:"name" binding:"required,min=1,max=32"`                         // 优惠券名称
	DiscountType     int       `json:"type" binding:"required,oneof=1 2"`                            // 1 = 满减，2 = 折扣
	UseRange         int       `json:"use_range" binding:"required,oneof=1 2 3 4 5"`                 // 使用范围：1 全平台，2商家 3 类别 4 具体商品 5 某个规格
	Threshold        int       `json:"threshold" binding:"gt=0"`                                     // 最低消费多少才能使用
	DiscountPrice    int       `json:"discount_price" binding:"gte=0,ltefield=Threshold"`            // 优惠券金额 * 100  或者 百分之多少的折扣
	DiscountRate     int       `json:"discount_rate" binding:"gt=1,lt=100"`                          // 打折比例，1-100 .这是里放大了100倍的
	ReceiveCount     int       `json:"receive_count" binding:"required,gt=0"`                        // 可以领取的数量
	ReceiveStartTime time.Time `json:"receive_start_time" binding:"required"`                        // 开始领取时间
	ReceiveEndTime   time.Time `json:"receive_end_time" binding:"required,gtfield=ReceiveStartTime"` // 领取结束时间
	UseStartTime     time.Time `json:"use_start_time" binding:"required"`                            // 优惠券开始时间
	UseEndTime       time.Time `json:"use_end_time" binding:"required,gtfield=UseStartTime"`         // 优惠券结束时间
	TotalCount       int       `json:"total_count" binding:"required,gte=1"`                         // 优惠券总数
}

type CouponUpdate struct {
	ID int `json:"id" binding:"required"`
	CouponCreate
}

type CouponSearch struct {
	Name             string    `json:"name" binding:"omitempty,min=1,max=32"`         // 优惠券名称
	DiscountType     int       `json:"type" binding:"omitempty,oneof=1 2"`            // 1 = 满减，2 = 折扣
	UseRange         int       `json:"min_spend" binding:"omitempty,oneof=1 2 3 4 5"` // 使用范围：1 全平台，2商家 3 类别 4 具体商品 5 某个规格
	Threshold        int       `json:"Threshold" binding:"omitempty,gte=0"`           // 最低消费多少才能使用
	ReceiveStartTime time.Time `json:"receive_start_time" binding:"omitempty"`        // 开始领取时间
	ReceiveEndTime   time.Time `json:"receive_end_time" binding:"omitempty"`          // 领取结束时间
	UseStartTime     time.Time `json:"use_start_time" binding:"omitempty"`            // 优惠券开始时间
	UseEndTime       time.Time `json:"use_end_time" binding:"omitempty"`              // 优惠券结束时间
	Paginate
}

type CouponDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}
