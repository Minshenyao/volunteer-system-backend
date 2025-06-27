package models

import (
	"time"
)

// Message 表示用户消息
type Message struct {
	ID      uint      `gorm:"primaryKey"`
	UserID  uint      `gorm:"not null"`                  // 关联的用户ID
	Title   string    `gorm:"not null"`                  // 消息标题
	Content string    `gorm:"not null"`                  // 消息内容
	Time    time.Time `gorm:"not null"`                  // 消息时间
	Status  string    `gorm:"not null;default:'unread'"` // 消息状态，默认未读
}
