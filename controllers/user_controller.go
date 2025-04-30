package controllers

import (
	database "keep_going/databases"
	"keep_going/models"
	jwt "keep_going/utils"
	"keep_going/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignIn(c *gin.Context) {
	var user models.User
	var allErrors []map[string]string
	err := c.ShouldBindJSON(&user)
	if err != nil {
		allErrors = append(allErrors, map[string]string{
			"field": "json",
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, allErrors)
		return
	}

	input_password := user.Password

	result := database.DB.Where("email = ?", user.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid credentials",
			"field":   "email",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input_password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid credentials",
			"field":   "email",
		})
		return
	}

	jwt, err := jwt.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong",
			"field":   "non_field",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"token":   jwt,
	})
	return

}

func SignUp(c *gin.Context) {
	var user models.User
	var allErrors []map[string]string
	err := c.ShouldBindJSON(&user)

	if err != nil {
		allErrors = append(allErrors, map[string]string{
			"field": "json",
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, allErrors)
		return
	}

	inputErrors := validators.ValidateUserSignUpInput(user)

	hashedPassword, err_hash := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err_hash != nil {

		allErrors = append(allErrors, map[string]string{
			"field": "json",
			"error": "Failed to encrypt password",
		})
		c.JSON(http.StatusInternalServerError, allErrors)
		return
	}

	user.Password = string(hashedPassword)

	allErrors = append(allErrors, inputErrors...)

	if len(allErrors) == 0 {
		user.Role = "user"
		result := database.DB.Create(&user)

		if result != nil && result.Error != nil {
			dbErrors := validators.ParseDatabaseUserSignUpError(result.Error)
			allErrors = append(allErrors, dbErrors...)
		}
	}

	if len(allErrors) > 0 {
		c.JSON(http.StatusBadRequest, allErrors)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully!",
		"email":   user.Email,
	})
}
