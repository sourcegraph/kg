package kg

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		into string
		from string
		exp  string
	}{{
		into: `{"one": "asdf"}`,
		from: `{"two": "sdf"}`,
		exp:  `{"one": "asdf", "two": "sdf"}`,
	}, {
		into: `{
			"one": {
				"two": {
					"three": "foo"
				}
			}
		}`,
		from: `{
			"one": {
				"two": {
					"three": "baz"
				},
				"two-b": "foo"
			},
			"one-b": "bar"
		}`,
		exp: `{
			"one": {
				"two": {
					"three": "baz"
				},
				"two-b": "foo"
			},
			"one-b": "bar"
		}`,
	}}

	for _, test := range tests {
		var into, from, exp map[string]interface{}
		json.Unmarshal([]byte(test.into), &into)
		json.Unmarshal([]byte(test.from), &from)
		json.Unmarshal([]byte(test.exp), &exp)

		m := make(map[string]interface{})
		Merge(m, into)
		if !reflect.DeepEqual(m, into) {
			t.Errorf("failed to merge map %+v into empty map", m)
			continue
		}
		Merge(m, from)
		if !reflect.DeepEqual(m, exp) {
			t.Errorf("expected %#v but got %#v", exp, m)
			continue
		}
	}
}
