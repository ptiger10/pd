package values

import (
	"math"
	"testing"
	"time"

	"github.com/ptiger10/pd/datatypes"
)

// func TestConvertAtomic(t *testing.T) {
// 	// testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	// testEpoch := 1556668800000000000
// 	// epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
// 	nan := math.NaN()
// 	var test = []struct {
// 		input        interface{}
// 		originalType datatypes.DataType
// 		convertTo    datatypes.DataType
// 		wantVal      interface{}
// 		wantNull     bool
// 	}{
// 		{math.NaN(), datatypes.Float64, datatypes.Float64, nan, true},
// 		{1.5, datatypes.Float64, datatypes.Float64, 1.5, false},
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
		convertTo datatypes.DataType
		wantVal   interface{}
		wantNull  bool
	}{
		// Float
		{newSliceFloat64([]float64{math.NaN()}), datatypes.Float64, nan, true},
		{newSliceFloat64([]float64{1.5}), datatypes.Float64, 1.5, false},

		{newSliceFloat64([]float64{math.NaN()}), datatypes.Int64, int64(0), true},
		{newSliceFloat64([]float64{1.5}), datatypes.Int64, int64(1), false},

		{newSliceFloat64([]float64{math.NaN()}), datatypes.String, "NaN", true},
		{newSliceFloat64([]float64{1.5}), datatypes.String, "1.5", false},

		{newSliceFloat64([]float64{math.NaN()}), datatypes.Bool, false, true},
		{newSliceFloat64([]float64{1.5}), datatypes.Bool, true, false},

		{newSliceFloat64([]float64{math.NaN()}), datatypes.DateTime, time.Time{}, true},
		{newSliceFloat64([]float64{float64(testEpoch)}), datatypes.DateTime, testDate, false},

		{newSliceFloat64([]float64{math.NaN()}), datatypes.Interface, nan, true},
		{newSliceFloat64([]float64{1.5}), datatypes.Interface, 1.5, false},

		// Int
		{NewSliceInt([]int64{1}), datatypes.Float64, 1.0, false},
		{NewSliceInt([]int64{1}), datatypes.Int64, int64(1), false},
		{NewSliceInt([]int64{1}), datatypes.String, "1", false},
		{NewSliceInt([]int64{1}), datatypes.Bool, true, false},

		{NewSliceInt([]int64{1}), datatypes.DateTime, epochDate, false},
		{NewSliceInt([]int64{int64(testEpoch)}), datatypes.DateTime, testDate, false},

		// String
		{newSliceString([]string{""}), datatypes.Float64, nan, true},
		{newSliceString([]string{"foo"}), datatypes.Float64, nan, true},
		{newSliceString([]string{"1.5"}), datatypes.Float64, 1.5, false},

		{newSliceString([]string{""}), datatypes.Int64, int64(0), true},
		{newSliceString([]string{"foo"}), datatypes.Int64, int64(0), true},
		{newSliceString([]string{"1.5"}), datatypes.Int64, int64(1), false},
		{newSliceString([]string{"1.0"}), datatypes.Int64, int64(1), false},
		{newSliceString([]string{"1"}), datatypes.Int64, int64(1), false},

		{newSliceString([]string{""}), datatypes.String, "NaN", true},
		{newSliceString([]string{"NaN"}), datatypes.String, "NaN", true},
		{newSliceString([]string{"n/a"}), datatypes.String, "NaN", true},
		{newSliceString([]string{"N/A"}), datatypes.String, "NaN", true},
		{newSliceString([]string{"1.5"}), datatypes.String, "1.5", false},
		{newSliceString([]string{"foo"}), datatypes.String, "foo", false},

		{newSliceString([]string{""}), datatypes.Bool, false, true},
		{newSliceString([]string{"foo"}), datatypes.Bool, true, false},

		{newSliceString([]string{"May 1, 2019"}), datatypes.DateTime, testDate, false},
		{newSliceString([]string{"5/1/2019"}), datatypes.DateTime, testDate, false},
		{newSliceString([]string{"2019-05-01"}), datatypes.DateTime, testDate, false},

		// Bool
		{newSliceBool([]bool{true}), datatypes.Float64, 1.0, false},
		{newSliceBool([]bool{false}), datatypes.Float64, 0.0, false},

		{newSliceBool([]bool{true}), datatypes.Int64, int64(1), false},
		{newSliceBool([]bool{false}), datatypes.Int64, int64(0), false},

		{newSliceBool([]bool{true}), datatypes.String, "true", false},
		{newSliceBool([]bool{false}), datatypes.String, "false", false},

		{newSliceBool([]bool{true}), datatypes.Bool, true, false},
		{newSliceBool([]bool{false}), datatypes.Bool, false, false},

		{newSliceBool([]bool{true}), datatypes.DateTime, epochDate, false},
		{newSliceBool([]bool{false}), datatypes.DateTime, epochDate, false},

		// DateTime
		{newSliceDateTime([]time.Time{testDate}), datatypes.Float64, float64(testEpoch), false},
		{newSliceDateTime([]time.Time{time.Time{}}), datatypes.Float64, nan, true},

		{newSliceDateTime([]time.Time{testDate}), datatypes.Int64, int64(testEpoch), false},
		{newSliceDateTime([]time.Time{time.Time{}}), datatypes.Int64, int64(0), true},

		{newSliceDateTime([]time.Time{testDate}), datatypes.String, "2019-05-01 00:00:00 +0000 UTC", false},
		{newSliceDateTime([]time.Time{time.Time{}}), datatypes.String, "NaN", true},

		{newSliceDateTime([]time.Time{testDate}), datatypes.Bool, true, false},
		{newSliceDateTime([]time.Time{time.Time{}}), datatypes.Bool, false, true},

		{newSliceDateTime([]time.Time{testDate}), datatypes.DateTime, testDate, false},
		{newSliceDateTime([]time.Time{time.Time{}}), datatypes.DateTime, time.Time{}, true},

		// Interface
		{newSliceInterface([]interface{}{math.NaN()}), datatypes.Float64, nan, true},
		{newSliceInterface([]interface{}{1.5}), datatypes.Float64, 1.5, false},

		{newSliceInterface([]interface{}{1}), datatypes.Int64, int64(1), false},

		{newSliceInterface([]interface{}{""}), datatypes.String, "NaN", true},
		{newSliceInterface([]interface{}{"NaN"}), datatypes.String, "NaN", true},
		{newSliceInterface([]interface{}{"n/a"}), datatypes.String, "NaN", true},
		{newSliceInterface([]interface{}{"N/A"}), datatypes.String, "NaN", true},
		{newSliceInterface([]interface{}{"1.5"}), datatypes.String, "1.5", false},
		{newSliceInterface([]interface{}{"foo"}), datatypes.String, "foo", false},

		{newSliceInterface([]interface{}{true}), datatypes.Bool, true, false},
		{newSliceInterface([]interface{}{false}), datatypes.Bool, false, false},

		{newSliceInterface([]interface{}{testDate}), datatypes.DateTime, testDate, false},
		{newSliceInterface([]interface{}{time.Time{}}), datatypes.DateTime, time.Time{}, true},
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
		datatype datatypes.DataType
	}{
		{datatypes.None},
		{datatypes.Unsupported},
	}
	for _, test := range tests {
		vals := newSliceFloat64([]float64{1.5})
		_, err := Convert(vals.Values, test.datatype)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.datatype)
		}
	}
}
