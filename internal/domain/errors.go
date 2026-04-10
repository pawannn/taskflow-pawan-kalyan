package domain

import "fmt"

var ErrUserAlreadyExists = fmt.Errorf("user already exists")
var ErrInvalidCredentials = fmt.Errorf("invalid credentials")
