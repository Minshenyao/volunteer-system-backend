package controllers

import (
	"volunteer-system-backend/dto"
	"volunteer-system-backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// MessageController 消息控制器
type MessageController struct {
	messageService services.MessageService
}

// NewMessageController 创建新的消息控制器实例
func NewMessageController(messageService services.MessageService) *MessageController {
	return &MessageController{
		messageService: messageService,
	}
}

// GetMessagesByUserID 根据用户ID获取消息列表
func (c *MessageController) GetMessagesByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户 ID"})
		return
	}

	messages, err := c.messageService.GetMessagesByUserID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取消息"})
		return
	}

	var messageResponses []dto.MessageResponse
	for _, message := range messages {
		messageResponse := dto.MessageResponse{
			ID:      message.ID,
			Title:   message.Title,
			Content: message.Content,
			Time:    message.Time.Format("2006-01-02 15:04"),
			Status:  message.Status,
		}
		messageResponses = append(messageResponses, messageResponse)
	}

	ctx.JSON(http.StatusOK, messageResponses)
}

// CreateMessage 创建新消息
func (c *MessageController) CreateMessage(ctx *gin.Context) {
	var request struct {
		UserID  uint   `json:"userID" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := c.messageService.CreateMessage(request.UserID, request.Title, request.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建消息失败"})
		return
	}

	messageResponse := dto.MessageResponse{
		ID:      message.ID,
		Title:   message.Title,
		Content: message.Content,
		Time:    message.Time.Format("2006-01-02 15:04"),
		Status:  message.Status,
	}

	ctx.JSON(http.StatusCreated, messageResponse)
}

// MarkMessageAsRead 将消息标记为已读
func (c *MessageController) MarkMessageAsRead(ctx *gin.Context) {
	messageIDStr := ctx.Param("messageID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "消息 ID 无效"})
		return
	}

	err = c.messageService.MarkMessageAsRead(uint(messageID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "无法将消息标记为已读"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "标记为已读的消息"})
}
