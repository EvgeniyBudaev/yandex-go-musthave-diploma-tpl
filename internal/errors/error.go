package errors

import (
	"errors"
	"fmt"
)

type OrderIsExistAnotherUserError struct {
	Message string
	Err     error
}

func NewOrderIsExistAnotherUserError(msg string, err error) error {
	return &OrderIsExistAnotherUserError{
		Message: msg,
		Err:     err,
	}
}

func (e *OrderIsExistAnotherUserError) Error() string {
	return fmt.Sprintf("[%s] %v", e.Message, e.Err)
}

type OrderIsExistThisUserError struct {
	Message string
	Err     error
}

func NewOrderIsExistThisUserError(msg string, err error) error {
	return &OrderIsExistThisUserError{
		Message: msg,
		Err:     err,
	}
}

func (e *OrderIsExistThisUserError) Error() string {
	return fmt.Sprintf("[%s] %v", e.Message, e.Err)
}

var ErrTooManyRequests = errors.New("too many requests")
var ErrNoContent = errors.New("no content")
