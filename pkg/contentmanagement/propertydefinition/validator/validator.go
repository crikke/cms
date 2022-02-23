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

type RequiredRule bool
type RegexRule string
type RangeRule struct {
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
			return RequiredRule(b), nil
		}
		return nil, errors.New("parse error: cannot parse into type RequiredField")
	case "pattern":
		if str, ok := val.(string); ok {
			return RegexRule(str), nil
		}
	case "range":
		if r, ok := val.(RangeRule); ok {
			return r, nil
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

func (r RangeRule) Validate(ctx context.Context, field interface{}) error {

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
