package index

import (
	"testing"
	"time"

	"github.com/ptiger10/pd/options"
)

func TestConvertIndex_int(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		lvl       Level
		convertTo options.DataType
	}{
		// Float
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), options.Float64},
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), options.Int64},
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), options.String},
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), options.Bool},
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), options.DateTime},

		// Int
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), options.Float64},
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), options.Int64},
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), options.String},
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), options.Bool},
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), options.DateTime},

		// String
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), options.Float64},
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), options.Int64},
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), options.String},
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), options.Bool},
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), options.DateTime},

		// Bool
		{MustCreateNewLevel([]bool{true, false, false}, ""), options.Float64},
		{MustCreateNewLevel([]bool{true, false, false}, ""), options.Int64},
		{MustCreateNewLevel([]bool{true, false, false}, ""), options.String},
		{MustCreateNewLevel([]bool{true, false, false}, ""), options.Bool},
		{MustCreateNewLevel([]bool{true, false, false}, ""), options.DateTime},

		// DateTime
		{MustCreateNewLevel([]time.Time{testDate}, ""), options.Float64},
		{MustCreateNewLevel([]time.Time{testDate}, ""), options.Int64},
		{MustCreateNewLevel([]time.Time{testDate}, ""), options.String},
		{MustCreateNewLevel([]time.Time{testDate}, ""), options.Bool},
		{MustCreateNewLevel([]time.Time{testDate}, ""), options.DateTime},

		// Interface
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), options.Float64},
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), options.Int64},
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), options.String},
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), options.Bool},
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), options.DateTime},
	}
	for _, test := range tests {
		lvl, err := test.lvl.Convert(test.convertTo)
		if err != nil {
			t.Error(err)
		}
		if lvl.DataType != test.convertTo {
			t.Errorf("Attempted conversion to %v returned %v", test.convertTo, lvl.DataType)
		}
	}
}

func TestConvert_Numeric_Datetime(t *testing.T) {
	n := int64(1556668800000000000)
	wantVal := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		lvl Level
	}{
		{MustCreateNewLevel([]int64{n}, "")},
		{MustCreateNewLevel([]float64{float64(n)}, "")},
	}
	for _, test := range tests {
		lvl, _ := test.lvl.Convert(options.DateTime)
		elem := lvl.Labels.Element(0)
		gotVal := elem.Value.(time.Time)
		if gotVal != wantVal {
			t.Errorf("Error converting %v to datetime: returned %v, want %v", test.lvl, gotVal, wantVal)
		}
	}
}

func TestConvert_Unsupported(t *testing.T) {
	var tests = []struct {
		datatype options.DataType
	}{
		{options.None},
		{options.Unsupported},
	}
	for _, test := range tests {
		lvl := MustCreateNewLevel([]float64{1, 2, 3}, "")
		_, err := lvl.Convert(test.datatype)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.datatype)
		}
	}
}
