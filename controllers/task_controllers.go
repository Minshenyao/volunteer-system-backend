package controllers

import (
	"volunteer-system-backend/dto"
	"volunteer-system-backend/services"
	"volunteer-system-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// JoinTask 加入任务
// @Summary 加入任务
// @Description 用户加入任务
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Param task body dto.TaskInfo true "任务信息"
// @Router /task/join [post]
func JoinTask(c *gin.Context) {
	nickname, exists := c.Get("Nickname")
	if !exists {
		utils.Respond(c, http.StatusUnauthorized, "error", "用户未登录", nil)
		return
	}
	email, _ := c.Get("Email")
	// 解析请求体中的任务信息
	var TaskRegistration dto.TaskRegistrationRequest
	if err := c.ShouldBindJSON(&TaskRegistration); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}

	// 调用服务层的 JoinTask 方法处理加入任务逻辑
	err := services.JoinTask(TaskRegistration, nickname, email)
	if err != nil {
		if err.Error() == "活动不存在" {
			utils.Respond(c, http.StatusNotFound, "error", err.Error(), nil)
		} else if err.Error() == "活动已达到人数上限" {
			utils.Respond(c, http.StatusForbidden, "error", err.Error(), nil)
		} else {
			utils.Respond(c, http.StatusInternalServerError, "error", "加入活动失败："+err.Error(), nil)
		}
		return
	}

	utils.Respond(c, http.StatusOK, "success", "加入活动成功", nil)
}

// CreateTask 创建新任务
// @Summary 创建新任务
// @Description 创建新的志愿者活动
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Param name body dto.TaskInfo true "任务信息"
// @Router /admin/create_task [post]
func CreateTask(c *gin.Context) {
	var input dto.CreateTaskInfo
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}

	task, err := services.CreateTask(input.Name, input.StartTime, input.EndTime, input.Location, input.Limit)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "创建活动失败："+err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "活动创建成功", gin.H{"task": task})
}

// GetTasks 获取所有任务
// @Summary 获取所有任务
// @Description 获取所有活动列表
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Router /task/tasks [get]
func GetTasks(c *gin.Context) {
	tasks, err := services.GetTasks()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "获取活动列表失败"+err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "获取活动列表成功", gin.H{"tasks": tasks})
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 删除指定活动
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Param id path string true "活动ID"
// @Router /admin/delete_task/{id} [delete]
func DeleteTask(c *gin.Context) {
	taskId := c.Query("TaskId")
	err := services.DeleteTask(taskId)
	if err != nil {
		utils.Respond(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "活动删除成功", nil)
}

// UpdateTask 更新任务
// @Summary 更新任务
// @Description 更新指定活动的信息
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Param name body dto.TaskInfo true "任务信息"
// @Router /admin/update [post]
func UpdateTask(c *gin.Context) {
	var input dto.CreateTaskInfo
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}

	task, err := services.UpdateTask(input.Name, input.StartTime, input.StartTime, input.Location, input.Limit)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "修改活动失败："+err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "修改创建成功", gin.H{"task": task})
}

// GetTaskDetails 获取任务详情
// @Summary 获取任务详情
// @Description 获取任务详情
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Router /admin/getTaskDetails [get]
func GetTaskDetails(c *gin.Context) {
	taskDetails, err := services.GetTaskDetails()
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "获取详细活动列表失败"+err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "获取详细活动列表成功", gin.H{"taskDetails": taskDetails})
}

// GetTaskStatus 获取任务状态
// @Summary 获取任务状态
// @Description 获取当前用户任务状态
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.UserProfileRequest true "Bearer 用户令牌"
// @Param taskID path int true "任务ID"
// @Router /task/getStatus [get]
func GetTaskStatus(c *gin.Context) {
	// 从上下文中获取用户名
	nickname, exists := c.Get("Nickname")
	if !exists {
		// 如果用户名不存在，说明用户未登录，返回未授权错误
		utils.Respond(c, http.StatusUnauthorized, "error", "用户未登录", nil)
		return
	}
	taskIdStr := c.Query("TaskId")
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		// 如果转换失败，返回错误响应给客户端
		utils.Respond(c, http.StatusBadRequest, "error", "任务ID格式错误，必须为有效的整数", nil)
		return
	}
	// 调用服务层的 GetTaskStatus 方法获取任务状态
	taskStatus, err := services.GetTaskStatus(taskId, nickname.(string))
	if err != nil {
		// 如果获取任务状态失败，返回内部服务器错误
		utils.Respond(c, http.StatusInternalServerError, "error", "获取当前用户任务状态失败："+err.Error(), nil)
		return
	}

	// 如果获取任务状态成功，返回成功信息和任务状态
	utils.Respond(c, http.StatusOK, "success", "获取当前用户任务状态成功", gin.H{"taskStatus": taskStatus})
}

// GetTaskAuditDetail 获取任务报名详情
// @Summary 获取任务报名详情
// @Description 获取任务报名详情
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.AuditRequest true "Bearer 用户令牌"
// @Param taskId path int true "任务ID"
// @Router /admin/GetTaskAuditDetail [post]
func GetTaskAuditDetail(c *gin.Context) {
	var input dto.AuditRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}
	taskDetails, err := services.GetTaskAuditDetail(input.TaskId)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "获取活动报名情况列表失败"+err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "获取活动报名情况列表成功", gin.H{"taskDetails": taskDetails})
}

// ApproveVolunteer 通过报名人审核
// @Summary 通过报名人审核
// @Description 通过报名人审核
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.HandleVolunteerRequest true "Bearer 用户令牌"
// @Param taskId path int true "任务ID"
// @Router /admin/approveVolunteer [post]
func ApproveVolunteer(c *gin.Context) {
	var input dto.HandleVolunteerRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}
	err := services.ApproveVolunteer(input.TaskId, input.Email)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "通过报名人审核失败"+err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "通过报名人审核成功", nil)
}

// RejectVolunteer 拒绝报名人审核
// @Summary 拒绝报名人审核
// @Description 拒绝报名人审核
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header dto.HandleVolunteerRequest true "Bearer 用户令牌"
// @Param taskId path int true "任务ID"
// @Router /admin/rejectVolunteer [post]
func RejectVolunteer(c *gin.Context) {
	var input dto.HandleVolunteerRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Respond(c, http.StatusBadRequest, "error", "请求参数错误："+err.Error(), nil)
		return
	}
	err := services.RejectVolunteer(input.TaskId, input.Email)
	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, "error", "拒绝报名人审核失败"+err.Error(), nil)
		return
	}
	utils.Respond(c, http.StatusOK, "success", "拒绝报名人审核成功", nil)
}
