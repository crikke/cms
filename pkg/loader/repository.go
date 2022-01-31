package loader

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type Repository interface {
	GetContent(ctx context.Context, id uuid.UUID) (contentData, error)
}

type contentData struct {
	ID       uuid.UUID
	ParentID uuid.UUID
	Version  int
	Created  time.Time
	Updated  time.Time
	Data     map[int]contentVersion
}
type contentVersion struct {
	Properties []contentProperty
	Name       map[language.Tag]string
	URLSegment map[language.Tag]string
}

type contentProperty struct {
	ID        uuid.UUID
	Name      string
	Type      string
	Localized bool
	Value     map[language.Tag]interface{}
}
