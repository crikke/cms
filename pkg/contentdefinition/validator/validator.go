package validator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

const (
	RuleRequired = "required"
	RuleRegex    = "regex"
	RuleRange    = "range"
)

type Required bool
type Regex string
type Range struct {
	Min *float64 `bson:"min, omitempty"`
	Max *float64 `bson:"max, omitempty"`
}

type Validator interface {
	Validate(ctx context.Context, field interface{}) error
}

func Parse(name string, val interface{}) (Validator, error) {

	switch name {
	case "required":
		if b, ok := val.(bool); ok {
			return Required(b), nil
		}
		return nil, errors.New("parse error: cannot parse into type RequiredField")
	case "pattern":
		if str, ok := val.(string); ok {
			return Regex(str), nil
		}
		return nil, errors.New("pattern is not of type string")
	case "range":
		if r, ok := val.(Range); ok {
			return r, nil
		}
	}

	return nil, errors.New("validator not found")
}

// Validators

// 0 is a valid number so wont validate
func (r Required) Validate(ctx context.Context, field interface{}) error {

	if !bool(r) {
		return nil
	}
	// check for nil
	if field == nil {
		return errors.New("required")
	}

	// check for empty string

	if str, ok := field.(string); ok && str == "" {
		return errors.New("required")
	}

	return nil
}

func (r Regex) Validate(ctx context.Context, field interface{}) error {

	if field == nil {
		return errors.New("pattern do not match")
	}
	str := fmt.Sprintf("%v", field)

	match, err := regexp.MatchString(string(r), str)

	if err != nil {
		return err
	}

	if !match {
		return errors.New("pattern do not match")

	}
	return nil
}

func (r Range) Validate(ctx context.Context, field interface{}) error {

	ln := 0.0

	// if field is string, then check character count
	if str, ok := field.(string); ok {
		ln = float64(len(str))
	}

	// if field is number, parse it and convert it to integer

	str := fmt.Sprintf("%v", field)
	if n, ok := strconv.ParseFloat(str, 32); ok == nil {
		ln = n
	}

	if r.Max != nil && ln > *r.Max {
		return errors.New("field greater than")
	}

	if r.Min != nil && ln < *r.Min {
		return errors.New("field less than")
	}

	return nil
}
