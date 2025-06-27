package models

import (
	"time"
)

// Task 志愿活动
type Task struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`               // 活动名称
	CreatedAt time.Time `gorm:"type:datetime;not null"` // 活动创建时间
	StartTime time.Time `gorm:"type:datetime;not null"` // 活动开始时间
	EndTime   time.Time `gorm:"type:datetime;not null"` // 活动结束时间
	Location  string    `gorm:"size:255"`               // 活动举行地点
	Limit     uint      `gorm:"not null;default:0"`     // 限制人数
	Joined    uint      `gorm:"default:0"`              // 已参加人数
	// 新增关联关系
	Participants []TaskParticipant `gorm:"foreignKey:TaskID"`
}

type Participant struct {
	Email         string `gorm:"unique, not null"`                 // 邮箱，用于登录
	Nickname      string `gorm:"not null,default:null"`            // 姓名
	Gender        string `gorm:"not null,default:null"`            // 性别
	Avatar        string `gorm:"default:null"`                     // 头像地址
	Phone         string `gorm:"default:null"`                     // 联系方式(默认+86)
	Duration      uint   `gorm:"default:0"`                        // 志愿时长（分钟）
	LastLoginTime string `json:"lastLoginTime" binging:"required"` // 最近一次登录时间
}
