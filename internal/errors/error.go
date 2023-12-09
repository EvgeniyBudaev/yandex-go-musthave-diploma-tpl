package errors

import (
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
