package validators

import "keep_going/models"

func ValidateHabitInput(input models.HabitTracker) []map[string]string {
	var errors []map[string]string

	if input.Title == "" {
		errors = append(errors, map[string]string{
			"field": "title",
			"error": "Title is required",
		})
	}

	if input.StartDate.IsZero() {
		errors = append(errors, map[string]string{
			"field": "start_date",
			"error": "Start date is required",
		})
	}
	return errors
}
