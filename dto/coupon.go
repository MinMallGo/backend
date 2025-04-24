package dto

import "time"

type CouponCreate struct {
	Name             string    `json:"name" binding:"required,min=1,max=32"`                         // 优惠券名称
	ScopeID          int       `json:"usage_scope_id" binding:"required,gt=0"`                       // 关联了表的，coupon_usage_scope
	ThresholdType    int       `json:"threshold_type" binding:"gt=0,oneof=1 2"`                      // 门槛：1 表示金额 2 表示件数
	ThresholdValue   int       `json:"threshold_value" binding:"required,gt=0"`                      // 门槛：200或者3件
	DiscountType     int       `json:"discount_type" binding:"required,oneof=1 2"`                   // 1 = 满减，2 = 折扣
	DiscountValue    int       `json:"discount_value" binding:"gte=0"`                               // 优惠券金额 * 100  或者 百分之多少的折扣
	ReceiveLimit     int       `json:"user_receive_limit" binding:"omitempty,gt=0,oneof=0 1"`        // 是否限制领取次数
	ReceiveStartTime time.Time `json:"receive_start_time" binding:"required"`                        // 开始领取时间
	ReceiveEndTime   time.Time `json:"receive_end_time" binding:"required,gtfield=ReceiveStartTime"` // 领取结束时间
	UseStartTime     time.Time `json:"use_start_time" binding:"required"`                            // 优惠券开始时间
	UseEndTime       time.Time `json:"use_end_time" binding:"required,gtfield=UseStartTime"`         // 优惠券结束时间
	TotalCount       int       `json:"total_count" binding:"required,gte=1"`                         // 优惠券总数
	IsLadder         int       `json:"is_ladder" binding:"omitempty,oneof=0 1"`                      // 是否是多梯度优惠券
	LadderRule       []Ladder  `json:"ladder_rule" binding:"omitempty,required_if=IsLadder 1,dive"`  // 多梯度规则
}

type Ladder struct {
	ThresholdValue int `json:"threshold_value" binding:"required,gt=0"`
	DiscountValue  int `json:"discount_value" binding:"required,gt=0"`
}

type CouponUpdate struct {
	ID int `json:"id" binding:"required"`
	CouponCreate
}

type CouponSearch struct {
	Name             string    `json:"name" binding:"omitempty,min=1,max=32"`                         // 优惠券名称
	ScopeID          int       `json:"usage_scope_id" binding:"omitempty,oneof=1 2 3 4 5"`            // 关联了表的，coupon_usage_scope
	ThresholdType    int       `json:"threshold_type" binding:"omitempty,gt=0,oneof=1 2"`             // 门槛：1 表示金额 2 表示件数
	ThresholdValue   int       `json:"threshold_value" binding:"omitempty,gt=0"`                      // 门槛：200或者3件
	DiscountType     int       `json:"type" binding:"omitempty,oneof=1 2"`                            // 1 = 满减，2 = 折扣
	DiscountValue    int       `json:"discount_value" binding:"omitempty,gte=0"`                      // 优惠券金额 * 100  或者 百分之多少的折扣
	ReceiveLimit     int       `json:"user_receive_limit" binding:"omitempty,gt=0,oneof=1 2"`         // 是否限制领取次数
	ReceiveStartTime time.Time `json:"receive_start_time" binding:"omitempty"`                        // 开始领取时间
	ReceiveEndTime   time.Time `json:"receive_end_time" binding:"omitempty,gtfield=ReceiveStartTime"` // 领取结束时间
	UseStartTime     time.Time `json:"use_start_time" binding:"omitempty"`                            // 优惠券开始时间
	UseEndTime       time.Time `json:"use_end_time" binding:"omitempty,gtfield=UseStartTime"`         // 优惠券结束时间
	Paginate
}

type CouponDelete struct {
	Id int `json:"id" binding:"required,gt=0"`
}
