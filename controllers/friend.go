package controllers

import (
	"net/http"
	"strconv"
	"time"

	"social-media-api/config"
	"social-media-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FriendRequestInput struct {
	ToUserID uint `json:"to_user_id" binding:"required"`
}

func SendFriendRequest(c *gin.Context) {
	var input FriendRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fromUserID, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)
	friendRequest := models.FriendRequest{
		FromUserID: uint(fromUserID),
		ToUserID:   input.ToUserID,
		Status:     "pending",
	}

	// Rate limit check: Ensure that the user sends no more than 3 friend requests per minute
	var count int64
	timeLimit := time.Now().Add(-time.Minute)
	config.DB.Model(&models.FriendRequest{}).Where("from_user_id = ? AND created_at >= ?", fromUserID, timeLimit).Count(&count)

	if count >= 3 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Friend request limit reached. Please try again later."})
		return
	}

	// Save friend request to the database
	config.DB.Create(&friendRequest)
	c.JSON(http.StatusCreated, friendRequest)
}

func ListFriendRequests(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)
	var requests []models.FriendRequest

	config.DB.Preload("FromUser").Where("to_user_id = ? AND status = ?", uint(userID), "pending").Find(&requests)
	c.JSON(http.StatusOK, requests)
}

func UpdateFriendRequestStatus(c *gin.Context) {
	var input struct {
		Action string `json:"action" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	friendRequestID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var friendRequest models.FriendRequest

	if err := config.DB.First(&friendRequest, friendRequestID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve friend request"})
		}
		return
	}

	if input.Action == "accept" {
		friendRequest.Status = "accepted"
	} else if input.Action == "reject" {
		friendRequest.Status = "rejected"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	config.DB.Save(&friendRequest)
	c.JSON(http.StatusOK, friendRequest)
}

func ListFriends(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)
	var friends []models.User

	// Find accepted friend requests where the user is either the sender or receiver
	config.DB.Raw(`
		SELECT * FROM users
		WHERE id IN (
			SELECT CASE WHEN from_user_id = ? THEN to_user_id ELSE from_user_id END
			FROM friend_requests
			WHERE (from_user_id = ? OR to_user_id = ?) AND status = 'accepted'
		)`, userID, userID, userID).Scan(&friends)

	c.JSON(http.StatusOK, friends)
}
