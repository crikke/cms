package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/content"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

	cmd := CreateContent{}
	handler := CreateContentHandler{
		ContentDefinitionRepository: contentdefinition.NewContentDefinitionRepository(c),
		ContentRepository:           content.NewContentRepository(c),
	}

	contentId, err := handler.Handle(context.Background(), cmd)
	assert.Error(t, err)
	assert.Equal(t, uuid.UUID{}, contentId)
}
