package query

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/command"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/contentmanagement/propertydefinition"
	"github.com/crikke/cms/pkg/contentmanagement/propertydefinition/validator"
	"github.com/crikke/cms/pkg/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_GetPropertyDefinitionValidationRule(t *testing.T) {

	tests := []struct {
		name      string
		typ       interface{}
		createcmd command.UpdateValidator
	}{
		{
			name:      "query RequiredRule",
			typ:       validator.RequiredRule(false),
			createcmd: command.UpdateValidator{ValidatorName: "required", Value: true},
		},
		{
			name:      "query RegexRule",
			typ:       validator.RegexRule(""),
			createcmd: command.UpdateValidator{ValidatorName: "pattern", Value: ""},
		},
	}
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cid, pid := createPropertyValidation(t, c)
			test.createcmd.ContentDefinitionID = cid
			test.createcmd.PropertyDefinitionID = pid

			repo := propertydefinition.NewPropertyDefinitionRepository(c)
			reqhandler := command.UpdateValidatorHandler{Repo: repo}
			reqhandler.Handle(context.Background(), test.createcmd)

			// query validation rule
			query := GetValidator{
				ContentDefinitionID:  cid,
				PropertyDefinitionID: pid,
				ValidatorName:        test.createcmd.ValidatorName,
			}
			queryHandler := GetValidatorHandler{
				repo: repo,
			}

			v, err := queryHandler.Handle(context.Background(), query)
			assert.NoError(t, err)

			assert.IsType(t, v, test.typ)
		})
	}

	cid, pid := createPropertyValidation(t, c)

	// create validation rule
	repo := propertydefinition.NewPropertyDefinitionRepository(c)
	reqcmd := command.UpdateValidator{ContentDefinitionID: cid, PropertyDefinitionID: pid, ValidatorName: "required", Value: true}
	reqhandler := command.UpdateValidatorHandler{Repo: repo}
	reqhandler.Handle(context.Background(), reqcmd)

	// query validation rule
	query := GetValidator{
		ContentDefinitionID:  cid,
		PropertyDefinitionID: pid,
		ValidatorName:        "required",
	}
	queryHandler := GetValidatorHandler{
		repo: repo,
	}

	v, err := queryHandler.Handle(context.Background(), query)
	assert.NoError(t, err)

	required, ok := v.(validator.RequiredRule)

	assert.True(t, ok)
	assert.True(t, bool(required))
}

func createPropertyValidation(t *testing.T, c *mongo.Client) (cid, pid uuid.UUID) {
	err := c.Database("cms").Collection("contentdefinition").Drop(context.Background())
	assert.NoError(t, err)

	// create contentdefinition
	contentRepo := contentdefinition.NewContentDefinitionRepository(c)
	cid, err = contentRepo.CreateContentDefinition(context.Background(), &contentdefinition.ContentDefinition{
		Name: "test",
	})

	assert.NoError(t, err)

	// create propertydefinition
	repo := propertydefinition.NewPropertyDefinitionRepository(c)
	handler := command.CreatePropertyDefinitionHandler{Repo: repo}

	testpd := command.CreatePropertyDefinition{
		Name:                "pd1",
		Description:         "pd2",
		Type:                "text",
		ContentDefinitionID: cid,
	}
	pid, err = handler.Handle(context.TODO(), testpd)
	assert.NoError(t, err)

	return
}
