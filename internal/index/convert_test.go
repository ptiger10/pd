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
		{mustCreateNewLevelFromSlice([]float64{1, 2, 3}), kinds.Float},
		{mustCreateNewLevelFromSlice([]float64{1, 2, 3}), kinds.Int},
		{mustCreateNewLevelFromSlice([]float64{1, 2, 3}), kinds.String},
		{mustCreateNewLevelFromSlice([]float64{1, 2, 3}), kinds.Bool},
		{mustCreateNewLevelFromSlice([]float64{1, 2, 3}), kinds.DateTime},

		// Int
		{mustCreateNewLevelFromSlice([]int64{1, 2, 3}), kinds.Float},
		{mustCreateNewLevelFromSlice([]int64{1, 2, 3}), kinds.Int},
		{mustCreateNewLevelFromSlice([]int64{1, 2, 3}), kinds.String},
		{mustCreateNewLevelFromSlice([]int64{1, 2, 3}), kinds.Bool},
		{mustCreateNewLevelFromSlice([]int64{1, 2, 3}), kinds.DateTime},

		// String
		{mustCreateNewLevelFromSlice([]string{"1", "2", "3"}), kinds.Float},
		{mustCreateNewLevelFromSlice([]string{"1", "2", "3"}), kinds.Int},
		{mustCreateNewLevelFromSlice([]string{"1", "2", "3"}), kinds.String},
		{mustCreateNewLevelFromSlice([]string{"1", "2", "3"}), kinds.Bool},
		{mustCreateNewLevelFromSlice([]string{"1", "2", "3"}), kinds.DateTime},

		// Bool
		{mustCreateNewLevelFromSlice([]bool{true, false, false}), kinds.Float},
		{mustCreateNewLevelFromSlice([]bool{true, false, false}), kinds.Int},
		{mustCreateNewLevelFromSlice([]bool{true, false, false}), kinds.String},
		{mustCreateNewLevelFromSlice([]bool{true, false, false}), kinds.Bool},
		{mustCreateNewLevelFromSlice([]bool{true, false, false}), kinds.DateTime},

		// DateTime
		{mustCreateNewLevelFromSlice([]time.Time{testDate}), kinds.Float},
		{mustCreateNewLevelFromSlice([]time.Time{testDate}), kinds.Int},
		{mustCreateNewLevelFromSlice([]time.Time{testDate}), kinds.String},
		{mustCreateNewLevelFromSlice([]time.Time{testDate}), kinds.Bool},
		{mustCreateNewLevelFromSlice([]time.Time{testDate}), kinds.DateTime},

		// Interface
		{mustCreateNewLevelFromSlice([]interface{}{1, "2", true}), kinds.Float},
		{mustCreateNewLevelFromSlice([]interface{}{1, "2", true}), kinds.Int},
		{mustCreateNewLevelFromSlice([]interface{}{1, "2", true}), kinds.String},
		{mustCreateNewLevelFromSlice([]interface{}{1, "2", true}), kinds.Bool},
		{mustCreateNewLevelFromSlice([]interface{}{1, "2", true}), kinds.DateTime},
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
		{mustCreateNewLevelFromSlice([]int64{n})},
		{mustCreateNewLevelFromSlice([]float64{float64(n)})},
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
		lvl := mustCreateNewLevelFromSlice([]float64{1, 2, 3})
		_, err := lvl.Convert(test.kind)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.kind)
		}
	}
}
