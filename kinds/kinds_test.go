package kinds

import (
	"testing"
)

func TestKind(t *testing.T) {
	var tests = []struct {
		kind     Kind
		expected string
	}{

		{None, "none"},
		{Float, "float64"},
		{Int, "int64"},
		{String, "string"},
		{Bool, "bool"},
		{DateTime, "time.Time"},
		{Interface, "interface"},
		{Unsupported, "unsupported"},
		{-1, "unknown"},
		{100, "unknown"},
	}
	for _, test := range tests {
		if test.kind.String() != test.expected {
			t.Errorf("Kind.String() for kind %v returned %v, want %v", test.kind, test.kind.String(), test.expected)
		}
	}
}
