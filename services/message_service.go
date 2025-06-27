package services

import (
	"volunteer-system-backend/models"
	"gorm.io/gorm"
	"log"
	"time"
)

// MessageService 消息服务接口
type MessageService interface {
	GetMessagesByUserID(userID uint) ([]models.Message, error)
	CreateMessage(userID uint, title, content string) (*models.Message, error)
	MarkMessageAsRead(messageID uint) error
}

// messageServiceImpl 消息服务实现
type messageServiceImpl struct {
	db *gorm.DB
}

// NewMessageService 创建新的消息服务实例
func NewMessageService() MessageService {
	if models.DB == nil {
		log.Println("数据库连接未初始化")
	}
	return &messageServiceImpl{
		db: models.DB,
	}
}

// GetMessagesByUserID 根据用户ID获取消息列表
func (s *messageServiceImpl) GetMessagesByUserID(userID uint) ([]models.Message, error) {
	var messages []models.Message
	err := s.db.Where("user_id = ?", userID).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// CreateMessage 创建新消息
func (s *messageServiceImpl) CreateMessage(userID uint, title, content string) (*models.Message, error) {
	message := &models.Message{
		UserID:  userID,
		Title:   title,
		Content: content,
		Time:    time.Now(),
		Status:  "unread",
	}
	err := s.db.Create(message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}

// MarkMessageAsRead 将消息标记为已读
func (s *messageServiceImpl) MarkMessageAsRead(messageID uint) error {
	err := s.db.Model(&models.Message{}).Where("id = ?", messageID).Update("status", "read").Error
	if err != nil {
		return err
	}
	return nil
}
