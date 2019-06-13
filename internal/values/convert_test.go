package values

import (
	"math"
	"testing"
	"time"

	"github.com/ptiger10/pd/options"
)

// func TestConvertAtomic(t *testing.T) {
// 	// testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	// testEpoch := 1556668800000000000
// 	// epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
// 	nan := math.NaN()
// 	var test = []struct {
// 		input        interface{}
// 		originalType options.DataType
// 		convertTo    options.DataType
// 		wantVal      interface{}
// 		wantNull     bool
// 	}{
// 		{math.NaN(), options.Float64, options.Float64, nan, true},
// 		{1.5, options.Float64, options.Float64, 1.5, false},
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
		convertTo options.DataType
		wantVal   interface{}
		wantNull  bool
	}{
		// Float
		{newSliceFloat64([]float64{math.NaN()}), options.Float64, nan, true},
		{newSliceFloat64([]float64{1.5}), options.Float64, 1.5, false},

		{newSliceFloat64([]float64{math.NaN()}), options.Int64, int64(0), true},
		{newSliceFloat64([]float64{1.5}), options.Int64, int64(1), false},

		{newSliceFloat64([]float64{math.NaN()}), options.String, "NaN", true},
		{newSliceFloat64([]float64{1.5}), options.String, "1.5", false},

		{newSliceFloat64([]float64{math.NaN()}), options.Bool, false, true},
		{newSliceFloat64([]float64{1.5}), options.Bool, true, false},

		{newSliceFloat64([]float64{math.NaN()}), options.DateTime, time.Time{}, true},
		{newSliceFloat64([]float64{float64(testEpoch)}), options.DateTime, testDate, false},

		{newSliceFloat64([]float64{math.NaN()}), options.Interface, nan, true},
		{newSliceFloat64([]float64{1.5}), options.Interface, 1.5, false},

		// Int
		{newSliceInt64([]int64{1}), options.Float64, 1.0, false},
		{newSliceInt64([]int64{1}), options.Int64, int64(1), false},
		{newSliceInt64([]int64{1}), options.String, "1", false},
		{newSliceInt64([]int64{1}), options.Bool, true, false},

		{newSliceInt64([]int64{1}), options.DateTime, epochDate, false},
		{newSliceInt64([]int64{int64(testEpoch)}), options.DateTime, testDate, false},

		// String
		{newSliceString([]string{""}), options.Float64, nan, true},
		{newSliceString([]string{"foo"}), options.Float64, nan, true},
		{newSliceString([]string{"1.5"}), options.Float64, 1.5, false},

		{newSliceString([]string{""}), options.Int64, int64(0), true},
		{newSliceString([]string{"foo"}), options.Int64, int64(0), true},
		{newSliceString([]string{"1.5"}), options.Int64, int64(1), false},
		{newSliceString([]string{"1.0"}), options.Int64, int64(1), false},
		{newSliceString([]string{"1"}), options.Int64, int64(1), false},

		{newSliceString([]string{""}), options.String, "NaN", true},
		{newSliceString([]string{"NaN"}), options.String, "NaN", true},
		{newSliceString([]string{"n/a"}), options.String, "NaN", true},
		{newSliceString([]string{"N/A"}), options.String, "NaN", true},
		{newSliceString([]string{"1.5"}), options.String, "1.5", false},
		{newSliceString([]string{"foo"}), options.String, "foo", false},

		{newSliceString([]string{""}), options.Bool, false, true},
		{newSliceString([]string{"foo"}), options.Bool, true, false},

		{newSliceString([]string{"May 1, 2019"}), options.DateTime, testDate, false},
		{newSliceString([]string{"5/1/2019"}), options.DateTime, testDate, false},
		{newSliceString([]string{"2019-05-01"}), options.DateTime, testDate, false},

		// Bool
		{newSliceBool([]bool{true}), options.Float64, 1.0, false},
		{newSliceBool([]bool{false}), options.Float64, 0.0, false},

		{newSliceBool([]bool{true}), options.Int64, int64(1), false},
		{newSliceBool([]bool{false}), options.Int64, int64(0), false},

		{newSliceBool([]bool{true}), options.String, "true", false},
		{newSliceBool([]bool{false}), options.String, "false", false},

		{newSliceBool([]bool{true}), options.Bool, true, false},
		{newSliceBool([]bool{false}), options.Bool, false, false},

		{newSliceBool([]bool{true}), options.DateTime, epochDate, false},
		{newSliceBool([]bool{false}), options.DateTime, epochDate, false},

		// DateTime
		{newSliceDateTime([]time.Time{testDate}), options.Float64, float64(testEpoch), false},
		{newSliceDateTime([]time.Time{time.Time{}}), options.Float64, nan, true},

		{newSliceDateTime([]time.Time{testDate}), options.Int64, int64(testEpoch), false},
		{newSliceDateTime([]time.Time{time.Time{}}), options.Int64, int64(0), true},

		{newSliceDateTime([]time.Time{testDate}), options.String, "2019-05-01 00:00:00 +0000 UTC", false},
		{newSliceDateTime([]time.Time{time.Time{}}), options.String, "NaN", true},

		{newSliceDateTime([]time.Time{testDate}), options.Bool, true, false},
		{newSliceDateTime([]time.Time{time.Time{}}), options.Bool, false, true},

		{newSliceDateTime([]time.Time{testDate}), options.DateTime, testDate, false},
		{newSliceDateTime([]time.Time{time.Time{}}), options.DateTime, time.Time{}, true},

		// Interface
		{newSliceInterface([]interface{}{math.NaN()}), options.Float64, nan, true},
		{newSliceInterface([]interface{}{1.5}), options.Float64, 1.5, false},

		{newSliceInterface([]interface{}{1}), options.Int64, int64(1), false},

		{newSliceInterface([]interface{}{""}), options.String, "NaN", true},
		{newSliceInterface([]interface{}{"NaN"}), options.String, "NaN", true},
		{newSliceInterface([]interface{}{"n/a"}), options.String, "NaN", true},
		{newSliceInterface([]interface{}{"N/A"}), options.String, "NaN", true},
		{newSliceInterface([]interface{}{"1.5"}), options.String, "1.5", false},
		{newSliceInterface([]interface{}{"foo"}), options.String, "foo", false},

		{newSliceInterface([]interface{}{true}), options.Bool, true, false},
		{newSliceInterface([]interface{}{false}), options.Bool, false, false},

		{newSliceInterface([]interface{}{testDate}), options.DateTime, testDate, false},
		{newSliceInterface([]interface{}{time.Time{}}), options.DateTime, time.Time{}, true},
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
		datatype options.DataType
	}{
		{options.None},
		{options.Unsupported},
	}
	for _, test := range tests {
		vals := newSliceFloat64([]float64{1.5})
		_, err := Convert(vals.Values, test.datatype)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.datatype)
		}
	}
}
