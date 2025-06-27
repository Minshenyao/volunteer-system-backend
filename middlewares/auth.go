package middlewares

import (
	"volunteer-system-backend/models"
	"volunteer-system-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "非法的请求"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		// 将用户信息保存到上下文
		c.Set("Email", claims.Email)
		c.Set("Nickname", claims.Nickname)
		c.Set("LoginTime", claims.LoginTime.Format("2006-01-02 15:04:05"))
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 首先验证 JWT token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "非法的请求"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 2. 查询用户信息，验证管理员权限
		var user models.User
		result := models.DB.Where("email = ?", claims.Email).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 3. 检查用户是否是管理员
		if !user.Admin {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		// 4. 将用户信息保存到上下文
		c.Set("Email", claims.Email)
		c.Set("Nickname", claims.Nickname)
		c.Set("LoginTime", claims.LoginTime.Format("2006-01-02 15:04:05"))
		c.Set("IsAdmin", true)

		c.Next()
	}
}
