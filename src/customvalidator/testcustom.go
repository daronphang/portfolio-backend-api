package customvalidator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var TestCustom validator.Func = func(fl validator.FieldLevel) bool {
	name, ok := fl.Field().Interface().(string)
	if ok {
		re := regexp.MustCompile("^[a-z]{5}$")
		if match := re.MatchString(name); match {
			return true
		}
	}
	return false
}
