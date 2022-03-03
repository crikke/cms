package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/db"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_CreatePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	assert.NoError(t, err)

	repo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := repo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name:                "test",
		Propertydefinitions: make(map[string]contentdefinition.PropertyDefinition),
	})

	assert.NoError(t, err)

	handler := CreatePropertyDefinitionHandler{Repo: repo}

	testpd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
	}

	id, err := handler.Handle(context.TODO(), testpd)

	assert.NoError(t, err)

	cd, err := repo.GetContentDefinition(context.TODO(), cid)
	assert.NoError(t, err)

	actual := cd.Propertydefinitions[testpd.Name]
	assert.Equal(t, testpd.Description, actual.Description)
	assert.Equal(t, testpd.Type, actual.Type)
	assert.Equal(t, id, actual.ID)
}

func Test_UpdatePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	assert.NoError(t, err)

	// create contentdefinition
	repo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := repo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name:                "test",
		Propertydefinitions: make(map[string]contentdefinition.PropertyDefinition),
	})
	assert.NoError(t, err)

	// create propertydefinition
	createhandler := CreatePropertyDefinitionHandler{Repo: repo}
	createcmd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
	}

	pid, err := createhandler.Handle(context.TODO(), createcmd)
	assert.NoError(t, err)

	// update propertydefiniton
	updatehandler := UpdatePropertyDefinitionHandler{repo: repo}

	str := "updated"
	b := true
	updatecmd := UpdatePropertyDefinition{
		ContentDefinitionID:  cid,
		PropertyDefinitionID: pid,
		Name:                 &str,
		Description:          &str,
		Localized:            &b,
	}

	err = updatehandler.Handle(context.Background(), updatecmd)
	assert.NoError(t, err)

	// get propertydefiniton
	cd, err := repo.GetContentDefinition(context.TODO(), cid)
	assert.NoError(t, err)

	actual := cd.Propertydefinitions[str]
	assert.NoError(t, err)
	// assert.Equal(t, str, actual.Name)
	assert.Equal(t, str, actual.Description)
	assert.Equal(t, createcmd.Type, actual.Type)
	assert.Equal(t, b, actual.Localized)
}

func Test_DeletePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	assert.NoError(t, err)

	// create contentdefinition
	repo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := repo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name:                "test",
		Propertydefinitions: make(map[string]contentdefinition.PropertyDefinition),
	})
	assert.NoError(t, err)

	// create propertydefinition
	createhandler := CreatePropertyDefinitionHandler{Repo: repo}
	createcmd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
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
	}

	err = deletehandler.Handle(context.Background(), deletecmd)
	assert.NoError(t, err)

	// get propertydefiniton
	_, err = repo.GetPropertyDefinition(context.TODO(), cid, pid)
	assert.Error(t, err)
	assert.Equal(t, err, mongo.ErrNoDocuments)
}

func Test_AddValidation(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	assert.NoError(t, err)

	repo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := repo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name:                "test",
		Propertydefinitions: make(map[string]contentdefinition.PropertyDefinition),
	})

	assert.NoError(t, err)

	handler := CreatePropertyDefinitionHandler{Repo: repo}

	testpd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
	}
	pid, err := handler.Handle(context.TODO(), testpd)
	assert.NoError(t, err)

	cmd1 := UpdateValidator{ContentDefinitionID: cid, PropertyDefinitionID: pid, ValidatorName: "required", Value: true}
	cmd2 := UpdateValidator{ContentDefinitionID: cid, PropertyDefinitionID: pid, ValidatorName: "pattern", Value: "^foo"}
	validationhandler := UpdateValidatorHandler{Repo: repo}

	err = validationhandler.Handle(context.Background(), cmd1)
	assert.NoError(t, err)
	err = validationhandler.Handle(context.Background(), cmd2)
	assert.NoError(t, err)

	cd, err := repo.GetContentDefinition(context.Background(), cid)

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
}
