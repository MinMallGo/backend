package service

import (
	"context"
	"encoding/json"
	"mall_backend/dao"
	"mall_backend/dao/model"
	"mall_backend/util"
	"strconv"
)

func Menu() {
	data, err := dao.NewMenuDao(util.DBClient()).Search()
	dataRole, err := dao.NewRoleDao(util.DBClient()).SearchRole()
	if err != nil {
		return
	}
	menuMap := make(map[int]model.MmMenu)
	for _, menu := range data {
		menuMap[int(menu.ID)] = menu
	}
	roleMenuMap := make(map[int][]model.MmMenu)
	for _, roleMenu := range dataRole {
		if menu, ok := menuMap[int(roleMenu.MenuID)]; ok {
			roleMenuMap[int(roleMenu.RoleID)] = append(roleMenuMap[int(roleMenu.RoleID)], menu)
		}
	}
	ctx := context.Background()
	cacheInit := util.CacheClient()
	for i, menu := range roleMenuMap {
		menuJSON, err := json.Marshal(menu)
		if err != nil {
			return
		}
		cacheInit.Set(ctx, strconv.Itoa(i), menuJSON, 0)
	}
	return
}
