package loader

import (
	"context"

	"github.com/crikke/cms/pkg/config"
	"github.com/crikke/cms/pkg/domain"
	"golang.org/x/text/language"
)

/*
	Loader is responsible for retreiving content from database and transforms it to content.Content
*/

type Loader interface {
	GetContent(ctx context.Context, contentReference domain.ContentReference) (domain.Content, error)
	GetChildNodes(ctx context.Context, contentReference domain.ContentReference) ([]domain.Content, error)
}

type loader struct {
	db  Repository
	cfg config.Configuration
}

func NewLoader(db Repository, cfg config.Configuration) loader {
	return loader{db, cfg}
}

func (l *loader) GetContent(ctx context.Context, contentReference domain.ContentReference) (domain.Content, error) {

	// TODO: should probably move getting version logic to database, locale should still be here for now since it contains fallback logic
	content, err := l.db.GetContent(ctx, contentReference.ID)

	if err != nil {
		panic(err)
	}

	t := l.cfg.Languages[0]

	if contentReference.Locale != nil {
		t = *contentReference.Locale
	}

	return convert(
		content,
		t,
		l.cfg.Languages[0],
		0)
}

// Converts a db entity to content.Content
func convert(entity contentData, lang language.Tag, fallbackLang language.Tag, version int) (domain.Content, error) {

	result := domain.Content{
		ID:       domain.ContentReference{ID: entity.ID, Locale: &lang, Version: version},
		ParentID: entity.ParentID,
		Created:  entity.Created,
		Updated:  entity.Updated,
	}

	data, exist := entity.Data[version]

	if !exist {
		return domain.Content{}, ContentError{entity.ID, version, "not found"}
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

		cp := domain.Property{
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
