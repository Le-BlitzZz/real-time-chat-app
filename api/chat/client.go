package chat

import (
	"fmt"
	"sync"
	"time"

	redisdb "github.com/Le-BlitzZz/real-time-chat-app/database/redis"
	"github.com/Le-BlitzZz/real-time-chat-app/model/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 2 * time.Second
)

type client struct {
	conn    *websocket.Conn
	onClose func(*client)
	write   chan redis.Message
	chatID  uint
	userID  uint
	once    sync.Once
}

// newClient creates a new WebSocket client.
func newClient(conn *websocket.Conn, userID, chatID uint, onClose func(*client)) *client {
	return &client{
		conn:    conn,
		write:   make(chan redis.Message, 10),
		userID:  userID,
		chatID:  chatID,
		onClose: onClose,
	}
}

// Close closes the WebSocket connection and its resources.
func (c *client) Close() {
	c.once.Do(func() {
		c.conn.Close()
		close(c.write)
	})
}

// NotifyClose closes the WebSocket connection and notifies that the connection was closed.
func (c *client) notifyClose() {
	c.once.Do(func() {
		c.conn.Close()
		close(c.write)
		c.onClose(c)
	})
}

// startReading reads messages from the WebSocket connection.
func (c *client) startReading(ctx *gin.Context, redisClient *redisdb.RedisDb, chatID uint, broadcastFunc func(uint, redis.Message)) {
	defer c.notifyClose()

	username, exists := ctx.Get("userName")
	if !exists {
		fmt.Println("Username not found in context")
		return
	}

	for {
		_, rawMessage, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		messageIncoming, err := redis.DeserializeMessageIncoming(string(rawMessage))
		if err != nil {
			fmt.Println("Error parsing message:", err)
			continue
		}

		message := redis.NewMessage(username.(string), messageIncoming.Message, chatID)

		redisClient.SaveMessage(ctx, message)
		redisClient.PublishMessage(ctx, message)

		broadcastFunc(chatID, message)
	}
}

func (c *client) startWriting() {
	defer c.notifyClose()

	for msg := range c.write {
		c.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := c.conn.WriteJSON(msg); err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}
