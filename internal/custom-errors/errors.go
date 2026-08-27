package customerrors

import "errors"

var (
	ErrNotFound   = errors.New("todo not found")
	ErrTitleEmpty = errors.New("title cannot be empty")
	ErrBadRequest = errors.New("invalid json")
)
