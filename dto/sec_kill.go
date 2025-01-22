package dto

// SecKillCreate 秒杀活动的添加的数据
type SecKillCreate struct {
	ActiveCreate
	SecKillInfo
}

// SecKillInfo 秒杀基本信息
type SecKillInfo struct {
	SpuId            int    `json:"spu_id" binding:"required,gt=0"`
	SkuId            int    `json:"sku_id" binding:"required,gt=0"`
	SpecId           int    `json:"spec_id" binding:"omitempty,gt=0"`
	Stock            int    `json:"stock" binding:"required,gt=0"`
	SecKillStartTime string `json:"sec_start_time" binding:"required,datetime=2006-01-02 15:04:05"`
	SecKillEndTime   string `json:"sec_end_time" binding:"required,datetime=2006-01-02 15:04:05,gtefield=SecKillStartTime"`
	Price            int    `json:"price" binding:"omitempty,gt=0"`
}

// SecKillUpdate 秒杀活动更新数据
type SecKillUpdate struct {
	ActiveCreate
	SpuId            int    `json:"spu_id" binding:"omitempty,gt=0"`
	SkuId            int    `json:"sku_id" binding:"omitempty,gt=0"`
	Stock            int    `json:"stock" binding:"omitempty,gt=0"`
	SecKillStartTime string `json:"sec_start_time" binding:"omitempty,datetime=2006-01-02 15:04:05"`
	SecKillEndTime   string `json:"sec_end_time" binding:"omitempty,datetime=2006-01-02 15:04:05,gtefield=SecKillStartTime"`
	Price            int    `json:"price" binding:"omitempty,gt=0"`
}
