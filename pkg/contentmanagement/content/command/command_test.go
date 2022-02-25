package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/siteconfiguration"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func Test_CreateContent(t *testing.T) {
	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	cdRepo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := cdRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test2",
	})
	assert.NoError(t, err)

	cmd := CreateContent{
		ContentDefinitionId: cid,
	}
	handler := CreateContentHandler{
		ContentDefinitionRepository: cdRepo,
		ContentRepository:           content.NewContentRepository(c),
	}

	contentId, err := handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, contentId)
}

func Test_CreateContent_Empty_ContentDefinition(t *testing.T) {
	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	c.Database("cms").Collection("content").Drop(context.Background())

	cmd := CreateContent{}
	handler := CreateContentHandler{
		ContentDefinitionRepository: contentdefinition.NewContentDefinitionRepository(c),
		ContentRepository:           content.NewContentRepository(c),
	}

	contentId, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Equal(t, uuid.UUID{}, contentId)
}

func Test_UpdateContent(t *testing.T) {
	c, err := db.Connect(context.Background(), "mongodb://0.0.0.0")
	assert.NoError(t, err)
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	c.Database("cms").Collection("content").Drop(context.Background())

	cdRepo := contentdefinition.NewContentDefinitionRepository(c)
	cdid, err := cdRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test2",
	})
	assert.NoError(t, err)

	contentRepo := content.NewContentRepository(c)
	contentId, err := contentRepo.CreateContent(context.Background(), content.Content{
		ContentDefinitionID: cdid,
	})
	assert.NoError(t, err)

	cmd := UpdateContent{
		Id: contentId,
		Fields: []struct {
			Language string
			Field    string
			Value    interface{}
		}{
			{
				Language: "sv-SE",
				Field:    content.NameField,
				Value:    "name test",
			},
			{
				Language: "sv-SE",
				Field:    content.UrlSegmentField,
				Value:    "url test",
			},
		},
	}

	cfg := siteconfiguration.SiteConfiguration{
		Languages: []language.Tag{
			language.MustParse("sv-SE"),
			language.MustParse("en-US"),
		},
	}
	handler := UpdateContentHandler{
		ContentDefinitionRepository: cdRepo,
		ContentRepository:           contentRepo,
		SiteConfiguration:           &cfg,
	}

	err = handler.Handle(context.Background(), cmd)
	assert.NoError(t, err)
	cont, err := contentRepo.GetContent(context.Background(), contentId)
	assert.NoError(t, err)
	assert.Equal(t, 1, cont.Version)
	assert.Equal(t, "name test", cont.Properties[cfg.Languages[0].String()][content.NameField])
	assert.Equal(t, "url-test", cont.Properties[cfg.Languages[0].String()][content.UrlSegmentField])
	assert.Equal(t, "name test", cont.Properties[cfg.Languages[1].String()][content.NameField])
	assert.Equal(t, "name-test", cont.Properties[cfg.Languages[1].String()][content.UrlSegmentField])
}
