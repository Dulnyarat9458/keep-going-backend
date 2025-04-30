package controllers

import (
	"fmt"
	"keep_going/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddNewHabit(c *gin.Context) {
	fmt.Println()
	if u, exists := c.Get("user"); exists {
		user := u.(models.User)
		fmt.Println(user.ID)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "not found user",
		})
		return
	}
}
