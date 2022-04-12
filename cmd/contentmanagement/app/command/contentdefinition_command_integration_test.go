//go:build integration

package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/workspace"
	"github.com/stretchr/testify/assert"
)

func Test_CreateContentDefinition(t *testing.T) {

	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")

	assert.NoError(t, err)

	wsRepo := workspace.NewWorkspaceRepository(c)
	ws, err := wsRepo.Create(context.Background(), workspace.Workspace{
		Name: "test",
	})

	repo := contentdefinition.NewContentDefinitionRepository(c)
	handler := CreateContentDefinitionHandler{Repo: repo}

	testcd := CreateContentDefinition{
		Name:        "Test",
		Description: "Test",
		WorkspaceId: ws,
	}

	id, err := handler.Handle(context.TODO(), testcd)

	assert.NoError(t, err)

	cd, err := repo.GetContentDefinition(context.TODO(), id, ws)
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

	client, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	wsRepo := workspace.NewWorkspaceRepository(client)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NoError(t, err)

			ws, err := wsRepo.Create(context.Background(), workspace.Workspace{
				Name: "test",
			})

			assert.NoError(t, err)

			repo := contentdefinition.NewContentDefinitionRepository(client)

			id, err := repo.CreateContentDefinition(context.Background(), &test.existing, ws)
			assert.NoError(t, err)

			cmd := test.updatecmd
			cmd.ContentDefinitionID = id
			cmd.WorkspaceId = ws

			handler := UpdateContentDefinitionHandler{Repo: repo}
			handler.Handle(context.TODO(), cmd)

			updated, err := repo.GetContentDefinition(context.TODO(), id, ws)
			assert.NoError(t, err)

			assert.Equal(t, test.expect.Name, updated.Name)
			assert.Equal(t, test.expect.Description, updated.Description)
		})
	}

	t.Cleanup(func() {
		workspaces, _ := wsRepo.ListAll(context.Background())

		for _, ws := range workspaces {
			wsRepo.Delete(context.Background(), ws.ID)
		}
	})
}
