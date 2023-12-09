package errors

import (
	"errors"
	"fmt"
)

// storage errors
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

var ErrOrderIsExistThisUser = errors.New("this order is exist the user")
var ErrOrderIsExistAnotherUser = errors.New("this order is exist another user")
