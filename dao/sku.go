package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/util"
	"strings"
	"time"
)

type SkuDao struct {
	db *gorm.DB
}

func NewSkuDao(db *gorm.DB) *SkuDao {
	return &SkuDao{db: db}
}

func (d *SkuDao) Create(create *dto.SkuCreate) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		title, specs, err := NewSpecKeyDao(d.db).GenSkuData(&create.Spec)
		if err != nil {
			return err
		}

		// title = 规格值的名字
		// specs = 规格上下级的关系
		param := &model.MmSku{
			Title:      title,
			SpuID:      int32(create.SpuID),
			Price:      int32(create.Price),
			Stock:      int32(create.Stock),
			Spces:      string(specs),
			Status:     constants.NormalStatus,
			CreateTime: time.Now(),
			UpdateTime: util.MinDateTime(),
			DeleteTime: util.MinDateTime(),
		}
		err = d.db.Model(&model.MmSku{}).Create(&param).Error
		if err != nil {
			return err
		}

		err = SpecCreate(create.SpuID, int(param.ID), &create.Spec)
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *SkuDao) Update(update *dto.SkuUpdate) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		title, specs, err := NewSpecKeyDao(d.db).GenSkuData(&update.Spec)
		if err != nil {
			return err
		}

		// title = 规格值的名字
		// specs = 规格上下级的关系
		param := &model.MmSku{
			Title:      title,
			SpuID:      int32(update.SpuID),
			Price:      int32(update.Price),
			Stock:      int32(update.Stock),
			Spces:      string(specs),
			UpdateTime: time.Now(),
		}
		// 更新sku
		err = d.db.Model(&model.MmSku{}).Where("id = ?", update.Id).Updates(param).Error
		if err != nil {
			log.Println("update sku with error:", err)
			return err
		}
		// 更新spec
		err = SpecUpdate(update.SpuID, update.SkuID, &update.Spec)
		if err != nil {
			log.Println("updateSpec  with error:", err)
			return err
		}

		return nil

	})
}

func (d *SkuDao) Exists(id ...int) bool {
	if len(id) == 0 {
		return false
	}
	return d.db.
		Model(&model.MmSku{}).
		Select("id").
		Where("status = ?", constants.NormalStatus).
		Where("id in ?", id).
		Find(&[]model.MmCategory{}).
		RowsAffected == int64(len(id))
}

type ExistsAndStock struct {
	ID    int
	Stock int
}

// Enough 规格存在且存量足够
func (d *SkuDao) Enough(skus []ExistsAndStock) error {
	if len(skus) == 0 {
		return errors.New("创建订单失败：商品不能为空")
	}
	// 构造 SQL 条件
	conditions := make([]string, 0, len(skus))
	values := make([]interface{}, 0, len(skus)*2)
	for _, s := range skus {
		conditions = append(conditions, "(id = ? AND stock >= ? AND status = ?)")
		values = append(values, s.ID, s.Stock, constants.NormalStatus)
	}

	// 拼接成 WHERE 子句
	whereClause := strings.Join(conditions, " OR ")

	var count int64
	tx := d.db.Model(&model.MmSku{}).Where(whereClause, values...).Count(&count)
	if tx.Error != nil {
		return tx.Error
	}

	if count != int64(len(skus)) {
		return errors.New("商品或库存不足")
	}
	return nil
}

