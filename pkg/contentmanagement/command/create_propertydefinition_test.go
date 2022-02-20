package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentdelivery/db"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/contentmanagement/propertydefinition"
	"github.com/stretchr/testify/assert"
)

func Test_CreatePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")

	assert.NoError(t, err)

	contentRepo := contentdefinition.NewContentDefinitionRepository(c)

	cid, err := contentRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test",
	})

	assert.NoError(t, err)

	repo := propertydefinition.NewPropertyDefinitionRepository(c)
	handler := CreatePropertyDefinitionHandler{repo: repo}

	testpd := CreatePropertyDefinition{
		Name:                "Test",
		Description:         "Test",
		Type:                "Test",
		ContentDefinitionID: cid,
	}

	id, err := handler.Handle(context.TODO(), testpd)

	assert.NoError(t, err)

	cd, err := repo.GetPropertyDefinition(context.TODO(), id)
	assert.NoError(t, err)

	assert.Equal(t, testpd.Name, cd.Name)
	assert.Equal(t, testpd.Description, cd.Description)
	assert.Equal(t, testpd.Type, cd.Type)
}
