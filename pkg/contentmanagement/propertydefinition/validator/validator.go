package validator

import (
	"context"
	"errors"
	"regexp"
)

/*
	propertyDefinition: {
		validators: [
			{
				"required": true
			},
			{
				"pattern": "^foo"
			},
		]
	}
*/
// Validator runs when content is Saved

type RequiredRule bool
type RegexRule string

type Validator interface {
	Validate(ctx context.Context, field interface{}) error
}

func Parse(name string, val interface{}) (Validator, error) {

	switch name {
	case "required":
		if b, ok := val.(bool); ok {
			return RequiredRule(b), nil
		}
		return nil, errors.New("parse error: cannot parse into type RequiredField")
	case "pattern":
		if str, ok := val.(string); ok {
			return RegexRule(str), nil
		}
	}

	return nil, errors.New("validator not found")
}

// Validators

func (r RequiredRule) Validate(ctx context.Context, field interface{}) error {
	if bool(r) && field == nil {
		return errors.New("required")
	}
	return nil
}

func (r RegexRule) Validate(ctx context.Context, field interface{}) error {

	str := string(r)
	b, ok := field.([]byte)

	if !ok {
		return errors.New("cannot validate unknown field")
	}

	match, err := regexp.Match(str, b)

	if err != nil {
		return err
	}

	if !match {
		return errors.New("pattern do not match")

	}
	return nil
}
