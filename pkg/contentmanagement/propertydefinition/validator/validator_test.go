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
		rule   RangeRule
		inputs map[interface{}]bool
	}{
		{
			name: "max range set",
			rule: RangeRule{
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
	}

	for _, test := range tests {

		for input, expect := range test.inputs {
			t.Run(fmt.Sprintf("%s_%v", test.name, input), func(t *testing.T) {
				err := test.rule.Validate(context.Background(), input)

				if expect {
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
