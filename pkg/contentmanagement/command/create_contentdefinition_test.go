package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentdelivery/db"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/stretchr/testify/assert"
)

func Test_CreateContentDefinition(t *testing.T) {

	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")

	assert.NoError(t, err)

	repo := contentdefinition.NewContentDefinitionRepository(c)
	handler := CreateContentDefinitionHandler{repo: repo}

	testcd := CreateContentDefinition{
		Name:        "Test",
		Description: "Test",
	}

	id, err := handler.Handle(context.TODO(), testcd)

	assert.NoError(t, err)

	cd, err := repo.GetContentDefinition(context.TODO(), id)
	assert.NoError(t, err)

	assert.Equal(t, testcd.Name, cd.Name)
	assert.Equal(t, testcd.Description, cd.Description)
}
