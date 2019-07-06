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

func TestGetDataType(t *testing.T) {
	var tests = []struct {
		expected DataType
		dataType string
	}{
		{Float64, "float"},
		{Float64, "float64"},
		{Float64, "Float64"},
		{Int64, "int"},
		{Int64, "int64"},
		{String, "string"},
		{String, "STRING"},
		{Bool, "bool"},
		{DateTime, "dateTime"},
		{Interface, "interface"},
		{Unsupported, "other"},
	}
	for _, tt := range tests {
		got := DT(tt.dataType)
		if got != tt.expected {
			t.Errorf("DT() = %v, want %v", got, tt.expected)
		}
	}
}
