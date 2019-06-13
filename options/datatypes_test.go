package options

import (
	"testing"
)

func TestDataType(t *testing.T) {
	var tests = []struct {
		DataType DataType
		expected string
	}{

		{None, "none"},
		{Float64, "float64"},
		{Int64, "int64"},
		{String, "string"},
		{Bool, "bool"},
		{DateTime, "dateTime"},
		{Interface, "interface"},
		{Unsupported, "unsupported"},
		{-1, "unknown"},
		{100, "unknown"},
	}
	for _, test := range tests {
		if test.DataType.String() != test.expected {
			t.Errorf("DataType.String() for DataType %v returned %v, want %v", test.DataType, test.DataType.String(), test.expected)
		}
	}
}
