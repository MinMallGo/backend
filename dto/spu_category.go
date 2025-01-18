package dto

import (
	"time"
)

type CategoryCreate struct {
	Name     string `json:"name" binding:"required,min=1,max=32"`
	ParentId int    `json:"parent_id" binding:"gte=0"`
	Path     string `json:"path"`
}

type CategoryDelete struct {
	Id int `json:"id"`
}

type CategoryUpdate struct {
	CategoryDelete
	CategoryCreate
	UpdateTime time.Time
}

type CategorySearch struct {
	Id    int    `json:"id" binding:"omitempty,gt=0"`
	Name  string `json:"name"`
	Page  int    `json:"page" binding:"gt=0"`
	Limit int    `json:"limit"`
}
