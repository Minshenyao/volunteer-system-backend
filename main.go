package main

import (
	_ "volunteer-system-backend/docs"

	"volunteer-system-backend/config"
	"volunteer-system-backend/models"
	"volunteer-system-backend/routes"
)

// @title 志愿者系统
// @version 1.0
// @description 一个简单的志愿者管理系统，包括用户注册、登录和志愿任务的创建、删除、修改、查找等

// @contact.name Minshenyao
// @contact.url https://www.github.com/Minshenyao

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
func main() {
	//加载配置
	config.LoadConfig()
	//初始化数据库
	models.InitDB()
	//初始化路由
	router := routes.SetupRouter()
	routes.MessageRoutes(router)

	router.Run(":8080")
}
