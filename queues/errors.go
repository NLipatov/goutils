package queues

import "errors"

var (
	ErrInvalidCapacity = errors.New("capacity must be greater than zero")
)
