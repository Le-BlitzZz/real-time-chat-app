package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Le-BlitzZz/real-time-chat-app/database/redis"
	"github.com/gin-gonic/gin"
)

type ChatDatabase interface {
	CreateChat(ctx context.Context, chatID uint, meta map[string]interface{}) error
	AddUserToChat(ctx context.Context, userID uint, chatID uint) error
	RemoveUserFromChat(ctx context.Context, userID uint, chatID uint) error
	ChatExists(ctx context.Context, chatID uint) (bool, error)
	GetUserChats(ctx context.Context, userID uint) ([]uint, error)
	RemoveChat(ctx context.Context, chatID uint) error
	GetOrCreatePrivateChatID(ctx context.Context, userID1, userID2 uint) (uint, error)
	GetChatUsers(ctx context.Context, chatID uint) ([]uint, error)
}

type ChatAPI struct {
	DB     ChatDatabase
	UserDB UserDatabase
}

type chatRequest struct {
	ReceiverNickname string `json:"receiver_nickname,omitempty"`
	ChatID           uint   `json:"chat_id"`
}

func (api *ChatAPI) Start(ctx *gin.Context) {
	fmt.Println("IN START")
	var startChat = chatRequest{}

	if err := ctx.ShouldBindJSON(&startChat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	fmt.Println("startChat.ReceiverNickname", startChat.ReceiverNickname)

	if startChat.ReceiverNickname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "receiver nickname is required for private chat"})
		return
	}

	senderID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	receiver, err := api.UserDB.GetUserByName(startChat.ReceiverNickname)
	if err != nil || receiver == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "receiver not found"})
		return
	}

	chatID, err := api.DB.GetOrCreatePrivateChatID(ctx, senderID.(uint), receiver.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create or retrieve private chat"})
		return
	}

	if err := api.DB.AddUserToChat(ctx, senderID.(uint), chatID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not add sender to chat"})
		return
	}

	if err := api.DB.AddUserToChat(ctx, receiver.ID, chatID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not add receiver to chat"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "chat started", "chat_id": chatID})
}

func (api *ChatAPI) ListChats(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	chatIDs, err := api.DB.GetUserChats(ctx, userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve chats"})
		return
	}

	fmt.Println("chatIds", chatIDs)
	var chats []map[string]interface{}
	for _, chatID := range chatIDs {
		if chatID == redis.GlobalChatID {
			continue
		}

		userIDs, err := api.DB.GetChatUsers(ctx, chatID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve chat participants"})
			return
		}

		var otherUserID uint
		for _, id := range userIDs {
			if id != userID.(uint) {
				otherUserID = id
				break
			}
		}

		otherUser, err := api.UserDB.GetUserByID(otherUserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve user information"})
			return
		}

		chats = append(chats, map[string]interface{}{
			"chat_id":   chatID,
			"other_user_name": otherUser.Name,
		})
	}

	for _, ccc := range chats {
		fmt.Println(ccc)
	}
	ctx.JSON(http.StatusOK, chats)
}
