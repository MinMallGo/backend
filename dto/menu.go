package dto

type MenuList struct {
	Name      string `json:"name"`
	Component string `json:"component"`
}
type MenuSearch struct {
	Paginate
}
type RoleMenu struct {
	ID     int `gorm:"primaryKey;autoIncrement"`
	RoleID int
	MenuID int
}
