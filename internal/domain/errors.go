package domain

const (
	ErrUserAlreadyExists = "user already exists"
	ErrUserNotFound      = "User not found"
	ErrIncorrectPassword = "incorrect password"
	ErrLogin             = "Unable to login at the moment"
	ErrRegister          = "Unable to register at the moment"

	ErrProjectNotFound = "project not found"
	ErrCreateProject   = "unable to create project at the moment"
	ErrFetchProject    = "unable to fetch user projects at the moment"
	ErrUpdateProject   = "unable to update project at the moment"
	ErrDeleteProject   = "unable to delete project at the moment"

	ErrTaskNotFound        = "task not found"
	ErrFetchTask           = "unable to fetch tasks at the moment"
	ErrCreateTask          = "unable to create tasks at the moment"
	ErrRequiredTaskTitle   = "task title is required"
	ErrInvalidTaskPriority = "invalid task priority"
	ErrInvalidTaskStatus   = "invalid task status"
	ErrUpdateTask          = "unable to update task at the moment"
	ErrDeleteTask          = "unable to delete task at the moment"

	ErrForbidded        = "unauthorized action"
	ErrUnAuthorized     = "unauthenticated"
	ErrNotFound         = "not found"
	ErrValidationFailed = "validation failed"
	ErrInvalidReqBody   = "invalid request body"
	ErrInternalError    = "internal server error"
	ErrBadRequest       = "bad request"
	ErrConflict         = "already exist"
)
