package controllers

import (
	"fmt"
	"keep_going/databases"
	"keep_going/models"
	"keep_going/validators"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HabitList(c *gin.Context) {
	var habitTrackers []models.HabitTracker

	if u, exists := c.Get("user"); exists {
		user := u.(models.User)

		result := databases.DB.Where(&models.HabitTracker{UserID: user.ID}).Find(&habitTrackers)

		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "something went wrong",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok", "data": habitTrackers,
		})
		return

	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "something went wrong",
	})
	return

}

func HabitDetail(c *gin.Context) {
	var habitTracker models.HabitTracker
	habitIdStr := c.Param("id")
	habitIdUint64, err := strconv.ParseUint(habitIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid habit ID"})
		return
	}

	if u, exists := c.Get("user"); exists {
		habitId := uint(habitIdUint64)
		user := u.(models.User)

		result := databases.DB.Where(&models.HabitTracker{UserID: user.ID, Model: gorm.Model{
			ID: habitId,
		}}).First(&habitTracker)

		fmt.Println(result.Error)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "something went wrong",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"data":    habitTracker,
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
		var habitTrackers models.HabitTracker
		var allErrors []map[string]string
		err := c.ShouldBindJSON(&habitTrackers)

		if err != nil {
			allErrors = append(allErrors, map[string]string{
				"field": "json",
				"error": err.Error(),
			})
			c.JSON(http.StatusBadRequest, allErrors)
			return
		}

		inputErrors := validators.ValidateHabitInput(habitTrackers)

		habitTrackers.UserID = user.ID
		habitTrackers.LastResetDate = habitTrackers.StartDate
		result := databases.DB.Create(&habitTrackers)

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

func HabitEdit(c *gin.Context) {
	var habitTracker models.HabitTracker
	habitIdStr := c.Param("id")
	habitIdUint64, err := strconv.ParseUint(habitIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid habit ID"})
		return
	}

	if u, exists := c.Get("user"); exists {
		habitId := uint(habitIdUint64)
		user := u.(models.User)

		result := databases.DB.Where(&models.HabitTracker{UserID: user.ID, Model: gorm.Model{
			ID: habitId,
		}}).First(&habitTracker)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "habit not found"})
			return
		}

		var input struct {
			Title         *string    `json:"title"`
			StartDate     *time.Time `json:"start_date"`
			LastResetDate *time.Time `json:"last_reset_date"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON input"})
			return
		}

		if input.Title != nil {
			habitTracker.Title = *input.Title
		}
		if input.StartDate != nil {
			habitTracker.StartDate = *input.StartDate
		}
		if input.LastResetDate != nil {
			habitTracker.LastResetDate = *input.LastResetDate
		}

		databases.DB.Save(&habitTracker)

		c.JSON(http.StatusOK, gin.H{
			"message": "habit updated",
			"data":    habitTracker,
		})
		return

	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "something went wrong",
	})
	return

}
