package websocket

import (
	"sync"
)

// Message WebSocket消息结构
type Message struct {
	UserID string // 目标用户ID
	Data   []byte // 消息数据
}

// Hub WebSocket集线器，管理所有客户端连接
type Hub struct {
	// 注册的客户端
	Clients map[string]*Client

	// 注册请求通道
	Register chan *Client

	// 注销请求通道
	Unregister chan *Client

	// 广播消息通道
	Broadcast chan Message

	// 互斥锁，保护Clients
	mu sync.RWMutex
}

// HubInstance 全局Hub实例
var HubInstance = NewHub()

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

// Run 启动Hub的主循环
func (h *Hub) Run() {
	for {
		select {
		// 处理客户端注册
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()

		// 处理客户端注销
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()

		// 处理广播消息
		case msg := <-h.Broadcast:
			h.mu.RLock()
			client, ok := h.Clients[msg.UserID]
			h.mu.RUnlock()

			if ok {
				select {
				case client.Send <- msg.Data:
				default:
					close(client.Send)
					h.mu.Lock()
					delete(h.Clients, client.UserID)
					h.mu.Unlock()
				}
			}
		}
	}
}