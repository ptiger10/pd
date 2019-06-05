package values

import (
	"math"
	"testing"
	"time"

	"github.com/ptiger10/pd/kinds"
)

// func TestConvertAtomic(t *testing.T) {
// 	// testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	// testEpoch := 1556668800000000000
// 	// epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
// 	nan := math.NaN()
// 	var test = []struct {
// 		input        interface{}
// 		originalType kinds.Kind
// 		convertTo    kinds.Kind
// 		wantVal      interface{}
// 		wantNull     bool
// 	}{
// 		{math.NaN(), kinds.Float, kinds.Float, nan, true},
// 		{1.5, kinds.Float, kinds.Float, 1.5, false},
// 	}
// 	for _, test := range tests {
// 		converted, err := Convert(test.input.V, test.convertTo)
// 		if err != nil {
// 			t.Errorf("Unable to convert to string: %v", err)
// 		}
// 		elem := converted.Element(0)
// 	}

// }

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
		{newSliceFloat([]float64{math.NaN()}), kinds.Float, nan, true},
		{newSliceFloat([]float64{1.5}), kinds.Float, 1.5, false},

		{newSliceFloat([]float64{math.NaN()}), kinds.Int, int64(0), true},
		{newSliceFloat([]float64{1.5}), kinds.Int, int64(1), false},

		{newSliceFloat([]float64{math.NaN()}), kinds.String, "NaN", true},
		{newSliceFloat([]float64{1.5}), kinds.String, "1.5", false},

		{newSliceFloat([]float64{math.NaN()}), kinds.Bool, false, true},
		{newSliceFloat([]float64{1.5}), kinds.Bool, true, false},

		{newSliceFloat([]float64{math.NaN()}), kinds.DateTime, time.Time{}, true},
		{newSliceFloat([]float64{float64(testEpoch)}), kinds.DateTime, testDate, false},

		{newSliceFloat([]float64{math.NaN()}), kinds.Interface, nan, true},
		{newSliceFloat([]float64{1.5}), kinds.Interface, 1.5, false},

		// Int
		{NewSliceInt([]int64{1}), kinds.Float, 1.0, false},
		{NewSliceInt([]int64{1}), kinds.Int, int64(1), false},
		{NewSliceInt([]int64{1}), kinds.String, "1", false},
		{NewSliceInt([]int64{1}), kinds.Bool, true, false},

		{NewSliceInt([]int64{1}), kinds.DateTime, epochDate, false},
		{NewSliceInt([]int64{int64(testEpoch)}), kinds.DateTime, testDate, false},

		// String
		{newSliceString([]string{""}), kinds.Float, nan, true},
		{newSliceString([]string{"foo"}), kinds.Float, nan, true},
		{newSliceString([]string{"1.5"}), kinds.Float, 1.5, false},

		{newSliceString([]string{""}), kinds.Int, int64(0), true},
		{newSliceString([]string{"foo"}), kinds.Int, int64(0), true},
		{newSliceString([]string{"1.5"}), kinds.Int, int64(1), false},
		{newSliceString([]string{"1.0"}), kinds.Int, int64(1), false},
		{newSliceString([]string{"1"}), kinds.Int, int64(1), false},

		{newSliceString([]string{""}), kinds.String, "NaN", true},
		{newSliceString([]string{"NaN"}), kinds.String, "NaN", true},
		{newSliceString([]string{"n/a"}), kinds.String, "NaN", true},
		{newSliceString([]string{"N/A"}), kinds.String, "NaN", true},
		{newSliceString([]string{"1.5"}), kinds.String, "1.5", false},
		{newSliceString([]string{"foo"}), kinds.String, "foo", false},

		{newSliceString([]string{""}), kinds.Bool, false, true},
		{newSliceString([]string{"foo"}), kinds.Bool, true, false},

		{newSliceString([]string{"May 1, 2019"}), kinds.DateTime, testDate, false},
		{newSliceString([]string{"5/1/2019"}), kinds.DateTime, testDate, false},
		{newSliceString([]string{"2019-05-01"}), kinds.DateTime, testDate, false},

		// Bool
		{newSliceBool([]bool{true}), kinds.Float, 1.0, false},
		{newSliceBool([]bool{false}), kinds.Float, 0.0, false},

		{newSliceBool([]bool{true}), kinds.Int, int64(1), false},
		{newSliceBool([]bool{false}), kinds.Int, int64(0), false},

		{newSliceBool([]bool{true}), kinds.String, "true", false},
		{newSliceBool([]bool{false}), kinds.String, "false", false},

		{newSliceBool([]bool{true}), kinds.Bool, true, false},
		{newSliceBool([]bool{false}), kinds.Bool, false, false},

		{newSliceBool([]bool{true}), kinds.DateTime, epochDate, false},
		{newSliceBool([]bool{false}), kinds.DateTime, epochDate, false},

		// DateTime
		{newSliceDateTime([]time.Time{testDate}), kinds.Float, float64(testEpoch), false},
		{newSliceDateTime([]time.Time{time.Time{}}), kinds.Float, nan, true},

		{newSliceDateTime([]time.Time{testDate}), kinds.Int, int64(testEpoch), false},
		{newSliceDateTime([]time.Time{time.Time{}}), kinds.Int, int64(0), true},

		{newSliceDateTime([]time.Time{testDate}), kinds.String, "2019-05-01 00:00:00 +0000 UTC", false},
		{newSliceDateTime([]time.Time{time.Time{}}), kinds.String, "NaN", true},

		{newSliceDateTime([]time.Time{testDate}), kinds.Bool, true, false},
		{newSliceDateTime([]time.Time{time.Time{}}), kinds.Bool, false, true},

		{newSliceDateTime([]time.Time{testDate}), kinds.DateTime, testDate, false},
		{newSliceDateTime([]time.Time{time.Time{}}), kinds.DateTime, time.Time{}, true},

		// Interface
		{newSliceInterface([]interface{}{math.NaN()}), kinds.Float, nan, true},
		{newSliceInterface([]interface{}{1.5}), kinds.Float, 1.5, false},

		{newSliceInterface([]interface{}{1}), kinds.Int, int64(1), false},

		{newSliceInterface([]interface{}{""}), kinds.String, "NaN", true},
		{newSliceInterface([]interface{}{"NaN"}), kinds.String, "NaN", true},
		{newSliceInterface([]interface{}{"n/a"}), kinds.String, "NaN", true},
		{newSliceInterface([]interface{}{"N/A"}), kinds.String, "NaN", true},
		{newSliceInterface([]interface{}{"1.5"}), kinds.String, "1.5", false},
		{newSliceInterface([]interface{}{"foo"}), kinds.String, "foo", false},

		{newSliceInterface([]interface{}{true}), kinds.Bool, true, false},
		{newSliceInterface([]interface{}{false}), kinds.Bool, false, false},

		{newSliceInterface([]interface{}{testDate}), kinds.DateTime, testDate, false},
		{newSliceInterface([]interface{}{time.Time{}}), kinds.DateTime, time.Time{}, true},
	}

	for _, test := range tests {
		converted, err := Convert(test.input.Values, test.convertTo)
		if err != nil {
			t.Errorf("Unable to convert to string: %v", err)
		}
		elem := converted.Element(0)
		val := elem.Value
		null := elem.Null
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
		vals := newSliceFloat([]float64{1.5})
		_, err := Convert(vals.Values, test.kind)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.kind)
		}
	}
}
