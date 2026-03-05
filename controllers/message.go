package controllers

import (
	"encoding/json"
	"net/http"
	"LightChat/models"
	"LightChat/utils"
	"LightChat/websocket"
	"time"

	"github.com/redis/go-redis/v9"
)

// SendMessage 发送消息接口
func SendMessage(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	var req struct {
		FromUserID string `json:"from_user_id"`
		ToUserID   string `json:"to_user_id"`
		Content    string `json:"content"`
		MsgType    string `json:"msg_type"` // text, image, voice
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求参数: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证必填参数
	if req.FromUserID == "" || req.ToUserID == "" || req.Content == "" {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 构建消息模型
	msg := models.Message{
		ID:         utils.GenerateID(),
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Content:    req.Content,
		MsgType:    req.MsgType,
		CreatedAt:  time.Now().Unix(),
	}

	// 1. 保存消息到Redis
	redisClient := utils.GetRedisClient()
	msgJSON, _ := json.Marshal(msg)
	_ = redisClient.LPush(r.Context(), "msg:"+req.ToUserID, msgJSON).Err()

	// 2. 通过WebSocket推送消息
	websocket.HubInstance.Broadcast <- websocket.Message{
		UserID: req.ToUserID,
		Data:   msgJSON,
	}

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"msg":     "消息发送成功",
		"message": msg,
	})
}

// GetMessageHistory 获取消息历史
func GetMessageHistory(w http.ResponseWriter, r *http.Request) {
	// 获取用户ID
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	// 从Redis获取消息历史
	redisClient := utils.GetRedisClient()
	msgs, err := redisClient.LRange(r.Context(), "msg:"+userID, 0, -1).Result()
	if err != nil && err != redis.Nil {
		http.Error(w, "获取消息失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 解析消息列表
	var messageList []models.Message
	for _, msgStr := range msgs {
		var msg models.Message
		_ = json.Unmarshal([]byte(msgStr), &msg)
		messageList = append(messageList, msg)
	}

	// 返回响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  200,
		"msg":   "success",
		"data":  messageList,
		"count": len(messageList),
	})
}