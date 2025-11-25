package customError

import "errors"

var (
	ErrUserAlreadyExist = errors.New("user already exists")
	UnExpectedError     = errors.New("unexpected error")
	ErrUserNotFound     = errors.New("user not found for given email")
)
