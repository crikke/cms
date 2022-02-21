package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentdelivery/db"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
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

func Test_UpdatePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	assert.NoError(t, err)

	// create contentdefinition
	contentRepo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := contentRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test",
	})
	assert.NoError(t, err)
	propRepo := contentdefinition.NewPropertyDefinitionRepository(c)

	// create propertydefinition
	createhandler := CreatePropertyDefinitionHandler{repo: propRepo}
	createcmd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
	}

	pid, err := createhandler.Handle(context.TODO(), createcmd)
	assert.NoError(t, err)

	// update propertydefiniton
	updatehandler := UpdatePropertyDefinitionHandler{repo: propRepo}

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
	actual, err := propRepo.GetPropertyDefinition(context.TODO(), cid, pid)
	assert.NoError(t, err)
	assert.Equal(t, str, actual.Name)
	assert.Equal(t, str, actual.Description)
	assert.Equal(t, createcmd.Type, actual.Type)
	assert.Equal(t, b, actual.Localized)
}

func Test_DeletePropertyDefinition(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	assert.NoError(t, err)

	// create contentdefinition
	contentRepo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err := contentRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test",
	})
	assert.NoError(t, err)
	propRepo := contentdefinition.NewPropertyDefinitionRepository(c)

	// create propertydefinition
	createhandler := CreatePropertyDefinitionHandler{repo: propRepo}
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
		repo: propRepo,
	}
	deletecmd := DeletePropertyDefinition{
		ContentDefinitionID:  cid,
		PropertyDefinitionID: pid,
	}

	err = deletehandler.Handle(context.Background(), deletecmd)
	assert.NoError(t, err)

	// get propertydefiniton
	_, err = propRepo.GetPropertyDefinition(context.TODO(), cid, pid)
	assert.Error(t, err)
	assert.Equal(t, err, mongo.ErrNoDocuments)
}