func (d *SkuDao) Delete(id ...int) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		err := d.db.Model(&model.MmSku{}).Where("id in ?", id).Updates(map[string]interface{}{
			"status":      constants.BanStatus,
			"delete_time": time.Now().Format("2006-01-02 15:04:05"),
		}).Error
		if err != nil {
			return err
		}

		err = d.db.Model(&model.MmSpec{}).Where("sku_id in ?", id).Updates(map[string]interface{}{
			"status":      constants.BanStatus,
			"delete_time": time.Now().Format("2006-01-02 15:04:05"),
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *SkuDao) OneById(id int) (*model.MmSku, error) {
	res := &model.MmSku{}
	if err := d.db.Model(&model.MmSku{}).Where("status = ?", constants.NormalStatus).Where("spu_id = ?").First(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (d *SkuDao) More(search *dto.SkuSearch) (*dto.PaginateCount, error) {
	// TODO
	return nil, nil
}

func (d *SkuDao) MoreWithOrder(items *[]dto.ShoppingItem) (*[]model.MmSku, error) {
	res := &[]model.MmSku{}
	skuIds := make([]int, 0, len(*items))
	for _, item := range *items {
		skuIds = append(skuIds, item.SkuID)
	}

	tx := d.db.Model(&model.MmSku{}).Where("id in ?", skuIds).Where("status = ?", true).Find(&res)
	if tx.Error != nil {
		return res, tx.Error
	}

	return res, nil
}

// TODO 加库存减库存这里需要优化，不过以后再说
// TODO 改造成一个sql。用case when then

func (d *SkuDao) DecreaseStock(items *[]dto.ShoppingItem) bool {
	decr := make([]StockUpdate, 0, len(*items))
	for _, item := range *items {
		decr = append(decr, StockUpdate{
			ID:  item.SkuID,
			Num: item.Num,
		})
	}
	// 改造成使用 case when then end 来一条sql执行
	err := d.StockDecrease(&decr)
	if err != nil {
		return false
	}

	return true
}

func (d *SkuDao) IncreaseStock(items *[]model.MmSubOrder) error {
	incr := make([]StockUpdate, 0, len(*items))
	for _, item := range *items {
		incr = append(incr, StockUpdate{
			ID:  int(item.SkuID),
			Num: int(item.Nums),
		})
	}
	err := d.StockIncrease(&incr)
	if err != nil {
		return err
	}
	return nil
}

type StockUpdate struct {
	ID  int
	Num int
}

func (d *SkuDao) StockDecrease(u *[]StockUpdate) error {
	// update mm_sku set stock = case when id = 1 then stock - 1 when id = 2 then stock - 2 end where id in (1,2)
	if len(*u) == 0 {
		return errors.New("商品库存扣减失败")
	}
	str := " WHEN id = %d THEN stock - %d "
	when := ""
	orStr := "(id = %d AND stock >= %d)"
	where := ""
	for index, item := range *u {
		when += fmt.Sprintf(str, item.ID, item.Num)
		where += fmt.Sprintf(orStr, item.ID, item.Num)
		if index < len(*u)-1 {
			where += "OR"
		}
	}

	sql := `UPDATE mm_sku SET stock = CASE %s END WHERE %s`
	sql = fmt.Sprintf(sql, when, where)
	//res := &model.MmSku{}
	tx := d.db.Raw(sql).Exec(sql)

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != int64(len(*u)) {
		return errors.New("商品库存扣减失败")
	}

	return nil
}

func (d *SkuDao) StockIncrease(u *[]StockUpdate) error {
	if len(*u) == 0 {
		return nil
	}
	str := " WHEN id = %d THEN stock + %d "
	when := ""
	idx := ""
	for index, update := range *u {
		when += fmt.Sprintf(str, update.ID, update.Num)
		idx += fmt.Sprintf("%d", update.ID)
		if index < len(*u)-1 {
			idx += ","
		}
	}
	sql := `UPDATE mm_sku SET stock = CASE %s END WHERE id IN (%s) `
	res := &model.MmSku{}
	tx := d.db.Raw(fmt.Sprintf(sql, when, idx)).Scan(res)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// SecKillCreate 扣减库存
func (d *SkuDao) SecKillCreate(items *[]dto.SecKillCreate) error {
	if len(*items) == 0 {
		return errors.New("创建秒杀活动：预扣减库存失败")
	}

	decr := make([]StockUpdate, 0, len(*items))
	for _, item := range *items {
		decr = append(decr, StockUpdate{
			ID:  item.SkuID,
			Num: item.Stock,
		})
	}
	return d.StockDecrease(&decr)
}

func (d *SkuDao) SecKillUpdate(id int, items *[]dto.SecKillCreate) error {
	if len(*items) == 0 {
		return errors.New("更新秒杀活动失败")
	}
	
	//  还需要一个前置条件就是，把原来的查询出来，然后进行库存增加操作，再然后才是对库存进行减操作
	ex := &[]model.MmSeckillProduct{}
	tx := d.db.Model(&model.MmSeckillProduct{}).Where("seckill_id = ? AND status = ?", id, true).Find(ex)
	if tx.Error != nil {
		return tx.Error
	}

	if len(*ex) > 0 {
		incr := make([]StockUpdate, 0, len(*items))
		// 进行库存归还操作
		for _, product := range *ex {
			incr = append(incr, StockUpdate{
				ID:  int(product.SkuID),
				Num: int(product.Stock),
			})
		}
		err := d.StockIncrease(&incr)
		if err != nil {
			return errors.New("更新秒杀活动失败：归还原有商品库存失败")
		}
	}

	decr := make([]StockUpdate, 0, len(*items))
	for _, item := range *items {
		decr = append(decr, StockUpdate{
			ID:  item.SkuID,
			Num: item.Stock,
		})
	}
	return d.StockDecrease(&decr)
}

// TODO 还需要一张库存变动信息表应该。
