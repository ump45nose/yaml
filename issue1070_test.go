package yaml

import (
	"fmt"
	"testing"
)

func TestIssue1070FlowMapQuestionInPlain(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"a: {b: c?c}", "c?c"},
		{"a: {b: c?c, d: e}", "c?c"},
		{"[c?c]", "c?c"},
		{"b: c?c", "c?c"},
	}
	for _, tc := range cases {
		var v any
		if err := Unmarshal([]byte(tc.in), &v); err != nil {
			t.Fatalf("%q: unmarshal: %v", tc.in, err)
		}
		got := findString(v)
		if got != tc.want {
			t.Fatalf("%q: got %q want %q; tree=%#v", tc.in, got, tc.want, v)
		}
	}
}

func findString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case map[string]any:
		for _, k := range []string{"b", "a"} {
			if s, ok := x[k].(string); ok {
				return s
			}
			if m, ok := x[k].(map[string]any); ok {
				if s := findString(m); s != "" {
					return s
				}
			}
		}
		for _, vv := range x {
			if s := findString(vv); s != "" {
				return s
			}
		}
	case []any:
		if len(x) > 0 {
			return findString(x[0])
		}
	}
	return fmt.Sprint(v)
}
