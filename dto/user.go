package dto

type UserLogin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
