package models

// Message 消息模型
type Message struct {
	ID         string `json:"id"`          // 消息ID
	FromUserID string `json:"from_user_id"`// 发送者ID
	ToUserID   string `json:"to_user_id"`  // 接收者ID
	Content    string `json:"content"`     // 消息内容
	MsgType    string `json:"msg_type"`    // 消息类型：text/image/voice
	CreatedAt  int64  `json:"created_at"`  // 创建时间戳
}