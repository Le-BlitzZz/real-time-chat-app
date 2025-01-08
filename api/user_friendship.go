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

	response := make([]gin.H, len(friends))
	for i, friend := range friends {
		response[i] = gin.H{
			"id":   friend.ID,
			"name": friend.Name,
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func (api *UserAPI) SendFriendRequest(ctx *gin.Context) {
	senderID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	receiver := sql.FriendRequestReceiver{}
	if err := ctx.ShouldBindJSON(&receiver); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := api.DB.GetUserByName(receiver.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not find requested user"})
		return
	}

	err = api.DB.CreateFriendRequest(senderID.(uint), user.ID)
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

	response := make([]gin.H, len(requests))
	for i, req := range requests {
		response[i] = gin.H{
			"id":   req.ID,
			"name": req.Sender.Name,
		}
	}

	ctx.JSON(http.StatusOK, response)
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
