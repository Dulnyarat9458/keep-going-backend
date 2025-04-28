package validators

import (
	"fmt"
	"keep_going/models"
	"net/mail"
	"strings"
)

func ValidateUserInput(input models.User) []map[string]string {
	var errors []map[string]string

	if input.FirstName == "" {
		errors = append(errors, map[string]string{
			"field": "first_name",
			"error": "First Name is required",
		})
	}
	if input.LastName == "" {
		errors = append(errors, map[string]string{
			"field": "last_name",
			"error": "Last Name is required",
		})
	}
	if input.Email == "" {
		errors = append(errors, map[string]string{
			"field": "email",
			"error": "Email is required",
		})
	} else {
		_, err := mail.ParseAddress(input.Email)
		if err != nil {
			errors = append(errors, map[string]string{
				"field": "email",
				"error": "Invalid email format",
			})
		}
	}

	if len(input.Password) < 8 {

		fmt.Println("condition password")

		errors = append(errors, map[string]string{
			"field": "password",
			"error": "Password must be at least 8 characters long",
		})
	}

	return errors
}

func ParseDatabaseError(err error) []map[string]string {
	var errors []map[string]string

	errMsg := err.Error()

	if err == nil {
		return errors
	}

	if strings.Contains(errMsg, "duplicate key value violates unique constraint") {
		if strings.Contains(errMsg, "uni_users_email") {
			errors = append(errors, map[string]string{
				"field": "email",
				"error": "Email Already used",
			})
		}
	} else {
		errors = append(errors, map[string]string{
			"field": "database",
			"error": errMsg,
		})
	}

	return errors
}
