package controllers

import (
	"bytes"

	"fmt"
	"html/template"
	"keep_going/databases"
	"keep_going/models"
	"keep_going/utils"
	"keep_going/validators"

	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mail.v2"
)

func SignIn(c *gin.Context) {
	var user models.User
	var input validators.SignInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleBindError(err))
		return
	}

	inputPassword := input.Password

	result := databases.DB.Where("email = ?", input.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_credentials",
			"field": "email",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_credentials",
			"field": "email",
		})
		return
	}

	jwt, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "something went wrong",
			"field": "non_field",
		})
	}

	c.SetCookie(
		"access_token", // name
		jwt,            // value
		3600*24,        // maxAge in seconds (1 day)
		"/",            // path
		"",             // domain ("" = current domain)
		true,           // secure (true = only HTTPS)
		true,           // httpOnly
	)
	return
}

func SignOut(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.JSON(200, gin.H{"message": "logged out"})
	return
}

func SignUp(c *gin.Context) {
	var input validators.SignUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleBindError(err))
		return
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		c.JSON(http.StatusInternalServerError, []map[string]string{
			{
				"field": "json",
				"error": "Failed to encrypt password",
			},
		})
		return
	}

	var existingUser models.User

	if err := databases.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, []map[string]string{
			{
				"field": "email",
				"error": "exists_email",
			},
		})
		return
	}

	existingUser.FirstName = input.FirstName
	existingUser.LastName = input.LastName
	existingUser.Email = input.Email
	existingUser.Password = string(hashedPassword)
	existingUser.Role = "user"
	result := databases.DB.Create(&existingUser)

	if result != nil && result.Error != nil {
		c.JSON(http.StatusBadRequest, []map[string]string{
			{
				"field": "database",
				"error": result.Error.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully!",
		"email":   existingUser.Email,
	})
	return
}

func ForgetPassword(c *gin.Context) {
	var input validators.ForgetPasswordInput

	var user models.User
	var resetToken models.ResetToken

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleBindError(err))
		return
	}

	tmpl, err := template.ParseFiles("templates/reset_password.html")
	if err != nil {
		log.Fatal("Template error:", err)
	}

	token, err := utils.GenerateResetToken()
	if err != nil {
		log.Fatal("token error:", err)
	}

	result_user := databases.DB.Where("email = ?", input.Email).First(&user)

	if result_user.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_email",
			"field": "email",
		})
		return
	}

	resetToken.UserID = user.ID
	resetToken.Token = token
	resetToken.ExpiresAt = time.Now().Add(30 * time.Minute)

	result_token := databases.DB.Create(&resetToken)

	if result_token.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "fail",
			"field": "non_field",
		})
		return
	}
	frontendWeb := os.Getenv("FRONTEND_WEB")

	data := struct {
		ResetLink string
	}{
		ResetLink: fmt.Sprintf("%s/reset-password?token=%s", frontendWeb, token),
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Fatal("Execute template error:", err)
	}

	mailSender := os.Getenv("MAIL_SENDER")
	mailHost := os.Getenv("MAIL_HOST")
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	m := mail.NewMessage()

	m.SetHeader("From", mailSender)
	m.SetHeader("To", input.Email)
	m.SetHeader("Subject", "Reset Password")
	m.SetBody("text/html", body.String())

	if err != nil {
		panic(err)
	}

	d := mail.NewDialer(
		mailHost,
		mailPort,
		mailUsername,
		mailPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{"message": "OK"})
	return
}

func ResetPassword(c *gin.Context) {
	var user models.User
	var resetToken models.ResetToken
	var input validators.ResetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.HandleBindError(err))
		return
	}

	tokenStr := input.Token

	resultToken := databases.DB.Where("token = ?", tokenStr).First(&resetToken)

	if resultToken.Error != nil {
		c.JSON(400, gin.H{"error": "invalid_token"})
		return
	}

	inputPassword := input.Password

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	if errHash != nil {
		c.JSON(400, gin.H{"error": "fail"})
		return
	}

	databases.DB.First(&user, resetToken.UserID)
	user.Password = string(hashedPassword)
	databases.DB.Unscoped().Delete(&resetToken)
	databases.DB.Save(&user)

	c.JSON(200, gin.H{"message": "ok"})
	return
}

func MyUserInfo(c *gin.Context) {
	var user models.User

	if u, exists := c.Get("user"); exists {
		fmt.Println(u)
		result := databases.DB.First(&user, u)

		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "fail",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data": gin.H{
				"first_name": user.FirstName,
				"last_name":  user.FirstName,
				"email":      user.FirstName,
			},
		})
		return

	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": "fail",
	})
	return
}
