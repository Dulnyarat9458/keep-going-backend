package controllers

import (
	database "keep_going/databases"
	"keep_going/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func SignUp(c *gin.Context) {
	var input models.User

	err := c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&input)

	c.JSON(http.StatusOK, gin.H{
		"message":  "User registered successfully!",
		"username": input.Email,
	})
}
