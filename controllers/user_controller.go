package controllers

import (
	database "keep_going/databases"
	"keep_going/models"
	"keep_going/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func SignUp(c *gin.Context) {
	var input models.User
	var allErrors []map[string]string
	err := c.ShouldBindJSON(&input)

	if err != nil {
		allErrors = append(allErrors, map[string]string{
			"field": "json",
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, allErrors)
		return
	}

	hashedPassword, err_hash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err_hash != nil {

		allErrors = append(allErrors, map[string]string{
			"field": "json",
			"error": "Failed to encrypt password",
		})
		c.JSON(http.StatusInternalServerError, allErrors)
		return
	}

	input.Password = string(hashedPassword)

	inputErrors := validators.ValidateUserInput(input)
	allErrors = append(allErrors, inputErrors...)

	if len(allErrors) == 0 {
		result := database.DB.Create(&input)

		if result != nil && result.Error != nil {
			dbErrors := validators.ParseDatabaseError(result.Error)
			allErrors = append(allErrors, dbErrors...)
		}
	}

	if len(allErrors) > 0 {
		c.JSON(http.StatusBadRequest, allErrors)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "User registered successfully!",
		"username": input.Email,
	})
}
