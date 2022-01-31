package loader

import "github.com/google/uuid"

type ContentError struct {
	ID      uuid.UUID
	Version int
	Message string
}

func (e ContentError) Error() string {
	return e.Message
}
