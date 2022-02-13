package user

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID
	Name string
}

func (u User) GetID() uuid.UUID {
	return u.ID
}

func (u User) GetName() string {
	return u.Name
}
