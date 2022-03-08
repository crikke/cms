package content

import "github.com/go-openapi/strfmt"

type CreateContentRequest struct {
	ContentDefinitionId strfmt.UUID
}

type UpdateContentRequestBody struct {
	// Version
	Version int
	// Language
	Language string
	// Properties
	Fields map[string]interface{}
}

type OKResult struct {
}
