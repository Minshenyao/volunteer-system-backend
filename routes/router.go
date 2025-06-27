package routes

import (
	"volunteer-system-backend/controllers"
	"volunteer-system-backend/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	//// 配置跨域中间件
	//r.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"http://47.122.29.1:5173"}, // 前端地址
	//	AllowMethods:     []string{"POST", "GET", "OPTIONS"},
	//	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	//	AllowCredentials: true,
	//}))
	// CORS 中间件配置
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"} // 前端地址
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	// Swagger API 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 公共 API 路由（无需鉴权）
	api := r.Group("/api")
	{
		api.POST("/register", controllers.RegisterUser) //用户注册
		api.POST("/login", controllers.LoginUser)       //用户登录
	}

	user := r.Group("/user")
	user.Use(middlewares.AuthMiddleware())
	{
		user.GET("/profile", controllers.GetUserProfile)            // 获取用户信息
		user.GET("/volunteer_count", controllers.GetVolunteerCount) // 统计志愿者用户个数
		user.POST("/upload_avatar", controllers.UploadAvatar)       // 新增上传头像接口
		user.PUT("/change_password", controllers.ChangePassword)    // 修改密码
		user.PUT("/update_profile", controllers.UpdateUserInfo)     // 更新用户信息
	}
	// 需要 JWT 鉴权的路由
	task := r.Group("/task")
	task.Use(middlewares.AuthMiddleware()) // 应用 JWT 中间件
	{
		task.GET("/tasks", controllers.GetTasks)                // 获取志愿活动列表
		task.GET("/getTaskStatus", controllers.GetTaskStatus)   // 获取任务状态
		task.GET("/getTaskDetails", controllers.GetTaskDetails) // 获取任务详情
		task.POST("/join", controllers.JoinTask)                // 参加志愿者活动
	}

	//管理员特定操作路由
	admin := r.Group("/admin")
	admin.Use(middlewares.AdminMiddleware()) // 应用管理员中间件
	{
		admin.GET("/volunteer_count", controllers.GetVolunteerCount)      // 统计志愿者用户个数
		admin.GET("/profile", controllers.GetUserProfile)                 // 获取用户信息
		admin.GET("/tasks", controllers.GetTasks)                         // 获取志愿活动列表
		admin.GET("/getTaskDetails", controllers.GetTaskDetails)          // 获取任务详情
		admin.GET("/getTaskStatus", controllers.GetTaskStatus)            // 获取任务状态
		admin.POST("/GetTaskAuditDetail", controllers.GetTaskAuditDetail) //获取任务报名详情
		admin.POST("/approveVolunteer", controllers.ApproveVolunteer)     //审核通过
		admin.POST("/rejectVolunteer", controllers.RejectVolunteer)       //审核拒绝
		admin.POST("/create_task", controllers.CreateTask)                // 创建志愿活动
		admin.POST("/update", controllers.UpdateTask)                     // 修改志愿活动
		admin.POST("/join", controllers.JoinTask)                         // 参加志愿者活动
		admin.DELETE("/delete_task", controllers.DeleteTask)              // 删除志愿活动
	}

	return r
}
