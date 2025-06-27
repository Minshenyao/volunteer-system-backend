package dto

import "volunteer-system-backend/models"

// TaskInfo 任务信息
type TaskInfo struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	CreatedAt string `json:"created_at" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Location  string `json:"location" binding:"required"`
	Limit     uint   `json:"limit" binding:"required"`
	Joined    uint   `json:"joined" binding:"required"`
}

// CreateTaskInfo 任务信息
type CreateTaskInfo struct {
	Name      string `json:"name" binding:"required"`
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binging:"required"`
	Location  string `json:"location" binging:"required"`
	Limit     uint   `json:"limit" binging:"required"`
}

// UpdateTaskStatusRequest 任务状态更新请求
type UpdateTaskStatusRequest struct {
	TaskID    string `json:"taskId" binding:"required"`
	NewStatus string `json:"newStatus" binding:"required"`
}

// TaskRegistrationRequest 任务报名请求
type TaskRegistrationRequest struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time" binging:"required"`
	Location  string `json:"location" binging:"required"`
}

// TaskApprovalRequest 任务审核请求
type TaskApprovalRequest struct {
	RegistrationID string `json:"registrationId" binding:"required"`
	Approved       bool   `json:"approved"`
}

// TaskDetailResponse 任务详情响应结构体
type TaskDetailResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
	Limit     uint   `json:"limit"`
	Joined    uint   `json:"joined"`
	// 参加任务的用户用户名列表
	Participants []models.Participant `json:"participants"`
}

// TaskParticipant 表示任务的已参加人员信息
type TaskParticipant struct {
	TaskID   uint   `json:"taskID" binding:"required"`   // 任务ID
	Nickname string `json:"nickname" binding:"required"` // 参加人员的昵称
	Status   uint   `json:"status" binding:"required"`   // 0表示待审核 1表示审核通过 2表示审核不通过 3表示未参加
}

type ParticipantStatusRequest struct {
	TaskId uint `json:"TaskId" binding:"required"`
}
type ParticipantStatus struct {
	ID       uint   `json:"id" binding:"required"`
	Nickname string `json:"nickname" binding:"required"` // 参加人员的昵称
	Status   uint   `json:"status" binding:"required"`   // 0表示待审核 1表示审核通过 2表示审核不通过 3表示未参加
}

type AuditRequest struct {
	TaskId uint `json:"taskId" binding:"required"`
}

type JoinTaskVolunteer struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Status   uint   `json:"status" binding:"required"`
}

type AuditResponse struct {
	TaskId     uint                `json:"taskId" binding:"required"`
	Volunteers []JoinTaskVolunteer `json:"volunteers"`
}

type HandleVolunteerRequest struct {
	TaskId uint   `json:"taskId" binding:"required"`
	Email  string `json:"email" binding:"required"`
}
