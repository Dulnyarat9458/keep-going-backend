package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func HandleBindError(err error) []map[string]string {
	var ve validator.ValidationErrors
	var out []map[string]string

	if errors.As(err, &ve) {
		for _, fe := range ve {
			out = append(out, map[string]string{
				"field": ToSnakeCase(fe.Field()),
				"error": ValidationMessage(fe),
			})
		}
	} else {
		out = append(out, map[string]string{
			"field": "json",
			"error": err.Error(),
		})
	}

	return out
}
