package contentdefinition

import "github.com/crikke/cms/pkg/contentdefinition"

type ContentDefinitionBody struct {
	// Content definition Name
	Name string
	// Content definition description
	Description string
	// PropertyDefinitions

	PropertyDefinitions map[string]contentdefinition.PropertyDefinition
}

//! TODO Remove this
type CreatePropertyDefinitionBody struct {
	Name        string
	Description string
	Type        string
}

type PropertyDefinitionBody struct {
	Type        string
	Description string
	Localized   bool
	Validation  map[string]interface{}
}
