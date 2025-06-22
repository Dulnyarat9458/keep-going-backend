package utils

import "github.com/go-playground/validator/v10"

func ValidationCode(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "email":
		return "invalid_email"
	case "min":
		return "min_length"
	}
	return "invalid"
}
