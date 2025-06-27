package models

import (
	"time"
)

type User struct {
	ID            uint      `gorm:"primaryKey"`
	Email         string    `gorm:"unique"`                //邮箱
	Nickname      string    `gorm:"not null"`              // 姓名
	Gender        string    `gorm:"not null,default:null"` // 性别
	Password      string    `gorm:"not null"`              // 密码
	Avatar        string    `gorm:"default:null"`          // 头像地址
	Phone         string    `gorm:"default:null"`          // 联系方式手机号码（默认+86）
	CreatedAt     time.Time // 注册时间
	Duration      uint      `gorm:"default:0"`     // 志愿时长（分钟）
	Admin         bool      `gorm:"default:false"` // 是否管理员
	LastLoginTime time.Time // 最近一次登录时间
}
