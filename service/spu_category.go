package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	constants "mall_backend/constant"
	"mall_backend/dao/model"
	"mall_backend/dto"
	"mall_backend/response"
	"mall_backend/util"
	"strconv"
	"strings"
	"time"
)

// CategoryCreate 创建spu分类
func CategoryCreate(c *gin.Context, create *dto.CategoryCreate) {
	// 查询父级的id，并组装path
	levelPath, err := genCategoryPath(create.ParentId)
	if err != nil {
		response.Error(c, err)
		return
	}

	param := &model.MmCategory{
		Name:       create.Name,
		ParentID:   int32(create.ParentId),
		Code:       util.GenSalt(),
		Path:       levelPath,
		Status:     constants.NormalStatus,
		CreateTime: time.Now(),
		UpdateTime: util.MinDateTime(),
		DeleteTime: util.MinDateTime(),
	}

	tx := util.DBClient().Model(&model.MmCategory{}).Create(&param)
	if tx.Error != nil {
		response.Failure(c, "创建失败")
		return
	}

	response.Success(c, []string{})
	return
}

// CategoryDelete 删除spu分类
func CategoryDelete(c *gin.Context, delete *dto.CategoryDelete) {
	update := &model.MmCategory{
		Status:     constants.BanStatus,
		DeleteTime: time.Now(),
	}
	log.Println(update)

	tx := util.DBClient().Model(&model.MmCategory{}).Select("status", "delete_time").Where("id = ?", delete.Id).Updates(update)
	if tx.Error != nil {
		response.Failure(c, "删除失败")
		return
	}

	response.Success(c, []string{})
	return
}

// CategoryUpdate 修改spu的分类信息
func CategoryUpdate(c *gin.Context, update *dto.CategoryUpdate) {
	// 额 == ，这里似乎要重新组织一下path，因为有可能改了parent_id
	update.UpdateTime = time.Now()
	tx := util.DBClient().Model(&model.MmCategory{}).Where("id = ?", update.Id).Updates(update)
	if tx.Error != nil {
		response.Failure(c, tx.Error.Error())
		return
	}

	response.Success(c, []string{})
	return
}

// CategorySearch 查询
func CategorySearch(c *gin.Context, search *dto.CategorySearch) {
	// id, name, page, limit
	category := &[]model.MmCategory{}

	// 先写吧，优化等以后再来说

	whereStr := "status = ?"
	var param []interface{}
	param = append(param, constants.NormalStatus)

	if search.Id > 0 {
		whereStr += " AND id = ?"
		param = append(param, search.Id)
	}

	if len(search.Name) > 0 {
		whereStr += " AND name LIKE ?"
		param = append(param, "%"+search.Name+"%")
	}

	tx := util.DBClient().Model(&model.MmCategory{}).Offset((search.Page-1)*search.Limit).Limit(search.Limit).Where(whereStr, param...).Find(&category)
	if tx.Error != nil {
		response.Failure(c, tx.Error.Error())
		return
	}

	response.Success(c, category)
	return
}

// 根据提供的parent_id，构造
func genCategoryPath(parentId int) (string, error) {
	// 查询父级的id，并组装path
	levelPath := ""
	if parentId != 0 {
		res := &model.MmCategory{}
		tx := util.DBClient().Model(&model.MmCategory{}).Where("id = ? AND status = ?", parentId, constants.NormalStatus).First(&res)
		if tx.Error != nil {
			return levelPath, errors.New("请选择正确的父级数据")
		}

		res.Path = strings.TrimSpace(res.Path)
		if len(res.Path) > 1 {
			levelPath = res.Path + "," + strconv.Itoa(int(res.ID))
		}
	}

	return levelPath, nil
}

// CategoryExists 根据ID查询当前分类是否存在。 我建议让他们自己判断一下，如果是0的话就不要传过来
func CategoryExists(categoryID int) bool {
	tx := util.DBClient().Model(&model.MmCategory{}).
		Select("id").
		Where("status = ?", constants.NormalStatus).
		Where("id = ?", categoryID).
		First(&model.MmCategory{})
	return tx.RowsAffected != 0
}
