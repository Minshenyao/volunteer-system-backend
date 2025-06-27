package utils

import (
	"github.com/gin-gonic/gin"
)

// Response 通用 API 响应格式
type Response struct {
	Status  string      `json:"status"`  // 状态描述
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 数据部分
}

// Respond 通用响应方法
func Respond(c *gin.Context, code int, status string, message string, data interface{}) {
	c.JSON(code, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
}
