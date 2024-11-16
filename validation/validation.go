package validation

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FormError struct {
	Field   string
	Message string
	Value   any
}

type FormErrors []FormError

func (f FormErrors) HasField(field string) bool {
	var has bool
	for _, err := range f {
		if err.Field == field {
			has = true
		}
	}
	return has
}

func (f FormErrors) GetMessage(field string) string {
	var message string
	for _, err := range f {
		if err.Field == field {
			message = err.Message
		}
	}
	return message
}

// TODO: Implement other value types (number, bool, slices).
func (f FormErrors) GetValue(field string) string {
	var value any
	for _, err := range f {
		if err.Field == field {
			value = err.Value
		}
	}
	if s, ok := value.(string); ok {
		return s
	}
	return ""
}

// TODO: Implement other error messages.
func Message(tag string) string {
	switch tag {
	case "required":
		return "Field is required."
	default:
		return ""
	}
}

func New() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return v
}

func Default() *FormErrors {
	var errors FormErrors
	return &errors
}

func Errors(data any, err error) *FormErrors {
	var errors FormErrors

	v := reflect.ValueOf(data)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		errors = append(errors, FormError{
			Field: field.Tag.Get("form"),
			Value: value.Interface(),
		})
	}

	for _, err := range err.(validator.ValidationErrors) {
		for i, formErr := range errors {
			if err.Field() == formErr.Field {
				fr := FormError{
					Field:   formErr.Field,
					Message: Message(err.Tag()),
					Value:   formErr.Value,
				}
				log.Println(fr)
				errors[i] = fr
			}
		}
	}

	return &errors
}
