package laplace

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type LaplaceError error

var (
	ErrYouDoNotHaveAccessToEndpoint LaplaceError = errors.New("you do not have access to this endpoint")
	ErrLimitExceeded                LaplaceError = errors.New("limit exceeded")
	ErrEndpointIsNotActive          LaplaceError = errors.New("endpoint is not active")
	ErrInvalidToken                 LaplaceError = errors.New("invalid token")
	ErrInvalidID                    LaplaceError = errors.New("invalid object id")
)

func WrapError(err error) error {
	if err == nil {
		return nil
	}

	if httpErr, ok := err.(*LaplaceHTTPError); ok {
		getLaplaceError(httpErr)
		return httpErr
	}

	return err
}

func getLaplaceError(httpErr *LaplaceHTTPError) {
	switch httpErr.HTTPStatus {
	case http.StatusForbidden:
		switch httpErr.Message.Message {
		case "you don't have access to this endpoint":
			httpErr.InternalError = ErrYouDoNotHaveAccessToEndpoint
		case "endpoint is not active":
			httpErr.InternalError = ErrEndpointIsNotActive
		}
		if strings.Contains(httpErr.Message.Message, "limit exceeded") {
			httpErr.InternalError = ErrLimitExceeded
		}
	case http.StatusBadRequest:
		switch httpErr.Message.Message {
		case "invalid id":
			httpErr.InternalError = ErrInvalidID
		}
	case http.StatusUnauthorized:
		switch httpErr.Message.Message {
		case "invalid token":
			httpErr.InternalError = ErrInvalidToken
		}
	}
}

type LaplaceHTTPError struct {
	HTTPStatus    int                 `json:"code"`
	Message       LaplaceHTTPErrorMsg `json:"msg"`
	InternalError error               `json:"-"`
}

type LaplaceHTTPErrorMsg struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

func (e *LaplaceHTTPError) Error() string {
	if e.InternalError != nil {
		return fmt.Sprintf("%d: %s (%s)", e.HTTPStatus, e.Message, e.InternalError)
	}
	return fmt.Sprintf("%d: %s", e.HTTPStatus, e.Message)
}

func (e *LaplaceHTTPError) Is(target error) bool {
	if e.InternalError != nil {
		return errors.Is(e.InternalError, target)
	}
	return e.Error() == target.Error()
}

// Cause returns the root cause error
func (e *LaplaceHTTPError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

func (e *LaplaceHTTPError) Unwrap() error {
	if e.InternalError != nil {
		return e.InternalError
	}

	return nil
}

// WithInternalError adds internal error information to the error
func (e *LaplaceHTTPError) WithInternalError(err error) *LaplaceHTTPError {
	e.InternalError = err
	return e
}

func HttpError(httpStatus int, fmtString string, args ...interface{}) *LaplaceHTTPError {
	return &LaplaceHTTPError{
		HTTPStatus: httpStatus,
		Message: LaplaceHTTPErrorMsg{
			Message:   fmt.Sprintf(fmtString, args...),
			ErrorCode: "",
		},
		InternalError: nil,
	}
}
