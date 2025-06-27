package dto

// RegisterUserRequest 用户注册请求
type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required" example:"123@qq.com"`
	Nickname string `json:"nickname" binding:"required" example:"user01"`
	Gender   string `json:"gender" binding:"required" example:"男"`
	Phone    string `json:"phone" binding:"required" example:"13812345678"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// LoginUserRequest 用户登录请求
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required" example:"123@qq.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// UserProfileRequest 用户信息请求
type UserProfileRequest struct {
	Authorization string `json:"authorization" binding:"required" example:"Bearer eyJhbGciOiJIsInR5cCI6IkpXJ9..."`
}

// ChangePasswordRequest 更改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// UpdateUserRoleRequest 用户角色更新请求
type UpdateUserRoleRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// UpdateUserInfoRequest 用户信息更新请求
type UpdateUserInfoRequest struct {
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
}
