package utils

import "time"

// GetUserStatus 根据最近登录时间判断用户状态
func GetUserStatus(lastLoginTime time.Time) string {
	now := time.Now()
	diff := now.Sub(lastLoginTime)

	if diff <= 30*time.Minute {
		return "活跃"
	} else if diff <= time.Hour {
		return "空闲"
	} else {
		return "下线"
	}
}
