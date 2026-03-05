package routes

import (
	"net/http"

	"LightChat/controllers"
	"LightChat/websocket"

	"github.com/gorilla/websocket"
)

// 升级HTTP连接为WebSocket连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 开发环境允许所有跨域，生产环境需限制
		return true
	},
}

// SetupRoutes 配置所有路由
func SetupRoutes() {
	// 健康检查路由
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 消息相关路由
	http.HandleFunc("/message/send", controllers.SendMessage)
	http.HandleFunc("/message/history", controllers.GetMessageHistory)

	// 群组相关路由
	http.HandleFunc("/group/create", controllers.CreateGroup)
	http.HandleFunc("/group/join", controllers.JoinGroup)
	http.HandleFunc("/group/members", controllers.GetGroupMembers)

	// 通话相关路由
	http.HandleFunc("/call/initiate", controllers.InitiateCall)
	http.HandleFunc("/call/end", controllers.EndCall)

	// WebSocket连接路由
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// 获取用户ID（实际项目需从登录态获取，这里简化）
		userId := r.URL.Query().Get("user_id")
		if userId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("user_id is required"))
			return
		}

		// 升级连接
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 创建客户端并加入hub
		client := websocket.NewClient(conn, userId)
		websocket.HubInstance.Register <- client

		// 启动客户端读写循环
		go client.ReadPump()
		go client.WritePump()
	})
}