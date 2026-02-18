package errs

import (
	"encore.dev/beta/errs"
)

// Wrap returns an encounter error with the given code and message, wrapping the original error.
func Wrap(code errs.ErrCode, err error, msg string) error {
	return &errs.Error{
		Code:    code,
		Message: msg,
		Meta: map[string]interface{}{
			"cause": err.Error(),
		},
	}
}

// New returns a new encore error with the given code and message.
func New(code errs.ErrCode, msg string) error {
	return &errs.Error{
		Code:    code,
		Message: msg,
	}
}

// BadRequest returns a standardized InvalidArgument error
func BadRequest(msg string) error {
	return New(errs.InvalidArgument, msg)
}

// Internal returns a standardized Internal error
func Internal(err error, msg string) error {
	return Wrap(errs.Internal, err, msg)
}

// NotFound returns a standardized NotFound error
func NotFound(msg string) error {
	return New(errs.NotFound, msg)
}

// Unauthorized returns a standardized Unauthenticated error
func Unauthorized(msg string) error {
	return New(errs.Unauthenticated, msg)
}
