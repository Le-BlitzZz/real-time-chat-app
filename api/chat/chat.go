package chat

import (
	"net/http"
	"sync"

	"github.com/Le-BlitzZz/real-time-chat-app/database/redis"
	redismodel "github.com/Le-BlitzZz/real-time-chat-app/model/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type API struct {
	clients     map[uint][]*client // Map chatID -> clients
	lock        sync.RWMutex
	upgrader    *websocket.Upgrader
	redisClient *redis.RedisDb
}

// New creates a new WebSocket stream API.
func New(redisClient *redis.RedisDb) *API {
	return &API{
		clients:     make(map[uint][]*client),
		upgrader:    &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		redisClient: redisClient,
	}
}

// Initialize upgrades HTTP to WebSocket and registers a client.
func (api *API) Initialize(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	chatIDstr := ctx.Query("chatID")
	if chatIDstr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "chatID is required"})
		return
	}
	chatID := redis.ParseUint(chatIDstr)

	conn, err := api.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
		return
	}

	client := newClient(conn, userID.(uint), chatID, api.removeClient)
	api.addClientToChat(client, chatID)

	go func() {
		messages, _ := api.redisClient.GetRecentMessages(ctx, chatID)
		for _, msg := range messages {
			client.write <- msg
		}
	}()

	go client.startReading(ctx, api.redisClient, chatID, api.broadcast)
	go client.startWriting()
}

func (api *API) broadcast(chatID uint, message redismodel.Message) {
	api.lock.RLock()
	defer api.lock.RUnlock()

	for _, client := range api.clients[chatID] {
		select {
		case client.write <- message:
		default:
		}
	}
}

// addClientToChat adds a client to the chat.
func (api *API) addClientToChat(client *client, chatID uint) {
	api.lock.Lock()
	defer api.lock.Unlock()
	api.clients[chatID] = append(api.clients[chatID], client)
}

// removeClient removes a client from chats.
func (api *API) removeClient(client *client) {
	api.lock.Lock()
	defer api.lock.Unlock()
	clients := api.clients[client.chatID]
	for i, c := range clients {
		if c == client {
			api.clients[client.chatID] = append(clients[:i], clients[i+1:]...)
			return
		}
	}
}
