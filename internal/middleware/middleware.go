package middleware

import (
	"encore.dev/beta/errs"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return &errs.Error{
				Code:    errs.Internal,
				Message: "internal validation error",
			}
		}
		return &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "validation failed: " + err.Error(),
		}
	}

	return nil
}
