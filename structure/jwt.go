package structure

type JWTUserInfo struct {
	Username string `json:"username" binding:"required"`
	UserID   int    `json:"userid" binding:"required"`
	Role     string `json:"role" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}
