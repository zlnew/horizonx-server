package rest

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(payload any) map[string]string {
	err := validate.Struct(payload)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := strings.ToLower(fieldError.Field())
			switch fieldError.Tag() {
			case "required":
				errors[fieldName] = fmt.Sprintf("The %s field is required.", fieldName)
			case "email":
				errors[fieldName] = "The email must be a valid email address."
			case "min":
				errors[fieldName] = fmt.Sprintf("The %s must be at least %s characters.", fieldName, fieldError.Param())
			default:
				errors[fieldName] = fmt.Sprintf("The %s field is invalid.", fieldName)
			}
		}
	}

	return errors
}
