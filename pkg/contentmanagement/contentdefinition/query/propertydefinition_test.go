package query

import (
	"context"
	"testing"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition/command"
	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition/validator"
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

			repo := contentdefinition.NewContentDefinitionRepository(c)
			reqhandler := command.UpdateValidatorHandler{Repo: repo}
			reqhandler.Handle(context.Background(), test.createcmd)

			// query validation rule
			query := GetValidatorForProperty{
				ContentDefinitionID:  cid,
				PropertyDefinitionID: pid,
				ValidatorName:        test.createcmd.ValidatorName,
			}
			queryHandler := GetValidatorForPropertyHandler{
				Repo: repo,
			}

			v, err := queryHandler.Handle(context.Background(), query)
			assert.NoError(t, err)

			assert.IsType(t, v, test.typ)
		})
	}

	cid, pid := createPropertyValidation(t, c)

	// create validation rule
	repo := contentdefinition.NewContentDefinitionRepository(c)
	reqcmd := command.UpdateValidator{ContentDefinitionID: cid, PropertyDefinitionID: pid, ValidatorName: "required", Value: true}
	reqhandler := command.UpdateValidatorHandler{Repo: repo}
	reqhandler.Handle(context.Background(), reqcmd)

	// query validation rule
	query := GetValidatorForProperty{
		ContentDefinitionID:  cid,
		PropertyDefinitionID: pid,
		ValidatorName:        "required",
	}
	queryHandler := GetValidatorForPropertyHandler{
		Repo: repo,
	}

	v, err := queryHandler.Handle(context.Background(), query)
	assert.NoError(t, err)

	required, ok := v.(validator.RequiredRule)

	assert.True(t, ok)
	assert.True(t, bool(required))
}

func Test_GetAllPropertyDefinitionValidationRules(t *testing.T) {

	tests := []struct {
		name      string
		createcmd []command.UpdateValidator
	}{
		{
			name: "query RequiredRule",
			createcmd: []command.UpdateValidator{
				{ValidatorName: "required", Value: true},
				{ValidatorName: "pattern", Value: "true"},
			},
		},
	}
	c, err := db.Connect(context.TODO(), "mongodb://0.0.0.0")
	assert.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cid, pid := createPropertyValidation(t, c)

			repo := contentdefinition.NewContentDefinitionRepository(c)
			for _, cmd := range test.createcmd {
				reqhandler := command.UpdateValidatorHandler{Repo: repo}
				cmd.PropertyDefinitionID = pid
				cmd.ContentDefinitionID = cid

				err := reqhandler.Handle(context.Background(), cmd)
				assert.NoError(t, err)
			}

			// query validation rule
			query := GetAllValidatorsForProperty{
				ContentDefinitionID:  cid,
				PropertyDefinitionID: pid,
			}
			queryHandler := GetAllValidatorsForPropertyHandler{
				Repo: repo,
			}

			v, err := queryHandler.Handle(context.Background(), query)
			assert.NoError(t, err)
			assert.Len(t, v, len(test.createcmd))
		})
	}

	cid, pid := createPropertyValidation(t, c)

	// create validation rule
	repo := contentdefinition.NewContentDefinitionRepository(c)
	reqcmd := command.UpdateValidator{ContentDefinitionID: cid, PropertyDefinitionID: pid, ValidatorName: "required", Value: true}
	reqhandler := command.UpdateValidatorHandler{Repo: repo}
	reqhandler.Handle(context.Background(), reqcmd)

	// query validation rule
	query := GetValidatorForProperty{
		ContentDefinitionID:  cid,
		PropertyDefinitionID: pid,
		ValidatorName:        "required",
	}
	queryHandler := GetValidatorForPropertyHandler{
		Repo: repo,
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
	repo := contentdefinition.NewContentDefinitionRepository(c)
	handler := command.CreatePropertyDefinitionHandler{Repo: repo}

	testpd := command.CreatePropertyDefinition{
		Name:                "pd111",
		Description:         "pd222",
		Type:                "text",
		ContentDefinitionID: cid,
	}
	pid, err = handler.Handle(context.TODO(), testpd)
	assert.NoError(t, err)

	return
}
