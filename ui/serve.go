package ui

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.Static("/static", "./ui/public")
	r.LoadHTMLGlob("./ui/public/*.html")

	// Serve HTML pages directly for key routes
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "auth.html", nil)
	})

	r.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", nil)
	})

	r.GET("/chat", func(ctx *gin.Context) {
		chatID := ctx.DefaultQuery("chatID", "0")
		ctx.HTML(http.StatusOK, "chat.html", gin.H{"chatID": chatID})
	})

	r.GET("/chat/:chatID", func(ctx *gin.Context) {
		chatID := ctx.Param("chatID")
		ctx.HTML(http.StatusOK, "private_chat.html", gin.H{"chatID": chatID})
	})
}
