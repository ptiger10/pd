package index

import (
	"testing"
	"time"

	"github.com/ptiger10/pd/kinds"
)

func TestConvertIndex_int(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		lvl       Level
		convertTo kinds.Kind
	}{
		// Float
		{mustCreateNewLevel([]float64{1, 2, 3}), kinds.Float64},
		{mustCreateNewLevel([]float64{1, 2, 3}), kinds.Int64},
		{mustCreateNewLevel([]float64{1, 2, 3}), kinds.String},
		{mustCreateNewLevel([]float64{1, 2, 3}), kinds.Bool},
		{mustCreateNewLevel([]float64{1, 2, 3}), kinds.DateTime},

		// Int
		{mustCreateNewLevel([]int64{1, 2, 3}), kinds.Float64},
		{mustCreateNewLevel([]int64{1, 2, 3}), kinds.Int64},
		{mustCreateNewLevel([]int64{1, 2, 3}), kinds.String},
		{mustCreateNewLevel([]int64{1, 2, 3}), kinds.Bool},
		{mustCreateNewLevel([]int64{1, 2, 3}), kinds.DateTime},

		// String
		{mustCreateNewLevel([]string{"1", "2", "3"}), kinds.Float64},
		{mustCreateNewLevel([]string{"1", "2", "3"}), kinds.Int64},
		{mustCreateNewLevel([]string{"1", "2", "3"}), kinds.String},
		{mustCreateNewLevel([]string{"1", "2", "3"}), kinds.Bool},
		{mustCreateNewLevel([]string{"1", "2", "3"}), kinds.DateTime},

		// Bool
		{mustCreateNewLevel([]bool{true, false, false}), kinds.Float64},
		{mustCreateNewLevel([]bool{true, false, false}), kinds.Int64},
		{mustCreateNewLevel([]bool{true, false, false}), kinds.String},
		{mustCreateNewLevel([]bool{true, false, false}), kinds.Bool},
		{mustCreateNewLevel([]bool{true, false, false}), kinds.DateTime},

		// DateTime
		{mustCreateNewLevel([]time.Time{testDate}), kinds.Float64},
		{mustCreateNewLevel([]time.Time{testDate}), kinds.Int64},
		{mustCreateNewLevel([]time.Time{testDate}), kinds.String},
		{mustCreateNewLevel([]time.Time{testDate}), kinds.Bool},
		{mustCreateNewLevel([]time.Time{testDate}), kinds.DateTime},

		// Interface
		{mustCreateNewLevel([]interface{}{1, "2", true}), kinds.Float64},
		{mustCreateNewLevel([]interface{}{1, "2", true}), kinds.Int64},
		{mustCreateNewLevel([]interface{}{1, "2", true}), kinds.String},
		{mustCreateNewLevel([]interface{}{1, "2", true}), kinds.Bool},
		{mustCreateNewLevel([]interface{}{1, "2", true}), kinds.DateTime},
	}
	for _, test := range tests {
		lvl, err := test.lvl.Convert(test.convertTo)
		if err != nil {
			t.Error(err)
		}
		if lvl.Kind != test.convertTo {
			t.Errorf("Attempted conversion to %v returned %v", test.convertTo, lvl.Kind)
		}
	}
}

func TestConvert_Numeric_Datetime(t *testing.T) {
	n := int64(1556668800000000000)
	wantVal := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		lvl Level
	}{
		{mustCreateNewLevel([]int64{n})},
		{mustCreateNewLevel([]float64{float64(n)})},
	}
	for _, test := range tests {
		lvl, _ := test.lvl.Convert(kinds.DateTime)
		elem := lvl.Labels.Element(0)
		gotVal := elem.Value.(time.Time)
		if gotVal != wantVal {
			t.Errorf("Error converting %v to datetime: returned %v, want %v", test.lvl, gotVal, wantVal)
		}
	}
}

func TestConvert_Unsupported(t *testing.T) {
	var tests = []struct {
		kind kinds.Kind
	}{
		{kinds.None},
		{kinds.Unsupported},
	}
	for _, test := range tests {
		lvl := mustCreateNewLevel([]float64{1, 2, 3})
		_, err := lvl.Convert(test.kind)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.kind)
		}
	}
}
