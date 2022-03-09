//go:build unit

package validator

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RangeRule(t *testing.T) {

	tests := []struct {
		name   string
		rule   Range
		inputs map[interface{}]bool
	}{
		{
			name: "max only",
			rule: Range{
				Max: makeFloat64(10),
			},
			inputs: map[interface{}]bool{
				10:      true,
				"banan": true,
				9.11:    true,
				11:      false,
				10.1:    false,
			},
		},
		{
			name: "max & min set",
			rule: Range{
				Max: makeFloat64(10),
				Min: makeFloat64(5),
			},
			inputs: map[interface{}]bool{
				10:      true,
				"aaaa":  false,
				9.11:    true,
				"aaaaa": true,
				3.14:    false,
				3:       false,
				nil:     false,
				"":      false,
			},
		},
	}

	for _, test := range tests {

		for input, ok := range test.inputs {
			t.Run(fmt.Sprintf("%s_%v", test.name, input), func(t *testing.T) {
				err := test.rule.Validate(context.Background(), input)

				if ok {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		}
	}
}

func makeFloat64(n float64) *float64 {
	return &n
}

func Test_RegexRule(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		inputs  map[interface{}]bool
	}{
		{
			name:    `match ^(\d|\w)`,
			pattern: `^(\d|\w)`,
			inputs: map[interface{}]bool{
				"foobar": true,
				3.14:     true,
				"!!foo":  false,
				"":       false,
				nil:      false,
			},
		},
	}

	for _, test := range tests {

		for input, ok := range test.inputs {

			t.Run(fmt.Sprintf("%s_%v", test.name, input), func(t *testing.T) {
				r := Regex(test.pattern)

				err := r.Validate(context.Background(), input)
				if ok {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		}
	}
}

func Test_RequiredRule(t *testing.T) {

	tests := []struct {
		name     string
		required bool
		inputs   map[interface{}]bool
	}{
		{
			name:     "required",
			required: true,
			inputs: map[interface{}]bool{
				"foo": true,
				"":    false,
				nil:   false,
			},
		},
		{
			name:     "not required",
			required: false,
			inputs: map[interface{}]bool{
				"foo": true,
				"":    true,
				nil:   true,
			},
		},
	}

	for _, test := range tests {

		for input, ok := range test.inputs {

			t.Run(fmt.Sprintf("%s_%v", test.name, input), func(t *testing.T) {

				r := Required(test.required)

				err := r.Validate(context.Background(), input)
				if ok {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		}
	}
}
