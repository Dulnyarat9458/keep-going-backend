package validators

import "time"

type AddHabitInput struct {
	Title     string     `json:"title" binding:"required"`
	StartDate *time.Time `json:"start_date" binding:"required"`
}

type EditHabitOutput struct {
	Title string `json:"title" binding:"required"`
}
