package index

import (
	"testing"
	"time"

	"github.com/ptiger10/pd/datatypes"
)

func TestConvertIndex_int(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		lvl       Level
		convertTo datatypes.DataType
	}{
		// Float
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), datatypes.Float64},
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), datatypes.Int64},
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), datatypes.String},
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), datatypes.Bool},
		{MustCreateNewLevel([]float64{1, 2, 3}, ""), datatypes.DateTime},

		// Int
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), datatypes.Float64},
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), datatypes.Int64},
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), datatypes.String},
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), datatypes.Bool},
		{MustCreateNewLevel([]int64{1, 2, 3}, ""), datatypes.DateTime},

		// String
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), datatypes.Float64},
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), datatypes.Int64},
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), datatypes.String},
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), datatypes.Bool},
		{MustCreateNewLevel([]string{"1", "2", "3"}, ""), datatypes.DateTime},

		// Bool
		{MustCreateNewLevel([]bool{true, false, false}, ""), datatypes.Float64},
		{MustCreateNewLevel([]bool{true, false, false}, ""), datatypes.Int64},
		{MustCreateNewLevel([]bool{true, false, false}, ""), datatypes.String},
		{MustCreateNewLevel([]bool{true, false, false}, ""), datatypes.Bool},
		{MustCreateNewLevel([]bool{true, false, false}, ""), datatypes.DateTime},

		// DateTime
		{MustCreateNewLevel([]time.Time{testDate}, ""), datatypes.Float64},
		{MustCreateNewLevel([]time.Time{testDate}, ""), datatypes.Int64},
		{MustCreateNewLevel([]time.Time{testDate}, ""), datatypes.String},
		{MustCreateNewLevel([]time.Time{testDate}, ""), datatypes.Bool},
		{MustCreateNewLevel([]time.Time{testDate}, ""), datatypes.DateTime},

		// Interface
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), datatypes.Float64},
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), datatypes.Int64},
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), datatypes.String},
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), datatypes.Bool},
		{MustCreateNewLevel([]interface{}{1, "2", true}, ""), datatypes.DateTime},
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
		lvl, _ := test.lvl.Convert(datatypes.DateTime)
		elem := lvl.Labels.Element(0)
		gotVal := elem.Value.(time.Time)
		if gotVal != wantVal {
			t.Errorf("Error converting %v to datetime: returned %v, want %v", test.lvl, gotVal, wantVal)
		}
	}
}

func TestConvert_Unsupported(t *testing.T) {
	var tests = []struct {
		datatype datatypes.DataType
	}{
		{datatypes.None},
		{datatypes.Unsupported},
	}
	for _, test := range tests {
		lvl := MustCreateNewLevel([]float64{1, 2, 3}, "")
		_, err := lvl.Convert(test.datatype)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.datatype)
		}
	}
}
