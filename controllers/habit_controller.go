package controllers

import (
	"fmt"
	"keep_going/databases"
	"keep_going/models"
	"keep_going/utils"

	"keep_going/validators"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

		result := databases.DB.Where(&models.HabitTracker{UserID: user.ID,
			ID: habitId,
		}).First(&habitTracker)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Habit not found",
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
		var input validators.AddHabitInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, utils.HandleBindError(err))
			return
		}

		habitTrackers.UserID = user.ID
		habitTrackers.Title = input.Title
		habitTrackers.LastResetDate = *input.StartDate
		habitTrackers.StartDate = *input.StartDate

		result := databases.DB.Create(&habitTrackers)

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

		result := databases.DB.Where(&models.HabitTracker{
			UserID: user.ID,
			ID:     habitId,
		}).First(&habitTracker)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "habit not found"})
			return
		}

		var input validators.EditHabitOutput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, utils.HandleBindError(err))
			return
		}

		habitTracker.Title = input.Title
		habitTracker.StartDate = *input.StartDate
		habitTracker.LastResetDate = *input.LastResetDate

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

func HabitDelete(c *gin.Context) {
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

		result := databases.DB.Where(&models.HabitTracker{UserID: user.ID,
			ID: habitId,
		}).Delete(&habitTracker)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "habit not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "habit Deleted",
			"data":    habitTracker,
		})
		return

	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "something went wrong",
	})
	return
}
