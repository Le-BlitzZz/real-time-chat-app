package api

import (
	"net/http"

	"github.com/Le-BlitzZz/real-time-chat-app/model/sql"
	"github.com/gin-gonic/gin"
)

func (api *UserAPI) GetFriends(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	friends, err := api.DB.GetFriends(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve friends"})
		return
	}

	ctx.JSON(http.StatusOK, friends)
}

func (api *UserAPI) SendFriendRequest(ctx *gin.Context) {
	receiver := sql.FriendRequestReceiverId{}
	if err := ctx.ShouldBindJSON(&receiver); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	senderID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := api.DB.CreateFriendRequest(senderID.(uint), receiver.ReceiverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not send friend request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "friend request sent"})
}

func (api *UserAPI) AcceptFriendRequest(ctx *gin.Context) {
	friendRequest := sql.FriendRequestAcceptReject{}
	if err := ctx.ShouldBindJSON(&friendRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := api.DB.AcceptFriendRequest(friendRequest.RequestID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not accept friend request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
}

func (api *UserAPI) GetFriendRequests(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	requests, err := api.DB.GetFriendRequests(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve friend requests"})
		return
	}

	ctx.JSON(http.StatusOK, requests)
}

func (api *UserAPI) RejectFriendRequest(ctx *gin.Context) {
	friendRequest := sql.FriendRequestAcceptReject{}
	if err := ctx.ShouldBindJSON(&friendRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := api.DB.RejectFriendRequest(friendRequest.RequestID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not reject friend request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "friend request rejected"})
}
