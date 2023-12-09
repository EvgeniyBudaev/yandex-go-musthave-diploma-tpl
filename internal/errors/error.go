package errors

import (
	"errors"
)

var ErrOrderIsExistThisUser = errors.New("this order is exist the user")
var ErrOrderIsExistAnotherUser = errors.New("this order is exist another user")
