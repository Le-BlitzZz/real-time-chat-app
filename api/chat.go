package api

import (
	"context"
	"net/http"

	"github.com/Le-BlitzZz/real-time-chat-app/database/redis"
	"github.com/Le-BlitzZz/real-time-chat-app/model/sql"
	"github.com/gin-gonic/gin"
)

type ChatSQLDatabase interface {
	CreateChat(chat *sql.Chat) error
	GetChatByID(chatID uint) (*sql.Chat, error)
}

type ChatRedisDatabase interface {
	AddUserToChat(ctx context.Context, userID uint, chatID uint) error
	GetUserChats(ctx context.Context, userID uint) ([]string, error)
	RemoveUserFromChat(ctx context.Context, userID uint, chatID uint) error
}

type ChatAPI struct {
	SQLDB   ChatSQLDatabase
	RedisDB ChatRedisDatabase
}

func (api *ChatAPI) CreateChat(ctx *gin.Context) {
	var chatRequest struct {
		Name string `json:"name,omitempty"`          // Optional for group chats
		Type string `json:"type" binding:"required"` // Must be "private" or "group"
	}
	if err := ctx.ShouldBindJSON(&chatRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if chatRequest.Type != sql.TypePrivate && chatRequest.Type != sql.TypeGroup {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat type"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var chat *sql.Chat
	if chatRequest.Type == sql.TypePrivate {
		chat = sql.PrivateChat
	} else {
		chat = sql.NewGroupChat(chatRequest.Name)
	}

	if err := api.SQLDB.CreateChat(chat); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create chat"})
		return
	}

	// Add the creator to the chat in Redis
	if err := api.RedisDB.AddUserToChat(ctx, userID.(uint), chat.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not add user to chat"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "chat created", "chat": chat})
}

func (api *ChatAPI) ListChats(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	chatIDsStr, err := api.RedisDB.GetUserChats(ctx, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve chats"})
		return
	}

	var chats []sql.Chat
	for _, chatIDStr := range chatIDsStr {
		chatID := redis.ParseUint(chatIDStr)
		chat, err := api.SQLDB.GetChatByID(chatID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve chat metadata"})
			return
		}
		chats = append(chats, *chat)
	}

	ctx.JSON(http.StatusOK, chats)
}

func (api *ChatAPI) LeaveChat(ctx *gin.Context) {
	var leaveRequest struct {
		ChatID uint `json:"chat_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&leaveRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if leaveRequest.ChatID == sql.GlobalChatID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot leave the global chat"})
		return
	}

	if err := api.RedisDB.RemoveUserFromChat(ctx, userID.(uint), leaveRequest.ChatID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not leave chat"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "left chat successfully"})
}
