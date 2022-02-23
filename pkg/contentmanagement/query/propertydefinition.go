package query

import (
	"context"
	"errors"

	"github.com/crikke/cms/pkg/contentmanagement/propertydefinition"
	"github.com/crikke/cms/pkg/contentmanagement/propertydefinition/validator"
	"github.com/google/uuid"
)

type GetValidator struct {
	ContentDefinitionID  uuid.UUID
	PropertyDefinitionID uuid.UUID
	ValidatorName        string
}

type GetValidatorHandler struct {
	repo propertydefinition.PropertyDefinitionRepository
}

func (h GetValidatorHandler) Handle(ctx context.Context, query GetValidator) (validator.Validator, error) {

	pd, err := h.repo.GetPropertyDefinition(ctx, query.ContentDefinitionID, query.PropertyDefinitionID)

	if err != nil {
		return nil, err
	}

	v, ok := pd.Validators[query.ValidatorName]

	if !ok {
		return nil, errors.New("validator not found")
	}

	val, err := validator.Parse(query.ValidatorName, v)

	if err != nil {
		return nil, err
	}
	return val, nil
}
