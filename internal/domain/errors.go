package domain

const (
	ErrUserAlreadyExists  = "user already exists"
	ErrInvalidCredentials = "invalid credentials"
	ErrIncorrectPassword  = "incorrect password"
	ErrLogin              = "Unable to login at the moment"
	ErrRegister           = "Unable to register at the moment"

	ErrProjectNotFound = "project not found"
	ErrCreateProject   = "unable to create project at the moment"
	ErrFetchProject    = "unable to fetch user projects at the moment"
	ErrUpdateProject   = "unable to update project at the moment"
	ErrDeleteProject   = "unable to delete project at the moment"

	ErrTaskNotFound = "task not found"
	ErrFetchTask    = "unable to fetch tasks at the moment"

	ErrForbidded = "You are not allowed to perform this operation"
)
