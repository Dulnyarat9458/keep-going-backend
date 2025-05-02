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

	result := databases.DB.Where("email = ?", user.Email).First(&user)

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

	jwt, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong",
			"field":   "non_field",
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
		result := databases.DB.Create(&user)

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
	return
}

func ForgetPassword(c *gin.Context) {

	var user models.User
	var reset_token models.ResetToken

	err := c.ShouldBindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something went wrong",
			"field":   "non_field",
		})
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

	result_user := databases.DB.Where("email = ?", user.Email).First(&user)

	if result_user.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid email",
			"field":   "email",
		})
		return
	}

	reset_token.UserID = user.ID
	reset_token.Token = token
	reset_token.ExpiresAt = time.Now().Add(30 * time.Minute)

	result_token := databases.DB.Create(&reset_token)

	if result_token.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "something wrong with token save",
			"field":   "non_field",
		})
		return
	}

	data := struct {
		ResetLink string
	}{
		ResetLink: fmt.Sprintf("https://example_frontend.com/reset?token=%s", token),
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Fatal("Execute template error:", err)
	}

	mail_sender := os.Getenv("MAIL_SENDER")
	mail_host := os.Getenv("MAIL_HOST")
	mail_port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mail_username := os.Getenv("MAIL_USERNAME")
	mail_password := os.Getenv("MAIL_PASSWORD")

	m := mail.NewMessage()

	m.SetHeader("From", mail_sender)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Reset Password")
	m.SetBody("text/html", body.String())

	if err != nil {
		panic(err)
	}

	d := mail.NewDialer(
		mail_host,
		mail_port,
		mail_username,
		mail_password,
	)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{"message": "OK"})
	return
}

func ResetPassword(c *gin.Context) {
	var user models.User
	var reset_token models.ResetToken

	type ResetPasswordInput struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	var input ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	tokenStr := input.Token

	fmt.Println("reset_token", reset_token)

	resultToken := databases.DB.Where("token = ?", tokenStr).First(&reset_token)

	if resultToken.Error != nil {
		c.JSON(400, gin.H{"error": "invalid token"})
		return
	}

	inputPassword := input.Password

	hashedPassword, err_hash := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	if err_hash != nil {
		c.JSON(400, gin.H{"message": "BAD"})
		return
	}

	databases.DB.First(&user, reset_token.UserID)
	user.Password = string(hashedPassword)
	databases.DB.Unscoped().Delete(&reset_token)
	databases.DB.Save(&user)

	c.JSON(200, gin.H{"message": "RESET COMPLETE"})
	return
}
