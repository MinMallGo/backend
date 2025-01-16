package dto

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserLogout struct{}

type UserLogoutResponse struct{}

type UserRegister struct {
	Account        string `json:"account" binding:"required,alphanum,min=1,max=32"`
	Username       string `json:"username" binding:"required,min=2,max=16"`
	Password       string `json:"password" binding:"required,min=8,max=32"`
	RepeatPassword string `json:"repeat_password" binding:"required,eqfield=Password"`
	Sex            int    `json:"sex" binding:"required,oneof=0 1"`
	Birthday       string `json:"birthday" binding:"required,datetime=2006-01-02"`
	Phone          string `json:"phone" binding:"omitempty,regexpPhone"`
}

type UserRegisterResponse struct {
	Account  string `json:"account"`
	Username string `json:"username"`
}
