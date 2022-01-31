package loader

import (
	"context"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/content"
	"github.com/crikke/cms/pkg/locale"
	"github.com/google/uuid"
	"golang.org/x/text/language"
)

/*
	Loader is responsible for retreiving content from database and transforms it to content.Content
*/

type Loader interface {
	GetContent(ctx context.Context, id uuid.UUID) (content.Content, error)
	GetChildNodes(ctx context.Context, id uuid.UUID) ([]content.Content, error)
}

type loader struct {
	db  Repository
	cfg config.Configuration
}

func NewLoader(db Repository, cfg config.Configuration) loader {
	return loader{db, cfg}
}

func (l *loader) GetContent(ctx context.Context, id uuid.UUID) (content.Content, error) {

	content, err := l.db.GetContent(ctx, id)

	if err != nil {
		panic(err)
	}

	t := locale.FromContext(ctx)

	if t == language.Und {
		t = l.cfg.Languages[0]
	}
	return convert(
		content,
		t,
		l.cfg.Languages[0],
		0)
}

// Converts a db entity to content.Content
func convert(entity contentData, lang language.Tag, fallbackLang language.Tag, version int) (content.Content, error) {

	result := content.Content{
		ID:       entity.ID,
		ParentID: entity.ParentID,
		Created:  entity.Created,
		Updated:  entity.Updated,
	}

	data, exist := entity.Data[version]

	if !exist {
		return content.Content{}, ContentError{entity.ID, version, "not found"}
	}

	result.URLSegment = data.URLSegment[lang]
	result.Name = data.Name[lang]

	for _, prop := range data.Properties {

		localized := prop.Value[fallbackLang]
		if prop.Localized {

			p, exist := prop.Value[lang]

			if exist {
				localized = p
			}
		}

		cp := content.Property{
			ID:        prop.ID,
			Name:      prop.Name,
			Type:      prop.Type,
			Localized: prop.Localized,
			Value:     localized,
		}
		result.Properties = append(result.Properties, cp)
		continue
	}

	return result, nil
}
