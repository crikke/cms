//go:build integration

package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentdefinition"
	"github.com/crikke/cms/pkg/db"
	"github.com/crikke/cms/pkg/workspace"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_CreatePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	wsRepo := workspace.NewWorkspaceRepository(c)
	ws, err := wsRepo.Create(context.Background(), workspace.Workspace{
		Name: "test",
	})

	repo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := repo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name:                "test",
		Propertydefinitions: make(map[string]contentdefinition.PropertyDefinition),
	}, ws)

	assert.NoError(t, err)

	handler := CreatePropertyDefinitionHandler{Repo: repo}

	testpd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
		WorkspaceID:         ws,
	}

	id, err := handler.Handle(context.TODO(), testpd)
	assert.NoError(t, err)

	cd, err := repo.GetContentDefinition(context.TODO(), cid, ws)
	assert.NoError(t, err)

	actual := cd.Propertydefinitions[testpd.Name]
	assert.Equal(t, testpd.Description, actual.Description)
	assert.Equal(t, testpd.Type, actual.Type)
	assert.Equal(t, id, actual.ID)

	t.Cleanup(func() {
		workspaces, _ := wsRepo.ListAll(context.Background())

		for _, ws := range workspaces {
			wsRepo.Delete(context.Background(), ws.ID)
		}
	})
}

func Test_UpdatePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")

	assert.NoError(t, err)
	wsRepo := workspace.NewWorkspaceRepository(c)
	ws, err := wsRepo.Create(context.Background(), workspace.Workspace{
		Name: "test",
	})
	// create contentdefinition
	repo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := repo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name:                "test",
		Propertydefinitions: make(map[string]contentdefinition.PropertyDefinition),
	}, ws)
	assert.NoError(t, err)

	// create propertydefinition
	createhandler := CreatePropertyDefinitionHandler{Repo: repo}
	createcmd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
		WorkspaceID:         ws,
	}

	pid, err := createhandler.Handle(context.TODO(), createcmd)
	assert.NoError(t, err)

	// update propertydefiniton
	updatehandler := UpdatePropertyDefinitionHandler{Repo: repo}

	str := "updated"
	b := true
	updatecmd := UpdatePropertyDefinition{
		ContentDefinitionID:  cid,
		PropertyDefinitionID: pid,
		Name:                 &str,
		Description:          &str,
		Localized:            &b,
		WorkspaceID:          ws,
	}

	err = updatehandler.Handle(context.Background(), updatecmd)
	assert.NoError(t, err)

	// get propertydefiniton
	cd, err := repo.GetContentDefinition(context.TODO(), cid, ws)
	assert.NoError(t, err)

	actual := cd.Propertydefinitions[str]
	assert.NoError(t, err)
	// assert.Equal(t, str, actual.Name)
	assert.Equal(t, str, actual.Description)
	assert.Equal(t, createcmd.Type, actual.Type)
	assert.Equal(t, b, actual.Localized)

	t.Cleanup(func() {
		workspaces, _ := wsRepo.ListAll(context.Background())

		for _, ws := range workspaces {
			wsRepo.Delete(context.Background(), ws.ID)
		}
	})
}

func Test_DeletePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	assert.NoError(t, err)
	wsRepo := workspace.NewWorkspaceRepository(c)
	ws, err := wsRepo.Create(context.Background(), workspace.Workspace{
		Name: "test",
	})
	// create contentdefinition
	repo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := repo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name:                "test",
		Propertydefinitions: make(map[string]contentdefinition.PropertyDefinition),
	}, ws)
	assert.NoError(t, err)

	// create propertydefinition
	createhandler := CreatePropertyDefinitionHandler{Repo: repo}
	createcmd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
		WorkspaceID:         ws,
	}

	pid, err := createhandler.Handle(context.TODO(), createcmd)
	assert.NoError(t, err)

	// delete propertydefinition
	deletehandler := DeletePropertyDefinitionHandler{
		repo: repo,
	}
	deletecmd := DeletePropertyDefinition{
		ContentDefinitionID:  cid,
		PropertyDefinitionID: pid,
		WorkspaceID:          ws,
	}

	err = deletehandler.Handle(context.Background(), deletecmd)
	assert.NoError(t, err)

	// get propertydefiniton
	_, err = repo.GetPropertyDefinition(context.TODO(), cid, pid, ws)
	assert.Error(t, err)
	assert.Equal(t, err, mongo.ErrNoDocuments)
	t.Cleanup(func() {
		workspaces, _ := wsRepo.ListAll(context.Background())

		for _, ws := range workspaces {
			wsRepo.Delete(context.Background(), ws.ID)
		}
	})
}

func Test_AddValidation(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")

	wsRepo := workspace.NewWorkspaceRepository(c)
	ws, err := wsRepo.Create(context.Background(), workspace.Workspace{
		Name: "test",
	})

	assert.NoError(t, err)

	repo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := repo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name:                "test",
		Propertydefinitions: make(map[string]contentdefinition.PropertyDefinition),
	}, ws)

	assert.NoError(t, err)

	handler := CreatePropertyDefinitionHandler{Repo: repo}

	testpd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
		WorkspaceID:         ws,
	}
	pid, err := handler.Handle(context.TODO(), testpd)
	assert.NoError(t, err)

	cmd1 := UpdateValidator{ContentDefinitionID: cid, PropertyDefinitionID: pid, ValidatorName: "required", Value: true, WorkspaceID: ws}
	cmd2 := UpdateValidator{ContentDefinitionID: cid, PropertyDefinitionID: pid, ValidatorName: "pattern", Value: "^foo", WorkspaceID: ws}
	validationhandler := UpdateValidatorHandler{Repo: repo}

	err = validationhandler.Handle(context.Background(), cmd1)
	assert.NoError(t, err)
	err = validationhandler.Handle(context.Background(), cmd2)
	assert.NoError(t, err)

	cd, err := repo.GetContentDefinition(context.Background(), cid, ws)

	pd := contentdefinition.PropertyDefinition{}
	for _, p := range cd.Propertydefinitions {
		if p.ID == pid {
			pd = p
			break
		}
	}

	assert.NoError(t, err)

	req, ok := pd.Validators["required"]
	assert.True(t, ok)
	assert.True(t, req.(bool))

	pattern, ok := pd.Validators["pattern"]
	assert.True(t, ok)
	assert.Equal(t, "^foo", pattern)

	t.Cleanup(func() {
		workspaces, _ := wsRepo.ListAll(context.Background())

		for _, ws := range workspaces {
			wsRepo.Delete(context.Background(), ws.ID)
		}
	})
}
