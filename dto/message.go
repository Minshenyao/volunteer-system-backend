package dto

// MessageResponse 消息响应结构体
type MessageResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
	Status  string `json:"status"`
}
