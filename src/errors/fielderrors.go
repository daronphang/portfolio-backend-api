package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Substitutes tag specified in struct for error message if thrown
func msgForTag(fe validator.FieldError) string {
	var msg string
	switch fe.Tag() {
	case "required":
		msg = "field is required"
	case "email":
		msg = "invalid email"
	case "boolean":
		msg = "field must be boolean"
	case "number":
		msg = "field must be number"
	case "numeric":
		msg = "field must be numeric"
	case "url":
		msg = "field must be url"
	case "len":
		msg = fmt.Sprintf("field must have length of %s", fe.Param())
	case "oneof":
		msg = fmt.Sprintf("field must be one of %s", fe.Param())
	case "gt":
		msg = fmt.Sprintf("field must be greater than %s", fe.Param())
	case "lt":
		msg = fmt.Sprintf("field must be less than %s", fe.Param())
	case "max":
		msg = fmt.Sprintf("maximum allowed is %s", fe.Param())
	case "min":
		msg = fmt.Sprintf("minimum allowed is %s", fe.Param())
	case "eq":
		msg = fmt.Sprintf("field must be equal to %s", fe.Param())
	case "ne":
		msg = fmt.Sprintf("field must not be equal to %s", fe.Param())
	default:
		msg = fe.Error()
	}
	return msg
}

func TranslateErrors(ve *validator.ValidationErrors) (map[string]string) {
	errors := make(map[string]string, len(*ve))
	for _, fe := range *ve {
		errors[fe.Field()] = msgForTag(fe)
	}
	return errors
}