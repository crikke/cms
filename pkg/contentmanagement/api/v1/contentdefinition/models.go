package contentdefinition

type ContentDefinitionBody struct {
	// Content definition Name
	Name string
	// Content definition description
	Description string
}

type CreatePropertyDefinitionBody struct {
	Name        string
	Description string
	Type        string
}

type UpdatePropertyDefinitionBody struct {
	Name        string
	Description string
	Localized   bool
	Validation  map[string]interface{}
}
