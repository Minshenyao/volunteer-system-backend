package models

import (
	"volunteer-system-backend/config"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

func InitDB() {
	dbConfig := config.ProjectConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Pass, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	// 自动迁移
	err = DB.AutoMigrate(&User{}, &Task{}, &TaskParticipant{}, &Message{})
	if err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}
	err = CreateAdminUser()
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTestTask()
	if err != nil {
		log.Fatal(err)
	}
}

func CreateAdminUser() error {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("无法对密码进行哈希处理")
	}
	// 创建新用户
	user := User{
		Email:         "admin@admin.com",
		Password:      string(hashedPassword),
		Nickname:      "管理员",
		Gender:        "保密",
		Avatar:        "",
		Phone:         "188888888888",
		CreatedAt:     time.Now().Local(),
		Duration:      999,
		Admin:         true,
		LastLoginTime: time.Now().Local(),
	}
	if err := DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
		if err = DB.Create(&user).Error; err != nil {
			return errors.New("无法创建管理员用户")
		}
	}
	return nil
}

func CreateTestTask() error {
	task := Task{
		Name:      "冬至晚会",
		CreatedAt: time.Now().Local(),
		StartTime: time.Now().Local(),
		EndTime:   time.Now().Add(60 * time.Minute).Local(),
		Location:  "大礼堂",
		Limit:     50,
	}
	if err := DB.Where("name = ?", task.Name).First(&Task{}).Error; err != nil {
		if err = DB.Create(&task).Error; err != nil {
			return errors.New("无法创建测试数据")
		}
	}
	return nil
}
