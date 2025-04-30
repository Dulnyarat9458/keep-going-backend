package controllers

import (
	"fmt"
	"keep_going/databases"
	"keep_going/models"
	"keep_going/validators"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HabitList(c *gin.Context) {
	var habit_trackers []models.HabitTracker

	if u, exists := c.Get("user"); exists {
		user := u.(models.User)

		result := databases.DB.Where(&models.HabitTracker{UserID: user.ID}).Find(&habit_trackers)

		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "something went wrong",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok", "data": habit_trackers,
		})
		return

	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "something went wrong",
	})
	return

}

func AddNewHabit(c *gin.Context) {
	if u, exists := c.Get("user"); exists {
		user := u.(models.User)
		var habit_tracker models.HabitTracker
		var allErrors []map[string]string
		err := c.ShouldBindJSON(&habit_tracker)

		if err != nil {
			allErrors = append(allErrors, map[string]string{
				"field": "json",
				"error": err.Error(),
			})
			c.JSON(http.StatusBadRequest, allErrors)
			return
		}

		inputErrors := validators.ValidateHabitInput(habit_tracker)

		habit_tracker.UserID = user.ID
		habit_tracker.LastResetDate = habit_tracker.StartDate
		result := databases.DB.Create(&habit_tracker)

		allErrors = append(allErrors, inputErrors...)

		if len(allErrors) > 0 {
			c.JSON(http.StatusBadRequest, allErrors)
			return
		}

		if result.Error != nil {
			fmt.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong",
				"field":   "non_field",
			})
			return
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "not found user",
		})
		return
	}
}
