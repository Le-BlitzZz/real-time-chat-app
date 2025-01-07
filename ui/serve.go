package ui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Embed the public directory into the binary.
//
//go:embed public/*
var public embed.FS

// Register sets up the UI routes.
func Register(r *gin.Engine) {
	// Serve files from the embedded public directory
	publicFS, err := fs.Sub(public, "public")
	if err != nil {
		panic(err)
	}

	// Serve static assets
	r.StaticFS("/static", http.FS(publicFS))

	// Serve auth.html for the root route
	r.GET("/", func(ctx *gin.Context) {
		ctx.FileFromFS("auth.html", http.FS(publicFS))
	})

	r.GET("/register", func(ctx *gin.Context) {
		ctx.FileFromFS("register.html", http.FS(publicFS))
	})

	// Serve chat.html for the /chat route
	r.GET("/chat", func(ctx *gin.Context) {
		ctx.FileFromFS("chat.html", http.FS(publicFS))
	})
}
