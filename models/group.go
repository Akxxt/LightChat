package models

// Group 群组模型
type Group struct {
	ID        string   `json:"id"`         // 群组ID
	Name      string   `json:"name"`       // 群组名称
	CreatorID string   `json:"creator_id"` // 创建者ID
	MemberIDs []string `json:"member_ids"` // 成员ID列表
	CreatedAt int64    `json:"created_at"` // 创建时间戳
}