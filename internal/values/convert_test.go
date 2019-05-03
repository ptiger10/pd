package values_test

import (
	"math"
	"testing"
	"time"

	"github.com/ptiger10/pd/internal/values"
	constructVal "github.com/ptiger10/pd/internal/values/constructors"

	"github.com/ptiger10/pd/kinds"
)

func TestConvert(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	testEpoch := 1556668800000000000
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	nan := math.NaN()
	var tests = []struct {
		input     values.Values
		convertTo kinds.Kind
		wantVal   interface{}
		wantNull  bool
	}{
		// Float
		{constructVal.SliceFloat([]float64{math.NaN()}), kinds.Float, nan, true},
		{constructVal.SliceFloat([]float64{1.5}), kinds.Float, 1.5, false},

		{constructVal.SliceFloat([]float64{math.NaN()}), kinds.Int, int64(0), true},
		{constructVal.SliceFloat([]float64{1.5}), kinds.Int, int64(1), false},

		{constructVal.SliceFloat([]float64{math.NaN()}), kinds.String, "NaN", true},
		{constructVal.SliceFloat([]float64{1.5}), kinds.String, "1.5", false},

		{constructVal.SliceFloat([]float64{math.NaN()}), kinds.Bool, false, true},
		{constructVal.SliceFloat([]float64{1.5}), kinds.Bool, true, false},

		{constructVal.SliceFloat([]float64{math.NaN()}), kinds.DateTime, time.Time{}, true},
		{constructVal.SliceFloat([]float64{float64(testEpoch)}), kinds.DateTime, testDate, false},

		// Int
		{constructVal.SliceInt([]int64{1}), kinds.Float, 1.0, false},
		{constructVal.SliceInt([]int64{1}), kinds.Int, int64(1), false},
		{constructVal.SliceInt([]int64{1}), kinds.String, "1", false},
		{constructVal.SliceInt([]int64{1}), kinds.Bool, true, false},

		{constructVal.SliceInt([]int64{1}), kinds.DateTime, epochDate, false},
		{constructVal.SliceInt([]int64{int64(testEpoch)}), kinds.DateTime, testDate, false},

		// String
		{constructVal.SliceString([]string{""}), kinds.Float, nan, true},
		{constructVal.SliceString([]string{"foo"}), kinds.Float, nan, true},
		{constructVal.SliceString([]string{"1.5"}), kinds.Float, 1.5, false},

		{constructVal.SliceString([]string{""}), kinds.Int, int64(0), true},
		{constructVal.SliceString([]string{"foo"}), kinds.Int, int64(0), true},
		{constructVal.SliceString([]string{"1.5"}), kinds.Int, int64(1), false},
		{constructVal.SliceString([]string{"1.0"}), kinds.Int, int64(1), false},
		{constructVal.SliceString([]string{"1"}), kinds.Int, int64(1), false},

		{constructVal.SliceString([]string{""}), kinds.String, "NaN", true},
		{constructVal.SliceString([]string{"NaN"}), kinds.String, "NaN", true},
		{constructVal.SliceString([]string{"n/a"}), kinds.String, "NaN", true},
		{constructVal.SliceString([]string{"N/A"}), kinds.String, "NaN", true},
		{constructVal.SliceString([]string{"1.5"}), kinds.String, "1.5", false},
		{constructVal.SliceString([]string{"foo"}), kinds.String, "foo", false},

		{constructVal.SliceString([]string{""}), kinds.Bool, false, true},
		{constructVal.SliceString([]string{"foo"}), kinds.Bool, true, false},

		{constructVal.SliceString([]string{"May 1, 2019"}), kinds.DateTime, testDate, false},
		{constructVal.SliceString([]string{"5/1/2019"}), kinds.DateTime, testDate, false},
		{constructVal.SliceString([]string{"2019-05-01"}), kinds.DateTime, testDate, false},

		// Bool
		{constructVal.SliceBool([]bool{true}), kinds.Float, 1.0, false},
		{constructVal.SliceBool([]bool{false}), kinds.Float, 0.0, false},

		{constructVal.SliceBool([]bool{true}), kinds.Int, int64(1), false},
		{constructVal.SliceBool([]bool{false}), kinds.Int, int64(0), false},

		{constructVal.SliceBool([]bool{true}), kinds.String, "true", false},
		{constructVal.SliceBool([]bool{false}), kinds.String, "false", false},

		{constructVal.SliceBool([]bool{true}), kinds.Bool, true, false},
		{constructVal.SliceBool([]bool{false}), kinds.Bool, false, false},

		{constructVal.SliceBool([]bool{true}), kinds.DateTime, epochDate, false},
		{constructVal.SliceBool([]bool{false}), kinds.DateTime, epochDate, false},

		// DateTime
		{constructVal.SliceDateTime([]time.Time{testDate}), kinds.Float, float64(testEpoch), false},
		{constructVal.SliceDateTime([]time.Time{time.Time{}}), kinds.Float, nan, true},

		{constructVal.SliceDateTime([]time.Time{testDate}), kinds.Int, int64(testEpoch), false},
		{constructVal.SliceDateTime([]time.Time{time.Time{}}), kinds.Int, int64(0), true},

		{constructVal.SliceDateTime([]time.Time{testDate}), kinds.String, "2019-05-01 00:00:00 +0000 UTC", false},
		{constructVal.SliceDateTime([]time.Time{time.Time{}}), kinds.String, "NaN", true},

		{constructVal.SliceDateTime([]time.Time{testDate}), kinds.Bool, true, false},
		{constructVal.SliceDateTime([]time.Time{time.Time{}}), kinds.Bool, false, true},

		{constructVal.SliceDateTime([]time.Time{testDate}), kinds.DateTime, testDate, false},
		{constructVal.SliceDateTime([]time.Time{time.Time{}}), kinds.DateTime, time.Time{}, true},

		// Interface
		{constructVal.SliceInterface([]interface{}{math.NaN()}), kinds.Float, nan, true},
		{constructVal.SliceInterface([]interface{}{1.5}), kinds.Float, 1.5, false},

		{constructVal.SliceInterface([]interface{}{1}), kinds.Int, int64(1), false},

		{constructVal.SliceInterface([]interface{}{""}), kinds.String, "NaN", true},
		{constructVal.SliceInterface([]interface{}{"NaN"}), kinds.String, "NaN", true},
		{constructVal.SliceInterface([]interface{}{"n/a"}), kinds.String, "NaN", true},
		{constructVal.SliceInterface([]interface{}{"N/A"}), kinds.String, "NaN", true},
		{constructVal.SliceInterface([]interface{}{"1.5"}), kinds.String, "1.5", false},
		{constructVal.SliceInterface([]interface{}{"foo"}), kinds.String, "foo", false},

		{constructVal.SliceInterface([]interface{}{true}), kinds.Bool, true, false},
		{constructVal.SliceInterface([]interface{}{false}), kinds.Bool, false, false},

		{constructVal.SliceInterface([]interface{}{testDate}), kinds.DateTime, testDate, false},
		{constructVal.SliceInterface([]interface{}{time.Time{}}), kinds.DateTime, time.Time{}, true},
	}

	for _, test := range tests {
		converted, err := values.Convert(test.input, test.convertTo)
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
		{kinds.Invalid},
		{kinds.Unsupported},
	}
	for _, test := range tests {
		vals := constructVal.SliceFloat([]float64{1.5})
		_, err := values.Convert(vals, test.kind)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.kind)
		}
	}
}
