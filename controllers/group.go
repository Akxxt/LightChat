package controllers

import (
	"encoding/json"
	"net/http"
	"LightChat/models"
	"LightChat/utils"
	"time"
)

// CreateGroup 创建群组
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CreatorID string   `json:"creator_id"`
		GroupName string   `json:"group_name"`
		MemberIDs []string `json:"member_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求参数: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 验证参数
	if req.CreatorID == "" || req.GroupName == "" || len(req.MemberIDs) == 0 {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	// 构建群组模型
	group := models.Group{
		ID:        utils.GenerateID(),
		Name:      req.GroupName,
		CreatorID: req.CreatorID,
		MemberIDs: append(req.MemberIDs, req.CreatorID), // 把创建者加入成员
		CreatedAt: time.Now().Unix(),
	}

	// 保存群组信息到Redis
	redisClient := utils.GetRedisClient()
	groupJSON, _ := json.Marshal(group)
	_ = redisClient.Set(r.Context(), "group:"+group.ID, groupJSON, 0).Err()

	// 返回响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":  200,
		"msg":   "群组创建成功",
		"group": group,
	})
}

// JoinGroup 加入群组
func JoinGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GroupID string `json:"group_id"`
		UserID  string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求参数: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 获取群组信息
	redisClient := utils.GetRedisClient()
	groupStr, err := redisClient.Get(r.Context(), "group:"+req.GroupID).Result()
	if err != nil {
		http.Error(w, "群组不存在", http.StatusNotFound)
		return
	}

	var group models.Group
	_ = json.Unmarshal([]byte(groupStr), &group)

	// 检查是否已加入
	for _, mid := range group.MemberIDs {
		if mid == req.UserID {
			http.Error(w, "已加入该群组", http.StatusBadRequest)
			return
		}
	}

	// 添加成员
	group.MemberIDs = append(group.MemberIDs, req.UserID)
	groupJSON, _ := json.Marshal(group)
	_ = redisClient.Set(r.Context(), "group:"+req.GroupID, groupJSON, 0).Err()

	// 返回响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "加入群组成功",
	})
}

// GetGroupMembers 获取群组成员
func GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		http.Error(w, "group_id is required", http.StatusBadRequest)
		return
	}

	// 获取群组信息
	redisClient := utils.GetRedisClient()
	groupStr, err := redisClient.Get(r.Context(), "group:"+groupID).Result()
	if err != nil {
		http.Error(w, "群组不存在", http.StatusNotFound)
		return
	}

	var group models.Group
	_ = json.Unmarshal([]byte(groupStr), &group)

	// 返回响应
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":     200,
		"msg":      "success",
		"group_id": group.ID,
		"members":  group.MemberIDs,
		"count":    len(group.MemberIDs),
	})
}