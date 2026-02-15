package validator

import (
	"tms-core-service/internal/domain/errs"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct and returns ValidationErrors
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	valErrs := make(errs.ValidationErrors)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := e.Field()
			valErrs[field] = append(valErrs[field], e.Tag())
		}
	} else {
		return err
	}

	return valErrs
}
