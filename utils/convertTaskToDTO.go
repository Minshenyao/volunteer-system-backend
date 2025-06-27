package utils

import (
	"volunteer-system-backend/dto"
	"volunteer-system-backend/models"
)

// ConvertTaskToDTO 转换函数
func ConvertTaskToDTO(task models.Task) dto.TaskInfo {
	return dto.TaskInfo{
		ID:        task.ID,
		Name:      task.Name,
		CreatedAt: FormatTime2Str(task.CreatedAt),
		StartTime: FormatTime2Str(task.StartTime),
		EndTime:   FormatTime2Str(task.EndTime),
		Location:  task.Location,
		Limit:     task.Limit,
		Joined:    task.Joined,
	}
}
