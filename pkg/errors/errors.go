package errors

// TODO once this errors package available in salt, we need to remove this and use the one in salt instead

import (
	"errors"
	"fmt"
	"strings"
)

// These aliased values are added to avoid conflicting imports of standard `errors`
// package and this `errors` package where these functions are needed.
var (
	Is  = errors.Is
	As  = errors.As
	New = errors.New
)

// Common error categories. Use `ErrX.WithXXX()` to clone and add context.
var (
	ErrInvalid     = Error{Code: "bad_request", Message: "Request is not valid"}
	ErrNotFound    = Error{Code: "not_found", Message: "Requested entity not found"}
	ErrConflict    = Error{Code: "conflict", Message: "An entity with conflicting identifier exists"}
	ErrInternal    = Error{Code: "internal_error", Message: "Some unexpected error occurred"}
	ErrUnsupported = Error{Code: "unsupported", Message: "Requested feature is not supported"}
)

// Error represents any error returned by the Entropy components along with any
// relevant context.
type Error struct {
	Code    string `json:"code"`
	Cause   string `json:"cause,omitempty"`
	Message string `json:"message"`
}

// WithCausef returns clone of err with the cause added. Use this when
// you need to provide description of the underlying technical root-cause
// which may be written in log for debugging purposes. Cause will be shown
// to the user only when the Message is empty.
func (err Error) WithCausef(format string, args ...interface{}) Error {
	cloned := err
	cloned.Cause = fmt.Sprintf(format, args...)
	return cloned
}

// WithMsgf returns a clone of the error with message set. Use this when
// you need to provide a custom message that should be shown to the user.
// If the message is set to empty string, cause will be displayed to the
// user.
func (err Error) WithMsgf(format string, args ...interface{}) Error {
	cloned := err
	cloned.Message = fmt.Sprintf(format, args...)
	return cloned
}

// Is checks if 'other' is of type Error and has the same code.
// See https://blog.golang.org/go1.13-errors.
func (err Error) Is(other error) bool {
	if oe, ok := other.(Error); ok { // nolint
		return oe.Code == err.Code
	}

	// unknown error types are considered as internal errors.
	return err.Code == ErrInternal.Code
}

func (err Error) Error() string {
	if err.Message != "" {
		return strings.ToLower(err.Message)
	}
	return fmt.Sprintf("%s: %s", err.Code, err.Cause)
}

// Errorf returns a formatted error similar to `fmt.Errorf` but uses the
// Error type defined in this package. returned value is equivalent to
// ErrInternal (i.e., errors.Is(retVal, ErrInternal) = true).
func Errorf(format string, args ...interface{}) error {
	return ErrInternal.WithMsgf(format, args...)
}

// OneOf checks if the error is one of the 'others'.
func OneOf(err error, others ...error) bool {
	for _, other := range others {
		if errors.Is(err, other) {
			return true
		}
	}
	return false
}

// E converts any given error to the Error type. Unknown are converted
// to ErrInternal.
func E(err error) Error {
	var e Error
	if errors.As(err, &e) {
		return e
	}
	return ErrInternal.WithCausef(err.Error())
}

// Verbose returns a verbose error value.
func Verbose(err error) error {
	var e Error
	if errors.As(err, &e) {
		return e.WithMsgf("%s: %s (cause: %s)", e.Code, e.Message, e.Cause)
	}
	return err
}
