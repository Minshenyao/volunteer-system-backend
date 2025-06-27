package models

import "gorm.io/gorm"

// TaskParticipant 表示任务的已参加人员信息
type TaskParticipant struct {
	gorm.Model
	TaskID   uint   `gorm:"not null"`           // 关联的任务ID
	Nickname string `gorm:"not null"`           // 参加人员的用户名
	Email    string `gorm:"not null"`           //参加人员的邮箱
	Status   uint   `gorm:"not null default:3"` // 0表示待审核 1表示审核通过 2表示审核不通过 3表示未参加
}
