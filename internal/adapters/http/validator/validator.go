// Package validator
package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(data any) map[string]string
}

type DefaultValidator struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return &DefaultValidator{
		validate: validator.New(),
	}
}

func (v *DefaultValidator) Validate(data any) map[string]string {
	err := v.validate.Struct(data)
	if err == nil {
		return map[string]string{}
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{
			"_error": "invalid payload",
		}
	}

	errors := make(map[string]string)

	for _, e := range validationErrors {
		field := v.resolveFieldName(data, e.Field())
		errors[field] = v.messageFor(e)
	}

	return errors
}

func (v *DefaultValidator) messageFor(e validator.FieldError) string {
	messages := map[string]func(validator.FieldError) string{
		"required": func(e validator.FieldError) string {
			return fmt.Sprintf("%s is required", e.Field())
		},
		"email": func(e validator.FieldError) string {
			return fmt.Sprintf("%s must be a valid email address", e.Field())
		},
		"min": func(e validator.FieldError) string {
			return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
		},
		"eqfield": func(e validator.FieldError) string {
			return fmt.Sprintf("%s must match %s", e.Field(), e.Param())
		},
	}

	if msg, ok := messages[e.Tag()]; ok {
		return msg(e)
	}

	return fmt.Sprintf("%s is invalid", e.Field())
}

func (v *DefaultValidator) resolveFieldName(data any, field string) string {
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if f, ok := t.FieldByName(field); ok {
		tag := f.Tag.Get("json")
		if tag != "" && tag != "-" {
			return strings.Split(tag, ",")[0]
		}
	}

	return strings.ToLower(field)
}
