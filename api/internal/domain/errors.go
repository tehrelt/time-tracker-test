package domain

import "errors"

var (
	ErrNotImplemented     = errors.New("not implemented")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyWorking = errors.New("user already working")
	ErrUserNotWorking     = errors.New("user not working")
)
