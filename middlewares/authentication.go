package middlewares

import (
	"keep_going/databases"
	"keep_going/models"
	"keep_going/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User
		err := c.ShouldBindJSON(&user)

		token, err := c.Cookie("access_token")
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: access_token cookie missing",
			})
			return
		}

		email, err := utils.ParseJWT(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized:can't authorize access_token",
			})
		}

		result := databases.DB.Where("email = ?", email).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "not found user",
			})
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
