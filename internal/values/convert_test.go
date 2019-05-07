package values

import (
	"math"
	"testing"
	"time"

	"github.com/ptiger10/pd/kinds"
)

func TestConvert(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	testEpoch := 1556668800000000000
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	nan := math.NaN()
	var tests = []struct {
		input     Factory
		convertTo kinds.Kind
		wantVal   interface{}
		wantNull  bool
	}{
		// Float
		{SliceFloat([]float64{math.NaN()}), kinds.Float, nan, true},
		{SliceFloat([]float64{1.5}), kinds.Float, 1.5, false},

		{SliceFloat([]float64{math.NaN()}), kinds.Int, int64(0), true},
		{SliceFloat([]float64{1.5}), kinds.Int, int64(1), false},

		{SliceFloat([]float64{math.NaN()}), kinds.String, "NaN", true},
		{SliceFloat([]float64{1.5}), kinds.String, "1.5", false},

		{SliceFloat([]float64{math.NaN()}), kinds.Bool, false, true},
		{SliceFloat([]float64{1.5}), kinds.Bool, true, false},

		{SliceFloat([]float64{math.NaN()}), kinds.DateTime, time.Time{}, true},
		{SliceFloat([]float64{float64(testEpoch)}), kinds.DateTime, testDate, false},

		// Int
		{SliceInt([]int64{1}), kinds.Float, 1.0, false},
		{SliceInt([]int64{1}), kinds.Int, int64(1), false},
		{SliceInt([]int64{1}), kinds.String, "1", false},
		{SliceInt([]int64{1}), kinds.Bool, true, false},

		{SliceInt([]int64{1}), kinds.DateTime, epochDate, false},
		{SliceInt([]int64{int64(testEpoch)}), kinds.DateTime, testDate, false},

		// String
		{SliceString([]string{""}), kinds.Float, nan, true},
		{SliceString([]string{"foo"}), kinds.Float, nan, true},
		{SliceString([]string{"1.5"}), kinds.Float, 1.5, false},

		{SliceString([]string{""}), kinds.Int, int64(0), true},
		{SliceString([]string{"foo"}), kinds.Int, int64(0), true},
		{SliceString([]string{"1.5"}), kinds.Int, int64(1), false},
		{SliceString([]string{"1.0"}), kinds.Int, int64(1), false},
		{SliceString([]string{"1"}), kinds.Int, int64(1), false},

		{SliceString([]string{""}), kinds.String, "NaN", true},
		{SliceString([]string{"NaN"}), kinds.String, "NaN", true},
		{SliceString([]string{"n/a"}), kinds.String, "NaN", true},
		{SliceString([]string{"N/A"}), kinds.String, "NaN", true},
		{SliceString([]string{"1.5"}), kinds.String, "1.5", false},
		{SliceString([]string{"foo"}), kinds.String, "foo", false},

		{SliceString([]string{""}), kinds.Bool, false, true},
		{SliceString([]string{"foo"}), kinds.Bool, true, false},

		{SliceString([]string{"May 1, 2019"}), kinds.DateTime, testDate, false},
		{SliceString([]string{"5/1/2019"}), kinds.DateTime, testDate, false},
		{SliceString([]string{"2019-05-01"}), kinds.DateTime, testDate, false},

		// Bool
		{SliceBool([]bool{true}), kinds.Float, 1.0, false},
		{SliceBool([]bool{false}), kinds.Float, 0.0, false},

		{SliceBool([]bool{true}), kinds.Int, int64(1), false},
		{SliceBool([]bool{false}), kinds.Int, int64(0), false},

		{SliceBool([]bool{true}), kinds.String, "true", false},
		{SliceBool([]bool{false}), kinds.String, "false", false},

		{SliceBool([]bool{true}), kinds.Bool, true, false},
		{SliceBool([]bool{false}), kinds.Bool, false, false},

		{SliceBool([]bool{true}), kinds.DateTime, epochDate, false},
		{SliceBool([]bool{false}), kinds.DateTime, epochDate, false},

		// DateTime
		{SliceDateTime([]time.Time{testDate}), kinds.Float, float64(testEpoch), false},
		{SliceDateTime([]time.Time{time.Time{}}), kinds.Float, nan, true},

		{SliceDateTime([]time.Time{testDate}), kinds.Int, int64(testEpoch), false},
		{SliceDateTime([]time.Time{time.Time{}}), kinds.Int, int64(0), true},

		{SliceDateTime([]time.Time{testDate}), kinds.String, "2019-05-01 00:00:00 +0000 UTC", false},
		{SliceDateTime([]time.Time{time.Time{}}), kinds.String, "NaN", true},

		{SliceDateTime([]time.Time{testDate}), kinds.Bool, true, false},
		{SliceDateTime([]time.Time{time.Time{}}), kinds.Bool, false, true},

		{SliceDateTime([]time.Time{testDate}), kinds.DateTime, testDate, false},
		{SliceDateTime([]time.Time{time.Time{}}), kinds.DateTime, time.Time{}, true},

		// Interface
		{SliceInterface([]interface{}{math.NaN()}), kinds.Float, nan, true},
		{SliceInterface([]interface{}{1.5}), kinds.Float, 1.5, false},

		{SliceInterface([]interface{}{1}), kinds.Int, int64(1), false},

		{SliceInterface([]interface{}{""}), kinds.String, "NaN", true},
		{SliceInterface([]interface{}{"NaN"}), kinds.String, "NaN", true},
		{SliceInterface([]interface{}{"n/a"}), kinds.String, "NaN", true},
		{SliceInterface([]interface{}{"N/A"}), kinds.String, "NaN", true},
		{SliceInterface([]interface{}{"1.5"}), kinds.String, "1.5", false},
		{SliceInterface([]interface{}{"foo"}), kinds.String, "foo", false},

		{SliceInterface([]interface{}{true}), kinds.Bool, true, false},
		{SliceInterface([]interface{}{false}), kinds.Bool, false, false},

		{SliceInterface([]interface{}{testDate}), kinds.DateTime, testDate, false},
		{SliceInterface([]interface{}{time.Time{}}), kinds.DateTime, time.Time{}, true},
	}

	for _, test := range tests {
		converted, err := Convert(test.input.V, test.convertTo)
		if err != nil {
			t.Errorf("Unable to convert to string: %v", err)
		}
		elem := converted.Element(0)
		val := elem[0]
		null := elem[1].(bool)
		if val != test.wantVal {
			// special case to check for two floats equaling NaN
			f1, ok1 := val.(float64)
			f2, ok2 := test.wantVal.(float64)
			if !ok1 || !ok2 || (math.IsNaN(f1) && !math.IsNaN(f2)) || (!math.IsNaN(f1) && math.IsNaN(f2)) {
				t.Errorf("%v conversion of %T (%v) returned %v, want %v", test.convertTo, test.input, test.input, val, test.wantVal)
			}

		}
		if null != test.wantNull {
			t.Errorf("%v conversion of %T (%v) returned null as %v, want %v", test.convertTo, test.input, test.input, null, test.wantNull)
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
		vals := SliceFloat([]float64{1.5})
		_, err := Convert(vals.V, test.kind)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.kind)
		}
	}
}
