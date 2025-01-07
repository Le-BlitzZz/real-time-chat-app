package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type SessionManager struct {
	store *sessions.CookieStore
}

func NewSessionManager() *SessionManager {
	secret := generateSecret()
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production for HTTPS
	}
	return &SessionManager{store: store}
}

func (sm *SessionManager) SetSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("sessionManager", sm)
		ctx.Next()
	}
}

func (sm *SessionManager) RequireSession() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, _ := sm.store.Get(ctx.Request, "chatapp-session")
		userID, ok := session.Values["user_id"].(uint)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Set("userID", userID)
		ctx.Set("userName", session.Values["user_name"].(string))
		ctx.Set("admin", session.Values["admin"].(bool))
		ctx.Next()
	}
}

func (sm *SessionManager) CreateSession(w http.ResponseWriter, r *http.Request, userID uint, userName string, isAdmin bool) error {
	session, _ := sm.store.Get(r, "chatapp-session")
	session.Values["user_id"] = userID
	session.Values["user_name"] = userName
	session.Values["admin"] = isAdmin
	return session.Save(r, w)
}

func (sm *SessionManager) DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, _ := sm.store.Get(r, "chatapp-session")
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func generateSecret() string {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		log.Fatalf("Failed to generate session secret: %v", err)
	}
	return base64.URLEncoding.EncodeToString(secret)
}
