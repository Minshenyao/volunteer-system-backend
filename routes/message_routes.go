package routes

import (
	"volunteer-system-backend/controllers"
	"volunteer-system-backend/middlewares"
	"volunteer-system-backend/services"
	"github.com/gin-gonic/gin"
)

// MessageRoutes 消息相关的路由
func MessageRoutes(r *gin.Engine) {
	//// 配置跨域中间件
	//r.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"http://localhost:5173"}, // 前端地址
	//	AllowMethods:     []string{"POST", "GET", "OPTIONS"},
	//	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	//	AllowCredentials: true,
	//}))
	messageService := services.NewMessageService()
	messageController := controllers.NewMessageController(messageService)

	messageGroup := r.Group("/message")
	messageGroup.Use(middlewares.AuthMiddleware()) // 应用 JWT 中间件
	{
		messageGroup.GET("/:userID", messageController.GetMessagesByUserID)
		messageGroup.POST("/", messageController.CreateMessage)
		messageGroup.PUT("/:messageID/read", messageController.MarkMessageAsRead)
	}
}
