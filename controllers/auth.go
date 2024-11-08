package controllers

import (
	"net/http"
	"social-media-api/config"
	"social-media-api/models"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Username: input.Username, Email: input.Email, Password: input.Password}
	config.DB.Create(&user)

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	config.DB.Where("username = ? AND password = ?", input.Username, input.Password).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
