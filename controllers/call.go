package controllers

import (
	"encoding/json"
	"net/http"
	"LightChat/websocket"
)

// InitiateCall 发起通话
func InitiateCall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromUserID string `json:"from_user_id"`
		ToUserID   string `json:"to_user_id"`
		CallType   string `json:"call_type"` // voice, video
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求参数: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.FromUserID == "" || req.ToUserID == "" || req.CallType == "" {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 推送通话请求
	callReq := map[string]interface{}{
		"type":      "call_request",
		"from_user": req.FromUserID,
		"call_type": req.CallType,
		"timestamp": time.Now().Unix(),
	}
	callJSON, _ := json.Marshal(callReq)
	websocket.HubInstance.Broadcast <- websocket.Message{
		UserID: req.ToUserID,
		Data:   callJSON,
	}

	// 返回响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "通话请求已发送",
	})
}

// EndCall 结束通话
func EndCall(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromUserID string `json:"from_user_id"`
		ToUserID   string `json:"to_user_id"`
		Reason     string `json:"reason"` // normal, cancel, timeout
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求参数: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 推送结束通话通知
	callEnd := map[string]interface{}{
		"type":      "call_end",
		"from_user": req.FromUserID,
		"reason":    req.Reason,
		"timestamp": time.Now().Unix(),
	}
	callJSON, _ := json.Marshal(callEnd)
	websocket.HubInstance.Broadcast <- websocket.Message{
		UserID: req.ToUserID,
		Data:   callJSON,
	}

	// 返回响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "通话已结束",
	})
}