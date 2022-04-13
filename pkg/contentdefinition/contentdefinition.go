package contentdefinition

import (
	"errors"
	"time"

	"github.com/crikke/cms/pkg/contentdefinition/validator"
	"github.com/google/uuid"
)

const (
	PROPFIELD_NAME = "name"

	PropertyTypeText   = "text"
	PropertyTypeNumber = "number"
	PropertyTypeBool   = "bool"

	ErrPropertyAlreadyExists = "propertydefinition already exists on contentdefinition"
	ErrPropertyTypeNotExists = "propertydefinition type does not exist"
)

// swagger:model ContentDefinition
type ContentDefinition struct {
	ID                  uuid.UUID `bson:"_id"`
	Name                string    `bson:"name,omitempty"`
	Description         string    `bson:"description,omitempty"`
	Created             time.Time
	Propertydefinitions map[string]PropertyDefinition
}

// swagger:model PropertyDefinition
type PropertyDefinition struct {
	ID uuid.UUID `bson:"id"`
	// Name        string    `bson:"name,omitempty"`
	Description string `bson:"description,omitempty"`
	Type        string `bson:"type,omitempty"`
	Localized   bool   `bson:"localized,omitempty"`
	// instead of using map[strin]validator.Validator, interface{} is used
	// this wont be a problem becuase they will be translated to validator.Validator in GetValidatorQueury
	Validators map[string]interface{} `bson:"validators,omitempty"`
}

type ContentDefinitionFactory struct {
}

func (f ContentDefinitionFactory) NewContentDefinition(name, desc string) (ContentDefinition, error) {
	if name == "" {
		return ContentDefinition{}, errors.New("name required")
	}

	return ContentDefinition{
		Name:        name,
		Description: desc,
		Propertydefinitions: map[string]PropertyDefinition{
			PROPFIELD_NAME: {
				ID:        uuid.New(),
				Type:      "text",
				Localized: true,
				Validators: map[string]interface{}{
					"required": validator.Required(true),
				},
			},
		}}, nil
}

func (f ContentDefinitionFactory) NewPropertyDefinition(cd *ContentDefinition, name, propertyType, description string, localized bool) error {

	if _, exist := cd.Propertydefinitions[name]; exist {
		return errors.New(ErrPropertyAlreadyExists)
	}

	if cd.Propertydefinitions == nil {
		cd.Propertydefinitions = make(map[string]PropertyDefinition)
	}

	pd := PropertyDefinition{
		ID:          uuid.New(),
		Description: description,
		Type:        propertyType,
		Localized:   localized,
		Validators: map[string]interface{}{
			validator.RuleRequired: validator.Required(false),
		},
	}

	switch propertyType {
	case PropertyTypeText:
		pd.Validators[validator.RuleRegex] = validator.Regex("")
		pd.Validators[validator.RuleRange] = validator.Range{}
	case PropertyTypeBool:
		break
	case PropertyTypeNumber:
		pd.Validators[validator.RuleRange] = validator.Range{}
	default:
		return errors.New(ErrPropertyTypeNotExists)
	}

	cd.Propertydefinitions[name] = pd

	return nil
}

func (f ContentDefinitionFactory) UpdatePropertyDefinitionName(cd *ContentDefinition, id uuid.UUID, name string) error {

	if name == "" {
		return errors.New("name required")
	}

	if _, exists := cd.Propertydefinitions[name]; exists {
		return errors.New("property with name already exists")
	}

	pd := PropertyDefinition{}
	pdName := ""
	for n, p := range cd.Propertydefinitions {
		if p.ID == id {
			pd = p
			pdName = n
			break
		}
	}
	if pdName == "" {
		return errors.New("property not found")
	}

	delete(cd.Propertydefinitions, pdName)
	cd.Propertydefinitions[name] = pd

	return nil
}

// This is essentialy an HTTP PUT, every propertydefinition that must be included in the argument, otherwise it is assumed to be deleted.
func (f ContentDefinitionFactory) UpdatePropertyDefinitions(cd *ContentDefinition, propertyDefinitions map[string]PropertyDefinition) error {

	// All propertydefinitions that havent been updated will be deleted.
	updatedProps := make(map[uuid.UUID]PropertyDefinition, 0)

	for _, prop := range propertyDefinitions {
		err := f.UpdatePropertyDefinition(
			cd,
			prop.ID,
			prop.Description,
			prop.Localized,
			prop.Validators)

		updatedProps[prop.ID] = prop

		if err != nil {
			return err
		}
	}

	// check for props that should be deleted
	for name, prop := range cd.Propertydefinitions {

		if _, ok := updatedProps[prop.ID]; ok {
			continue
		}

		delete(cd.Propertydefinitions, name)

	}

	// if property exist, remove it from map and store it in collidingNames
	// next loop, set each property again, if still collides there is an actual collsion and return error
	// collidingNames := make(map[uuid.UUID]string)

	// update property names
	for name, prop := range propertyDefinitions {

		existing, ok := cd.Propertydefinitions[name]

		if ok {

			// if property with name exist and the IDs isnt equal, add property to map.
			// when all names have been updated
			// try to resolve colliding names
			if existing.ID != prop.ID {
				return errors.New(ErrPropertyAlreadyExists)
			}

			continue
		}

		if err := f.UpdatePropertyDefinitionName(cd, prop.ID, name); err != nil {
			return err
		}
	}

	return nil
}

func (f ContentDefinitionFactory) UpdatePropertyDefinition(cd *ContentDefinition, id uuid.UUID, desc string, localized bool, validationRules map[string]interface{}) error {
	pd := PropertyDefinition{}
	pdName := ""
	for n, p := range cd.Propertydefinitions {
		if p.ID == id {
			pd = p
			pdName = n

			break
		}
	}

	if pdName == "" {
		return errors.New("property not found")
	}

	pd.Localized = localized
	pd.Description = desc

	for k, v := range validationRules {
		_, ok := pd.Validators[k]

		if !ok {
			return errors.New("validator not found")
		}

		pd.Validators[k] = v
	}

	cd.Propertydefinitions[pdName] = pd
	return nil
}

func NewContentDefinition(name, desc string) (ContentDefinition, error) {

	if name == "" {
		return ContentDefinition{}, errors.New("name required")
	}

	return ContentDefinition{
		Name:        name,
		Description: desc,
		Propertydefinitions: map[string]PropertyDefinition{
			PROPFIELD_NAME: {
				ID:        uuid.New(),
				Type:      "text",
				Localized: true,
				Validators: map[string]interface{}{
					"required": validator.Required(true),
				},
			},
		}}, nil
}

func (cd ContentDefinition) PropertyValid(field, lang string, value interface{}) error {

	pd, ok := cd.Propertydefinitions[field]

	if !ok {
		return errors.New("property does not exist")
	}

	if !pd.Localized && lang != "" {
		return errors.New("content.ErrUnlocalizedPropLocalizedValue")
	}
	return nil
}
