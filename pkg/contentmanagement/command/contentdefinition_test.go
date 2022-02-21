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

func Test_UpdateContentDefinition(t *testing.T) {

	tests := []struct {
		name      string
		existing  contentdefinition.ContentDefinition
		updatecmd UpdateContentDefinition
		expect    contentdefinition.ContentDefinition
	}{
		{
			name: "update all fields",
			existing: contentdefinition.ContentDefinition{
				Name: "old",
			},
			updatecmd: UpdateContentDefinition{
				Name:        "updated",
				Description: "updated",
			},
			expect: contentdefinition.ContentDefinition{
				Name:        "updated",
				Description: "updated",
			},
		},
		{
			name: "update single field",
			existing: contentdefinition.ContentDefinition{
				Name:        "old",
				Description: "old",
			},
			updatecmd: UpdateContentDefinition{
				Name: "updated",
			},
			expect: contentdefinition.ContentDefinition{
				Name:        "updated",
				Description: "old",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
			assert.NoError(t, err)

			repo := contentdefinition.NewContentDefinitionRepository(client)

			id, err := repo.CreateContentDefinition(context.Background(), &test.existing)
			assert.NoError(t, err)

			cmd := test.updatecmd
			cmd.ID = id

			handler := UpdateContentDefinitionHandler{repo: repo}
			handler.Handle(context.TODO(), cmd)

			updated, err := repo.GetContentDefinition(context.TODO(), id)
			assert.NoError(t, err)

			assert.Equal(t, test.expect.Name, updated.Name)
			assert.Equal(t, test.expect.Description, updated.Description)
		})
	}
}
