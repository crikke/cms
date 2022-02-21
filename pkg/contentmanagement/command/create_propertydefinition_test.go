package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentdelivery/db"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/stretchr/testify/assert"
)

func Test_CreatePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	assert.NoError(t, err)

	contentRepo := contentdefinition.NewContentDefinitionRepository(c)

	cid, err := contentRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test",
	})

	assert.NoError(t, err)

	repo := contentdefinition.NewPropertyDefinitionRepository(c)
	handler := CreatePropertyDefinitionHandler{repo: repo}

	testpd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
	}

	id, err := handler.Handle(context.TODO(), testpd)

	assert.NoError(t, err)

	actual, err := repo.GetPropertyDefinition(context.TODO(), cid, id)
	assert.NoError(t, err)

	assert.Equal(t, testpd.Name, actual.Name)
	assert.Equal(t, testpd.Description, actual.Description)
	assert.Equal(t, testpd.Type, actual.Type)
}
