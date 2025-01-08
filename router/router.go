package router

import (
	"net/http"

	"github.com/Le-BlitzZz/real-time-chat-app/api"
	"github.com/Le-BlitzZz/real-time-chat-app/api/channel"
	"github.com/Le-BlitzZz/real-time-chat-app/auth"
	"github.com/Le-BlitzZz/real-time-chat-app/config"
	"github.com/Le-BlitzZz/real-time-chat-app/database"
	"github.com/Le-BlitzZz/real-time-chat-app/ui"
	"github.com/gin-gonic/gin"
)

func Create(db *database.Database, conf *config.Configuration) *gin.Engine {
	g := gin.Default()

	g.RemoteIPHeaders = []string{"X-Forwarded-For"}
	g.ForwardedByClientIP = true

	g.Use(func(ctx *gin.Context) {
		// Map sockets "@" to 127.0.0.1, because gin-gonic can only trust IPs.
		if ctx.Request.RemoteAddr == "@" {
			ctx.Request.RemoteAddr = "127.0.0.1:65535"
		}
	})

	sessionManager := auth.NewSessionManager()
	g.Use(sessionManager.SetSession())

	userAPI := api.UserAPI{DB: db.SQL}
	chatAPI := api.ChatAPI{SQLDB: db.SQL, RedisDB: db.Redis}
	
	channelAPI := channel.New(db.Redis)

	ui.Register(g)
	g.POST("/register", userAPI.Register)
	g.POST("/login", userAPI.Login)
	g.POST("/logout", sessionManager.RequireSession(), userAPI.Logout)

	g.POST("/friend-request", sessionManager.RequireSession(), userAPI.SendFriendRequest)
	g.POST("/friend-request/accept", sessionManager.RequireSession(), userAPI.AcceptFriendRequest)
	g.GET("/friends", sessionManager.RequireSession(), userAPI.GetFriends)
	g.GET("/friend-requests", sessionManager.RequireSession(), userAPI.GetFriendRequests)

	g.GET("/ws/chat", sessionManager.RequireSession(), func(c *gin.Context) {
		chatID := c.Query("chatID")
		if chatID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "chatID is required"})
			return
		}
		c.Set("chatID", chatID)
		channelAPI.Initialize(c)
	})

	g.POST("/chat", sessionManager.RequireSession(), chatAPI.CreateChat)
	g.GET("/chats", sessionManager.RequireSession(), chatAPI.ListChats)
	g.POST("/chat/leave", sessionManager.RequireSession(), chatAPI.LeaveChat)

	return g
}
