package services

import (
	"volunteer-system-backend/dto"
	"volunteer-system-backend/models"
	"volunteer-system-backend/utils"
	"errors"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

func CreateTask(name, startTime, endTime, location string, limit uint) (dto.TaskInfo, error) {
	// 检查活动名是否已存在
	var existingName models.Task
	if err := models.DB.Where("name = ?", strings.TrimSpace(name)).First(&existingName).Error; err == nil {
		return dto.TaskInfo{}, errors.New("该活动名称已经存在")
	}

	// 创建新活动
	task := models.Task{
		Name:      name,
		CreatedAt: utils.FormatStr2Time(time.Now().Local().Format("2006-01-02 15:04:05")),
		StartTime: utils.FormatStr2Time(startTime),
		EndTime:   utils.FormatStr2Time(endTime),
		Location:  location,
		Limit:     limit,
		Joined:    0,
		// 确保 Participants 字段为空
		Participants: []models.TaskParticipant{},
	}
	if err := models.DB.Create(&task).Error; err != nil {
		return dto.TaskInfo{}, errors.New("无法创建活动")
	}
	return utils.ConvertTaskToDTO(task), nil
}

func GetTasks() ([]dto.TaskInfo, error) {
	var tasks []models.Task
	if err := models.DB.Preload("Participants").Find(&tasks).Error; err != nil {
		return nil, err
	}
	newTask := make([]dto.TaskInfo, len(tasks))
	for i, task := range tasks {
		newTask[i] = utils.ConvertTaskToDTO(task)
	}
	return newTask, nil
}

func JoinTask(taskInfo dto.TaskRegistrationRequest, nickname, email any) error {
	// 检查活动是否存在
	var task models.Task
	if err := models.DB.Where("id = ? AND name = ? ", taskInfo.ID, taskInfo.Name).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("活动不存在")
		}
		return err
	}
	// 检查活动是否已达到人数上限
	if task.Joined >= task.Limit && task.Limit != 0 {
		return errors.New("活动已达到人数上限")
	}

	// 检查用户是否已经报名
	var participant models.TaskParticipant
	if err := models.DB.Where("task_id = ? AND email = ?", task.ID, email).First(&participant).Error; err == nil {
		return errors.New("用户已经报名该活动")
	}

	// 创建新的参加人员记录
	newParticipant := models.TaskParticipant{
		TaskID:   task.ID,
		Nickname: nickname.(string),
		Email:    email.(string),
	}
	if err := models.DB.Create(&newParticipant).Error; err != nil {
		return errors.New("无法记录用户报名信息")
	}

	// 更新已参加人数
	task.Joined++
	if err := models.DB.Save(&task).Error; err != nil {
		return errors.New("无法更新活动参加人数")
	}
	var TaskParticipant dto.TaskParticipant
	if err := models.DB.Model(&TaskParticipant).Where("task_id = ? AND email = ?", task.ID, email).Update("status", "0").Error; err != nil {
		return errors.New("无法更新任务状态")
	}
	_, err := CreateMessage(1, "新的待审核通知", "管理员您好，您有一条新的待审核活动")
	if err != nil {
		return errors.New("创建消息失败")
	}
	return nil
}

func DeleteTask(id string) error {
	result := models.DB.Delete(&models.Task{}, id)
	if result.RowsAffected == 0 {
		return errors.New("活动不存在")
	}
	return result.Error
}

func UpdateTask(name, startTime, endTime, location string, limit uint) (dto.TaskInfo, error) {
	var task models.Task
	if err := models.DB.Where("name = ?", strings.TrimSpace(name)).First(&task).Error; err != nil {
		return dto.TaskInfo{}, errors.New("该活动不存在，无法修改")
	}

	// 使用 map 更新特定字段
	updates := map[string]interface{}{
		"start_time": utils.FormatStr2Time(startTime),
		"end_time":   utils.FormatStr2Time(endTime),
		"location":   location,
		"limit":      limit,
	}

	if err := models.DB.Model(&task).Updates(updates).Error; err != nil {
		return dto.TaskInfo{}, errors.New("无法修改活动")
	}

	// 重新查询更新后的完整记录
	if err := models.DB.First(&task, task.ID).Error; err != nil {
		return dto.TaskInfo{}, errors.New("无法获取更新后的活动信息")
	}

	return utils.ConvertTaskToDTO(task), nil
}

// GetTaskDetails 获取任务详情
func GetTaskDetails() ([]dto.TaskDetailResponse, error) {
	var tasks []models.Task
	if err := models.DB.Preload("Participants").Find(&tasks).Error; err != nil {
		return nil, err
	}

	var taskDetails []dto.TaskDetailResponse
	for _, task := range tasks {
		var participantUsernames []models.Participant
		for _, participant := range task.Participants {
			currUser, err := GetUserProfile(participant.Email)
			// 调用服务层的 GetTaskStatus 方法获取任务状态
			taskStatus, err := GetTaskStatus(int(participant.TaskID), currUser.Nickname)
			if err != nil {
				return nil, err
			}
			if taskStatus.Status != 1 {
				continue
			}
			userInfo := models.Participant{
				Email:         currUser.Email,
				Nickname:      currUser.Nickname,
				Gender:        currUser.Gender,
				Avatar:        currUser.Avatar,
				Phone:         currUser.Phone,
				Duration:      currUser.Duration,
				LastLoginTime: utils.FormatTime2Str(currUser.LastLoginTime),
			}
			participantUsernames = append(participantUsernames, userInfo)
		}

		taskDetail := dto.TaskDetailResponse{
			ID:           task.ID,
			Name:         task.Name,
			CreatedAt:    utils.FormatTime2Str(task.CreatedAt),
			StartTime:    utils.FormatTime2Str(task.StartTime),
			EndTime:      utils.FormatTime2Str(task.EndTime),
			Location:     task.Location,
			Limit:        task.Limit,
			Joined:       task.Joined,
			Participants: participantUsernames,
		}
		taskDetails = append(taskDetails, taskDetail)
	}

	return taskDetails, nil
}

