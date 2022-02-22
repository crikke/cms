package command

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/contentmanagement/propertydefinition"
	"github.com/crikke/cms/pkg/db"
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

	repo := propertydefinition.NewPropertyDefinitionRepository(c)
	handler := CreatePropertyDefinitionHandler{repo: repo}

	testpd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
	}

	_, err = handler.Handle(context.TODO(), testpd)

	assert.NoError(t, err)

	// actual, err := repo.GetPropertyDefinition(context.TODO(), cid, id)
	assert.NoError(t, err)

	// assert.Equal(t, testpd.Name, actual.GetName())
	// assert.Equal(t, testpd.Description, actual.GetDescription())
	// assert.Equal(t, testpd.Type, actual.GetType())
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
	propRepo := propertydefinition.NewPropertyDefinitionRepository(c)

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
	// actual, err := propRepo.GetPropertyDefinition(context.TODO(), cid, pid)
	assert.NoError(t, err)
	// assert.Equal(t, str, actual.GetName())
	// assert.Equal(t, str, actual.GetDescription())
	// // assert.Equal(t, createcmd.Type, actual.Type)
	// assert.Equal(t, b, actual.GetLocalized())
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
	propRepo := propertydefinition.NewPropertyDefinitionRepository(c)

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

func Test_AddValidation(t *testing.T) {
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	assert.NoError(t, err)

	contentRepo := contentdefinition.NewContentDefinitionRepository(c)

	cid, err := contentRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test",
	})

	assert.NoError(t, err)

	repo := propertydefinition.NewPropertyDefinitionRepository(c)
	handler := CreatePropertyDefinitionHandler{repo: repo}

	testpd := CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
	}
	pid, err := handler.Handle(context.TODO(), testpd)
	assert.NoError(t, err)

	reqcmd := AddValidator{ContentDefinitionID: cid, PropertyDefinitionID: pid, ValidatorName: "required", Value: true}
	reqcmd2 := AddValidator{ContentDefinitionID: cid, PropertyDefinitionID: pid, ValidatorName: "pattern", Value: "^foo"}
	reqhandler := AddValidatorHandler{repo: repo}

	err = reqhandler.Handle(context.Background(), reqcmd)
	assert.NoError(t, err)
	err = reqhandler.Handle(context.Background(), reqcmd2)
	assert.NoError(t, err)

	pd, err := repo.GetPropertyDefinition(context.Background(), cid, pid)
	assert.NoError(t, err)

	req, ok := pd.Validators["required"]
	assert.True(t, ok)
	assert.True(t, req.(bool))

	pattern, ok := pd.Validators["pattern"]
	assert.True(t, ok)
	assert.Equal(t, "^foo", pattern)
}
