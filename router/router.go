package router

import (
	"github.com/Le-BlitzZz/real-time-chat-app/api"
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

	userHandler := api.UserAPI{DB: db.SQL}

	ui.Register(g)
	g.POST("/register", userHandler.Register)
	g.POST("/login", userHandler.Login)

	protected := g.Group("/")
	protected.Use(sessionManager.RequireSession())
	{
		// protected.GET("/user/me", userHandler.GetCurrentUser)
		protected.POST("/logout", userHandler.Logout)
	}

	return g
}

// func CreateOld(db *database.Database, conf *config.Configuration) *gin.Engine {
// 	g := gin.Default()

// 	g.Static("/static", "./frontend")
// 	g.LoadHTMLGlob("./frontend/*.html")

// 	g.GET("/", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "auth.html", nil)
// 	})

// 	g.POST("/register", func(c *gin.Context) {
// 		auth.Register(c, db)
// 	})
// 	g.POST("/login", func(c *gin.Context) {
// 		auth.Login(c, db)
// 	})

// 	g.GET(
// 		"/chat",
// 		func(c *gin.Context) {
// 			auth.Authenticate(c, db, conf.JWTSecret)
// 		},
// 		func(c *gin.Context) {
// 			stream.Chat(c, db)
// 		},
// 	)

// 	return g
// }