// GetTaskStatus 根据任务 ID 和用户名获取任务状态信息
func GetTaskStatus(taskID int, nickname string) (dto.ParticipantStatus, error) {
	// 首先根据任务 ID 查找任务
	var task models.Task
	if err := models.DB.Preload("Participants").Where("id = ?", taskID).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ParticipantStatus{}, errors.New("任务不存在")
		}
		return dto.ParticipantStatus{}, errors.New("任务不存在")
	}

	// 查找该用户在任务中的参与记录
	var participant models.TaskParticipant
	if err := models.DB.Where("task_id = ? AND nickname = ?", taskID, nickname).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ParticipantStatus{}, errors.New("用户未参与该任务")
		}
		return dto.ParticipantStatus{}, errors.New("用户未参与该任务")
	}
	taskParticipantStatus := dto.ParticipantStatus{
		ID:       participant.TaskID,
		Nickname: participant.Nickname,
		Status:   participant.Status,
	}
	return taskParticipantStatus, nil
}

func GetTaskAuditDetail(TaskId uint) (dto.AuditResponse, error) {
	var TaskParticipant models.TaskParticipant
	if err := models.DB.Where("task_id = ?", TaskId).First(&TaskParticipant).Error; err != nil {
		return dto.AuditResponse{}, errors.New("该活动暂时没有人报名或已结束")
	}
	var volunteers []dto.JoinTaskVolunteer
	if err := models.DB.Table("task_participants").Where("task_id = ? AND status = 0", TaskId).Find(&volunteers).Error; err != nil {
		return dto.AuditResponse{}, err
	}
	AuditResponse := dto.AuditResponse{
		TaskId:     TaskId,
		Volunteers: volunteers,
	}
	return AuditResponse, nil
}

// GetTaskNameByID 根据任务ID查询活动名称
func GetTaskNameByID(taskId uint) (string, error) {
	var task models.Task
	if err := models.DB.Where("id = ?", taskId).First(&task).Error; err != nil {
		// 如果查询出错，判断是否是记录未找到的错误
		return "", errors.New("活动不存在")
	}
	return task.Name, nil
}

// CreateMessage 创建新消息
func CreateMessage(userID uint, title, content string) (*models.Message, error) {
	message := &models.Message{
		UserID:  userID,
		Title:   title,
		Content: content,
		Time:    time.Now(),
		Status:  "unread",
	}
	err := models.DB.Create(message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}

// ApproveVolunteer 通过报名人审核
func ApproveVolunteer(taskId uint, email string) error {
	var TaskParticipant models.TaskParticipant
	if err := models.DB.Where("task_id = ? AND email = ?", taskId, email).First(&TaskParticipant).Error; err != nil {
		return errors.New("该活动已结束或该用户没有报名该活动")
	}
	if TaskParticipant.Status != 0 {
		if TaskParticipant.Status == 3 {
			return errors.New("该用户未参加该活动")
		}
		return errors.New("该用户无需审核")
	}
	if TaskParticipant.Status == 0 {
		if err := models.DB.Model(&TaskParticipant).Update("status", 1).Error; err != nil {
			return errors.New("更新任务状态失败")
		}
		// 创建消息
		userId, err := GetUserIDByEmail(email)
		if err != nil {
			return errors.New("获取用户ID失败")
		}
		taskName, err := GetTaskNameByID(taskId)
		if err != nil {
			return errors.New("获取任务名称失败")
		}
		_, err = CreateMessage(userId, "申请通过通知", "您的申请已被通过，活动名称: \""+taskName+"\"")
		if err != nil {
			log.Println(err)
			return errors.New("创建消息失败")
		}
		return nil
	}
	return nil
}

// RejectVolunteer 拒绝报名人审核
func RejectVolunteer(taskId uint, email string) error {
	var TaskParticipant models.TaskParticipant
	if err := models.DB.Where("task_id = ? AND email = ?", taskId, email).First(&TaskParticipant).Error; err != nil {
		return errors.New("该活动已结束或该用户没有报名该活动")
	}
	if TaskParticipant.Status != 0 {
		if TaskParticipant.Status == 3 {
			return errors.New("该用户未参加该活动")
		}
		return errors.New("该用户无需审核")
	}
	if TaskParticipant.Status == 0 {
		if err := models.DB.Model(&TaskParticipant).Update("status", 2).Error; err != nil {
			return errors.New("更新任务状态失败")
		}
		// 创建消息
		userId, err := GetUserIDByEmail(email)
		if err != nil {
			return errors.New("获取用户ID失败")
		}
		taskName, err := GetTaskNameByID(taskId)
		if err != nil {
			return errors.New("获取任务名称失败")
		}
		_, err = CreateMessage(userId, "申请拒绝通知", "很抱歉，您的申请已被拒绝，活动名称: \""+taskName+"\"")
		if err != nil {
			log.Println(err)
			return errors.New("创建消息失败")
		}
		return nil
	}
	return nil
}
