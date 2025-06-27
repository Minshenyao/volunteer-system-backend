package services

import (
	"volunteer-system-backend/models"
	"volunteer-system-backend/utils"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"mime/multipart"
	"time"
)

// RegisterUser 注册用户服务
func RegisterUser(email, nickname, gender, phone, password string) error {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := models.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.New("该用户已经存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("无法对密码进行哈希处理")
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// 创建新用户
	user := models.User{
		Email:         email,
		Nickname:      nickname,
		Password:      string(hashedPassword),
		CreatedAt:     time.Now().In(loc),
		Duration:      0,
		Admin:         false,
		LastLoginTime: time.Now().In(loc),
		Gender:        gender,
		Phone:         phone,
	}
	if err := models.DB.Create(&user).Error; err != nil {
		return errors.New("无法创建用户")
	}
	return nil
}

// LoginUser 用户登录服务
func LoginUser(email, password string) (string, error) {
	var user models.User

	// 查找用户
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("邮箱或密码无效")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("邮箱或密码无效")
	}

	localTime := time.Now().Local()
	// 生成 JWT
	token, err := utils.GenerateJWT(user.Email, user.Nickname, localTime)
	if err != nil {
		return "", errors.New("生成 token 失败")
	}
	user.LastLoginTime = localTime
	if err := models.DB.Save(&user).Error; err != nil {
		return "", errors.New("无法更新最后登录时间")
	}
	return token, nil
}

// GetUserIDByEmail 通过 email 查询用户 ID
func GetUserIDByEmail(email string) (uint, error) {
	var user models.User
	// 根据 email 查找用户
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("用户不存在")
		}
		return 0, err
	}
	return user.ID, nil
}

// GetUserProfile 获取当前登录用户的个人信息
func GetUserProfile(email any) (*models.User, error) {
	var user models.User
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

// ChangePassword 修改密码服务
func ChangePassword(email, oldPassword, newPassword string) error {
	var user models.User
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("无法对密码进行哈希处理")
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err := models.DB.Save(&user).Error; err != nil {
		return errors.New("无法更新密码")
	}

	return nil
}

// GetVolunteerCount 统计志愿者用户个数服务
func GetVolunteerCount() ([]map[string]interface{}, error) {
	var volunteers []map[string]interface{}
	if err := models.DB.Table("users").Select("email, id, nickname, gender, phone, avatar, duration, last_login_time").Find(&volunteers).Error; err != nil {
		return nil, err
	}
	// 判断用户状态
	for i, volunteer := range volunteers {
		// 确保 last_login_time 是 time.Time 类型
		lastLoginTime, ok := volunteer["last_login_time"].(time.Time)
		if !ok {
			return nil, errors.New("无法将 last_login_time 转换为 time.Time 类型")
		}
		// 获取用户状态
		status := utils.GetUserStatus(lastLoginTime)
		volunteers[i]["status"] = status
		formattedTime := lastLoginTime.Format("2006-01-02 15:04:05")
		volunteers[i]["last_login_time"] = formattedTime
	}

	return volunteers, nil
}

// UploadAvatar 上传用户头像
func UploadAvatar(email string, file multipart.File, fileHeader *multipart.FileHeader) error {
	var user models.User
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}
	// 上传文件到图床
	imageURL, err := utils.AliOSSUtils(file, fileHeader)
	if err != nil {
		return err
	}
	// 更新用户的头像字段
	user.Avatar = imageURL
	if err := models.DB.Save(&user).Error; err != nil {
		return errors.New("无法更新用户头像")
	}
	return nil
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(email any, nickname, gender, phone string) error {
	var user models.User
	if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 只在字段值发生变化时进行更新
	updated := false

	if nickname != "" && nickname != user.Nickname {
		user.Nickname = nickname
		updated = true
	}
	if gender != "" && gender != user.Gender {
		user.Gender = gender
		updated = true
	}
	if phone != "" && phone != user.Phone {
		user.Phone = phone
		updated = true
	}

	// 如果有字段更新，则保存
	if updated {
		if err := models.DB.Save(&user).Error; err != nil {
			return errors.New("无法更新用户信息")
		}
	}

	return nil
}
