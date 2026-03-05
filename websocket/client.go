package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端
type Client struct {
	Conn   *websocket.Conn // WebSocket连接
	UserID string          // 用户ID
	Send   chan []byte     // 待发送消息通道
}

// 心跳超时时间
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024000 // 1MB
)

// NewClient 创建新的客户端
func NewClient(conn *websocket.Conn, userID string) *Client {
	return &Client{
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []byte, 256),
	}
}

// ReadPump 读取客户端消息
func (c *Client) ReadPump() {
	defer func() {
		HubInstance.Unregister <- c
		c.Conn.Close()
	}()

	// 设置读取参数
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// 循环读取消息
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// 实际项目中可处理客户端发来的消息，这里简化
	}
}

// WritePump 发送消息给客户端
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		// 发送消息
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 发送缓冲区所有消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		// 发送心跳
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}