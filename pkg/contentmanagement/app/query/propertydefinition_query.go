package query

import (
	"context"

	"github.com/crikke/cms/pkg/contentmanagement/contentdefinition"
	"github.com/google/uuid"
)

type GetContentDefinition struct {
	ID uuid.UUID
}

type GetContentDefinitionHandler struct {
	Repo contentdefinition.ContentDefinitionRepository
}

func (h GetContentDefinitionHandler) Handle(ctx context.Context, query GetContentDefinition) (contentdefinition.ContentDefinition, error) {

	return h.Repo.GetContentDefinition(ctx, query.ID)
}

// type GetValidatorForProperty struct {
// 	ContentDefinitionID  uuid.UUID
// 	PropertyDefinitionID uuid.UUID
// 	ValidatorName        string
// }

// type GetValidatorForPropertyHandler struct {
// 	Repo contentdefinition.ContentDefinitionRepository
// }

// func (h GetValidatorForPropertyHandler) Handle(ctx context.Context, query GetValidatorForProperty) (validator.Validator, error) {

// 	pd, err := h.Repo.GetPropertyDefinition(ctx, query.ContentDefinitionID, query.PropertyDefinitionID)

// 	if err != nil {
// 		return nil, err
// 	}

// 	v, ok := pd.Validators[query.ValidatorName]

// 	if !ok {
// 		return nil, errors.New("validator not found")
// 	}

// 	val, err := validator.Parse(query.ValidatorName, v)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return val, nil
// }

// type GetAllValidatorsForProperty struct {
// 	ContentDefinitionID  uuid.UUID
// 	PropertyDefinitionID uuid.UUID
// }

// type GetAllValidatorsForPropertyHandler struct {
// 	Repo contentdefinition.ContentDefinitionRepository
// }

// func (h GetAllValidatorsForPropertyHandler) Handle(ctx context.Context, query GetAllValidatorsForProperty) ([]validator.Validator, error) {

// 	pd, err := h.Repo.GetPropertyDefinition(ctx, query.ContentDefinitionID, query.PropertyDefinitionID)

// 	if err != nil {
// 		return nil, err
// 	}

// 	result := []validator.Validator{}

// 	for name, v := range pd.Validators {
// 		val, err := validator.Parse(name, v)

// 		if err != nil {
// 			return nil, err
// 		}

// 		result = append(result, val)
// 	}
// 	return result, nil
// }
