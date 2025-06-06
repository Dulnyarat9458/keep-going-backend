package utils

import "github.com/go-playground/validator/v10"

func ValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	}
	return fe.Field() + " is not valid"
}
