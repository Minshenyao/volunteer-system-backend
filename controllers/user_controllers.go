package controllers

import (
	"volunteer-system-backend/dto"
	"volunteer-system-backend/services"
	"volunteer-system-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterUser 注册新用户
// @Summary 注册新用户
// @Description 用户注册
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.RegisterUserRequest true "注册信息"
// @Router /api/register [post]
func RegisterUser(c *gin.Context) {
	var input dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}

	// 调用服务层
	if err := services.RegisterUser(input.Email, input.Nickname, input.Gender, input.Phone, input.Password); err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "注册失败"+err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "注册成功", nil)
}

// LoginUser 登录用户
// @Summary 登录用户
// @Description 用户登录并生成 token
// @Tags user
// @Accept json
// @Produce json
// @Param request body dto.LoginUserRequest true "登录信息"
// @Router /api/login [post]
func LoginUser(c *gin.Context) {
	var input dto.LoginUserRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}

	token, err := services.LoginUser(input.Email, input.Password)
	if err != nil {
		utils.Respond(c, http.StatusUnauthorized, "error", "登录失败"+err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "登录成功", gin.H{"token": token})
}

// GetUserProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的个人信息
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Router /user/profile [get]
func GetUserProfile(c *gin.Context) {
	userName, exists := c.Get("Email")
	if !exists {
		utils.Respond(c, http.StatusUnauthorized, "error", "用户未登录", nil)
		return
	}
	_, exists = c.Get("LoginTime")
	if !exists {
		utils.Respond(c, http.StatusUnauthorized, "error", "用户未登录", nil)
		return
	}

	user, err := services.GetUserProfile(userName)
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "error", err.Error(), nil)
	}

	// 判断用户状态
	status := utils.GetUserStatus(user.LastLoginTime)

	utils.Respond(c, http.StatusOK, "success", "获取成功", gin.H{
		"email":        user.Email,
		"nickName":     user.Nickname,
		"gender":       user.Gender,
		"avatar":       user.Avatar,
		"volunteerID":  user.ID,
		"phone":        user.Phone,
		"duration":     user.Duration,
		"isAdmin":      user.Admin,
		"status":       status,
		"lastActivity": user.LastLoginTime,
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改用户密码
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Param request body dto.ChangePasswordRequest true "修改密码信息"
// @Router /user/change_password [put]
func ChangePassword(c *gin.Context) {
	var input dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}

	email, exists := c.Get("Email")
	if !exists {
		utils.Respond(c, http.StatusUnauthorized, "error", "用户未登录", nil)
		return
	}

	if email == "admin@admin.com" {
		utils.Respond(c, http.StatusForbidden, "error", "不允许修改默认管理员用户密码", nil)
		return
	}

	if err := services.ChangePassword(email.(string), input.OldPassword, input.NewPassword); err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "修改密码失败："+err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "success", "修改密码成功", nil)
}

// GetVolunteerCount 统计志愿者用户个数
// @Summary 统计志愿者用户个数
// @Description 统计所有志愿者用户个数，返回姓名、服务时长、状态、最近一次登录时间
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Router /user/volunteer_count [get]
func GetVolunteerCount(c *gin.Context) {
	volunteers, err := services.GetVolunteerCount()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "获取志愿者信息失败："+err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "success", "获取志愿者信息成功", gin.H{"volunteers": volunteers})
}

// UploadAvatar 上传用户头像
// @Summary 上传用户头像
// @Description 用户上传头像
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Param avatar formData file true "用户头像文件"
// @Router /user/upload_avatar [post]
func UploadAvatar(c *gin.Context) {
	// 获取上传的文件
	file, fileHeader, err := c.Request.FormFile("avatar")
	if err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误：无法获取上传的文件", nil)
		return
	}
	defer file.Close()

	email, exists := c.Get("Email")
	if !exists {
		utils.Respond(c, http.StatusUnauthorized, "error", "用户未登录", nil)
		return
	}

	// 调用服务层
	if err := services.UploadAvatar(email.(string), file, fileHeader); err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "上传头像失败："+err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "success", "上传头像成功", nil)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户的个人信息
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Param request body dto.UpdateUserInfoRequest true "更新用户信息"
// @Router /user/update_user_info [put]
func UpdateUserInfo(c *gin.Context) {
	var input dto.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}
	email, exists := c.Get("Email")
	if !exists {
		utils.Respond(c, http.StatusUnauthorized, "error", "用户未登录", nil)
		return
	}
	if err := services.UpdateUserInfo(email, input.Nickname, input.Gender, input.Phone); err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "更新用户信息失败："+err.Error(), nil)
		return
	}

	utils.Respond(c, http.StatusOK, "success", "更新用户信息成功", nil)
}
