package series

import (
	"testing"

	"github.com/ptiger10/pd/kinds"
)

func TestKind(t *testing.T) {
	var tests = []struct {
		kind     kinds.Kind
		expected string
	}{

		{kinds.Invalid, "invalid"},
		{kinds.Float, "float64"},
		{kinds.Int, "int64"},
		{kinds.String, "string"},
		{kinds.Bool, "bool"},
		{kinds.DateTime, "time.Time"},
		{kinds.Interface, "interface"},
		{kinds.Unsupported, "unsupported"},
		{-1, "unknown"},
		{100, "unknown"},
	}
	for _, test := range tests {
		s, _ := New([]int{1, 2, 3})
		s.kind = test.kind
		if s.Kind() != test.expected {
			t.Errorf("s.Kind() for kind %v returned %v, want %v", test.kind, test.kind.String(), test.expected)
		}
	}
}
