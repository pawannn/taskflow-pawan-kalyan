package domain

import (
	"errors"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrProjectNotFound = errors.New("project not found")
	ErrTaskNotFound    = errors.New("task not found")
)
