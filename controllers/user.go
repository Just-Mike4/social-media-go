package controllers

import (
	"net/http"
	"social-media-api/config"
	"social-media-api/models"

	"github.com/gin-gonic/gin"
)

func SearchUsers(c *gin.Context) {
	var users []models.User
	search := c.Query("search")
	pageSize := 10

	config.DB.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%").
		Limit(pageSize).Find(&users)

	c.JSON(http.StatusOK, users)
}
