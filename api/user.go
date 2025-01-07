package api

import (
	"net/http"

	"github.com/Le-BlitzZz/real-time-chat-app/auth"
	"github.com/Le-BlitzZz/real-time-chat-app/auth/password"
	"github.com/Le-BlitzZz/real-time-chat-app/model/sql"
	"github.com/gin-gonic/gin"
)

type UserDatabase interface {
	GetUserByName(name string) (*sql.User, error)
	GetUserByEmail(email string) (*sql.User, error)
	CreateUser(user *sql.User) error
	// UpdateUser(user *sql.User) error
}

type LoginUser struct {
	Identifier string `json:"identifier" binding:"required"` // Can be username or email
	Password   string `json:"password" binding:"required"`
}

type UserAPI struct {
	DB UserDatabase
}

func (api *UserAPI) Register(ctx *gin.Context) {
	user := sql.CreateUser{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if existingUser, err := api.DB.GetUserByName(user.Name); err == nil && existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	if existingUser, err := api.DB.GetUserByEmail(user.Email); err == nil && existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	newUser := &sql.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: password.GeneratePasswordHash(user.Password),
		Admin:    false,
	}

	if err := api.DB.CreateUser(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (api *UserAPI) Login(ctx *gin.Context) {
	user := LoginUser{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Try to find the user by name or email
	existingUser := &sql.User{}
	var err error
	if existingUser, err = api.DB.GetUserByName(user.Identifier); err != nil || existingUser == nil {
		if existingUser, err = api.DB.GetUserByEmail(user.Identifier); err != nil || existingUser == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
	}

	if !password.CompareHashPassword(existingUser.Password, []byte(user.Password)) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	sessionManager := ctx.MustGet("sessionManager").(*auth.SessionManager)
	if err := sessionManager.CreateSession(ctx.Writer, ctx.Request, existingUser.ID, existingUser.Name, existingUser.Admin); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	// Return user details in the response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": map[string]interface{}{
			"id":    existingUser.ID,
			"name":  existingUser.Name,
			"email": existingUser.Email,
			"admin": existingUser.Admin,
		},
	})
}

func (u *UserAPI) Logout(ctx *gin.Context) {
	sessionManager := ctx.MustGet("sessionManager").(*auth.SessionManager)
	if err := sessionManager.DestroySession(ctx.Writer, ctx.Request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log out"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
