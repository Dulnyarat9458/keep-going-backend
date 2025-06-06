package controllers

import (
	"bytes"
	"errors"
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
	"github.com/go-playground/validator/v10"
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

	inputPassword := user.Password

	result := databases.DB.Where("email = ?", user.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid credentials",
			"field":   "email",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)); err != nil {
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

	var input validators.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		var out []map[string]string

		if errors.As(err, &ve) {
			for _, fe := range ve {
				out = append(out, map[string]string{
					"field": utils.ToSnakeCase(fe.Field()),
					"error": utils.ValidationMessage(fe),
				})
			}
		} else {
			out = append(out, map[string]string{
				"field": "json",
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusBadRequest, out)
		return
	}

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
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
				"error": "Email already exists",
			},
		})
		return
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Email = input.Email
	user.Password = string(hashedPassword)
	user.Role = "user"
	result := databases.DB.Create(&user)

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
		"email":   user.Email,
	})
	return
}

func ForgetPassword(c *gin.Context) {

	var user models.User
	var resetToken models.ResetToken

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

	resetToken.UserID = user.ID
	resetToken.Token = token
	resetToken.ExpiresAt = time.Now().Add(30 * time.Minute)

	result_token := databases.DB.Create(&resetToken)

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

	mailSender := os.Getenv("MAIL_SENDER")
	mailHost := os.Getenv("MAIL_HOST")
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	m := mail.NewMessage()

	m.SetHeader("From", mailSender)
	m.SetHeader("To", user.Email)
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

	resultToken := databases.DB.Where("token = ?", tokenStr).First(&resetToken)

	if resultToken.Error != nil {
		c.JSON(400, gin.H{"error": "invalid token"})
		return
	}

	inputPassword := input.Password

	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	if errHash != nil {
		c.JSON(400, gin.H{"message": "BAD"})
		return
	}

	databases.DB.First(&user, resetToken.UserID)
	user.Password = string(hashedPassword)
	databases.DB.Unscoped().Delete(&resetToken)
	databases.DB.Save(&user)

	c.JSON(200, gin.H{"message": "RESET COMPLETE"})
	return
}
