package values

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ptiger10/pd/options"
)

// TestSharedFloat64 tests shared float64 interface methods
func TestSharedFloat64(t *testing.T) {
	factory := newSliceFloat64([]float64{1, 2})
	wantContainer := Container{&float64Values{float64Value{1, false}, float64Value{2, false}}, options.Float64}
	if !reflect.DeepEqual(factory, wantContainer) {
		t.Errorf("newSlice got %v, want %v", factory, wantContainer)
	}

	vals := factory.Values
	gotLen := vals.Len()
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Len() got %v, want %v", gotLen, wantLen)
	}
	gotVals := vals.Vals()
	wantVals := []float64{1, 2}
	if !reflect.DeepEqual(gotVals, wantVals) {
		t.Errorf("Vals() got %v, want %v", gotVals, wantVals)
	}
	gotValues := vals.Values()
	wantValues := []interface{}{1.0, 2.0}
	if !reflect.DeepEqual(gotValues, wantValues) {
		t.Errorf("Values() got %v, want %v", gotValues, wantValues)
	}

	e := vals.Element(0)
	wantElem := Elem{1.0, false}
	if !reflect.DeepEqual(e, wantElem) {
		t.Errorf("Element() got %v, want %v", e, wantElem)
	}
	gotLess := vals.Less(0, 1)
	wantLess := true
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}
	gotLess = vals.Less(1, 0)
	wantLess = false
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}

	v := vals.Copy()
	if reflect.ValueOf(v).Pointer() == reflect.ValueOf(vals).Pointer() {
		t.Errorf("Copy() retained reference to the original, want copy")
	}

	subset := vals.Subset([]int{0})
	wantSubset := &float64Values{float64Value{1, false}}
	if !reflect.DeepEqual(subset, wantSubset) {
		t.Errorf("Subset() got %v, want %v", subset, wantSubset)
	}

	vals.Swap(0, 1)
	wantSwap := &float64Values{float64Value{2, false}, float64Value{1, false}}
	if !reflect.DeepEqual(vals, wantSwap) {
		t.Errorf("Swap() got %v, want %v", vals, wantSwap)
	}

	vals.Set(0, 5)
	wantSet := &float64Values{float64Value{5, false}, float64Value{1, false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() got %v, want %v", vals, wantSet)
	}

	vals.Set(0, "")
	if elem := vals.Element(0); !math.IsNaN(elem.Value.(float64)) || elem.Null != true {
		t.Errorf("Set() on null value did not return null values")
	}

	vals.Drop(0)
	wantDrop := &float64Values{float64Value{1, false}}
	if !reflect.DeepEqual(vals, wantDrop) {
		t.Errorf("Drop() got %v, want %v", vals, wantDrop)
	}

	vals.Insert(1, 2)
	wantInsert := &float64Values{float64Value{1, false}, float64Value{2, false}}
	if !reflect.DeepEqual(vals, wantInsert) {
		t.Errorf("Insert() got %v, want %v", vals, wantInsert)
	}

}

// TestSharedInt64 tests shared int64 interface methods
func TestSharedInt64(t *testing.T) {
	factory := newSliceInt64([]int64{1, 2})
	wantContainer := Container{&int64Values{int64Value{1, false}, int64Value{2, false}}, options.Int64}
	if !reflect.DeepEqual(factory, wantContainer) {
		t.Errorf("newSlice got %v, want %v", factory, wantContainer)
	}

	vals := factory.Values
	gotLen := vals.Len()
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Len() got %v, want %v", gotLen, wantLen)
	}
	gotVals := vals.Vals()
	wantVals := []int64{1, 2}
	if !reflect.DeepEqual(gotVals, wantVals) {
		t.Errorf("Vals() got %v, want %v", gotVals, wantVals)
	}
	gotValues := vals.Values()
	wantValues := []interface{}{int64(1), int64(2)}
	if !reflect.DeepEqual(gotValues, wantValues) {
		t.Errorf("Values() got %v, want %v", gotValues, wantValues)
	}

	e := vals.Element(0)
	wantElem := Elem{int64(1), false}
	if !reflect.DeepEqual(e, wantElem) {
		t.Errorf("Element() got %v, want %v", e, wantElem)
	}
	gotLess := vals.Less(0, 1)
	wantLess := true
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}
	gotLess = vals.Less(1, 0)
	wantLess = false
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}

	v := vals.Copy()
	if reflect.ValueOf(v).Pointer() == reflect.ValueOf(vals).Pointer() {
		t.Errorf("Copy() retained reference to the original, want copy")
	}

	subset := vals.Subset([]int{0})
	wantSubset := &int64Values{int64Value{1, false}}
	if !reflect.DeepEqual(subset, wantSubset) {
		t.Errorf("Subset() got %v, want %v", subset, wantSubset)
	}

	vals.Swap(0, 1)
	wantSwap := &int64Values{int64Value{2, false}, int64Value{1, false}}
	if !reflect.DeepEqual(vals, wantSwap) {
		t.Errorf("Swap() got %v, want %v", vals, wantSwap)
	}

	vals.Set(0, 5)
	wantSet := &int64Values{int64Value{5, false}, int64Value{1, false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() got %v, want %v", vals, wantSet)
	}

	vals.Set(0, "")
	wantSet = &int64Values{int64Value{0, true}, int64Value{1, false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() on null value got %v, want %v", vals, wantSet)
	}

	vals.Drop(0)
	wantDrop := &int64Values{int64Value{1, false}}
	if !reflect.DeepEqual(vals, wantDrop) {
		t.Errorf("Drop() got %v, want %v", vals, wantDrop)
	}

	vals.Insert(1, 2)
	wantInsert := &int64Values{int64Value{1, false}, int64Value{2, false}}
	if !reflect.DeepEqual(vals, wantInsert) {
		t.Errorf("Insert() got %v, want %v", vals, wantInsert)
	}
}

// TestSharedString tests shared string interface methods
func TestSharedString(t *testing.T) {
	factory := newSliceString([]string{"bar", "foo"})
	wantContainer := Container{&stringValues{stringValue{"bar", false}, stringValue{"foo", false}}, options.String}
	if !reflect.DeepEqual(factory, wantContainer) {
		t.Errorf("newSlice got %v, want %v", factory, wantContainer)
	}

	vals := factory.Values
	gotLen := vals.Len()
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Len() got %v, want %v", gotLen, wantLen)
	}
	gotVals := vals.Vals()
	wantVals := []string{"bar", "foo"}
	if !reflect.DeepEqual(gotVals, wantVals) {
		t.Errorf("Vals() got %v, want %v", gotVals, wantVals)
	}
	gotValues := vals.Values()
	wantValues := []interface{}{"bar", "foo"}
	if !reflect.DeepEqual(gotValues, wantValues) {
		t.Errorf("Values() got %v, want %v", gotValues, wantValues)
	}

	e := vals.Element(0)
	wantElem := Elem{"bar", false}
	if !reflect.DeepEqual(e, wantElem) {
		t.Errorf("Element() got %v, want %v", e, wantElem)
	}
	gotLess := vals.Less(0, 1)
	wantLess := true
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}
	gotLess = vals.Less(1, 0)
	wantLess = false
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}

	v := vals.Copy()
	if reflect.ValueOf(v).Pointer() == reflect.ValueOf(vals).Pointer() {
		t.Errorf("Copy() retained reference to the original, want copy")
	}

	subset := vals.Subset([]int{0})
	wantSubset := &stringValues{stringValue{"bar", false}}
	if !reflect.DeepEqual(subset, wantSubset) {
		t.Errorf("Subset() got %v, want %v", subset, wantSubset)
	}

	vals.Swap(0, 1)
	wantSwap := &stringValues{stringValue{"foo", false}, stringValue{"bar", false}}
	if !reflect.DeepEqual(vals, wantSwap) {
		t.Errorf("Swap() got %v, want %v", vals, wantSwap)
	}

	vals.Set(0, "baz")
	wantSet := &stringValues{stringValue{"baz", false}, stringValue{"bar", false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() got %v, want %v", vals, wantSet)
	}

	vals.Set(0, "")
	wantSet = &stringValues{stringValue{"NaN", true}, stringValue{"bar", false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() on null value got %v, want %v", vals, wantSet)
	}

	vals.Drop(0)
	wantDrop := &stringValues{stringValue{"bar", false}}
	if !reflect.DeepEqual(vals, wantDrop) {
		t.Errorf("Drop() got %v, want %v", vals, wantDrop)
	}

	vals.Insert(1, "foo")
	wantInsert := &stringValues{stringValue{"bar", false}, stringValue{"foo", false}}
	if !reflect.DeepEqual(vals, wantInsert) {
		t.Errorf("Insert() got %v, want %v", vals, wantInsert)
	}
}

// TestSharedBool tests shared bool interface methods
func TestSharedBool(t *testing.T) {
	factory := newSliceBool([]bool{false, true})
	wantContainer := Container{&boolValues{boolValue{false, false}, boolValue{true, false}}, options.Bool}
	if !reflect.DeepEqual(factory, wantContainer) {
		t.Errorf("newSlice got %v, want %v", factory, wantContainer)
	}

	vals := factory.Values
	gotLen := vals.Len()
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Len() got %v, want %v", gotLen, wantLen)
	}
	gotVals := vals.Vals()
	wantVals := []bool{false, true}
	if !reflect.DeepEqual(gotVals, wantVals) {
		t.Errorf("Vals() got %v, want %v", gotVals, wantVals)
	}
	gotValues := vals.Values()
	wantValues := []interface{}{false, true}
	if !reflect.DeepEqual(gotValues, wantValues) {
		t.Errorf("Values() got %v, want %v", gotValues, wantValues)
	}

	e := vals.Element(0)
	wantElem := Elem{false, false}
	if !reflect.DeepEqual(e, wantElem) {
		t.Errorf("Element() got %v, want %v", e, wantElem)
	}
	gotLess := vals.Less(0, 1)
	wantLess := true
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}
	gotLess = vals.Less(1, 0)
	wantLess = false
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}

	v := vals.Copy()
	if reflect.ValueOf(v).Pointer() == reflect.ValueOf(vals).Pointer() {
		t.Errorf("Copy() retained reference to the original, want copy")
	}

	subset := vals.Subset([]int{0})
	wantSubset := &boolValues{boolValue{false, false}}
	if !reflect.DeepEqual(subset, wantSubset) {
		t.Errorf("Subset() got %v, want %v", subset, wantSubset)
	}

	vals.Swap(0, 1)
	wantSwap := &boolValues{boolValue{true, false}, boolValue{false, false}}
	if !reflect.DeepEqual(vals, wantSwap) {
		t.Errorf("Swap() got %v, want %v", vals, wantSwap)
	}

	vals.Set(0, false)
	wantSet := &boolValues{boolValue{false, false}, boolValue{false, false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() got %v, want %v", vals, wantSet)
	}

	vals.Set(0, "")
	wantSet = &boolValues{boolValue{false, true}, boolValue{false, false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() on null value got %v, want %v", vals, wantSet)
	}

	vals.Drop(0)
	wantDrop := &boolValues{boolValue{false, false}}
	if !reflect.DeepEqual(vals, wantDrop) {
		t.Errorf("Drop() got %v, want %v", vals, wantDrop)
	}

	vals.Insert(1, true)
	wantInsert := &boolValues{boolValue{false, false}, boolValue{true, false}}
	if !reflect.DeepEqual(vals, wantInsert) {
		t.Errorf("Insert() got %v, want %v", vals, wantInsert)
	}
}

// TestSharedDateTime tests shared time.Time interface methods
func TestSharedDateTime(t *testing.T) {
	dt1 := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	dt2 := time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)
	factory := newSliceDateTime([]time.Time{dt1, dt2})
	wantContainer := Container{&dateTimeValues{dateTimeValue{dt1, false}, dateTimeValue{dt2, false}}, options.DateTime}
	if !reflect.DeepEqual(factory, wantContainer) {
		t.Errorf("newSlice got %v, want %v", factory, wantContainer)
	}

	vals := factory.Values
	gotLen := vals.Len()
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Len() got %v, want %v", gotLen, wantLen)
	}
	gotVals := vals.Vals()
	wantVals := []time.Time{dt1, dt2}
	if !reflect.DeepEqual(gotVals, wantVals) {
		t.Errorf("Vals() got %v, want %v", gotVals, wantVals)
	}
	gotValues := vals.Values()
	wantValues := []interface{}{dt1, dt2}
	if !reflect.DeepEqual(gotValues, wantValues) {
		t.Errorf("Values() got %v, want %v", gotValues, wantValues)
	}

	e := vals.Element(0)
	wantElem := Elem{dt1, false}
	if !reflect.DeepEqual(e, wantElem) {
		t.Errorf("Element() got %v, want %v", e, wantElem)
	}
	gotLess := vals.Less(0, 1)
	wantLess := true
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}
	gotLess = vals.Less(1, 0)
	wantLess = false
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}

	v := vals.Copy()
	if reflect.ValueOf(v).Pointer() == reflect.ValueOf(vals).Pointer() {
		t.Errorf("Copy() retained reference to the original, want copy")
	}

	subset := vals.Subset([]int{0})
	wantSubset := &dateTimeValues{dateTimeValue{dt1, false}}
	if !reflect.DeepEqual(subset, wantSubset) {
		t.Errorf("Subset() got %v, want %v", subset, wantSubset)
	}

	vals.Swap(0, 1)
	wantSwap := &dateTimeValues{dateTimeValue{dt2, false}, dateTimeValue{dt1, false}}
	if !reflect.DeepEqual(vals, wantSwap) {
		t.Errorf("Swap() got %v, want %v", vals, wantSwap)
	}

	vals.Set(0, dt1)
	wantSet := &dateTimeValues{dateTimeValue{dt1, false}, dateTimeValue{dt1, false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() on null value got %v, want %v", vals, wantSet)
	}

	vals.Set(0, "")
	wantSet = &dateTimeValues{dateTimeValue{time.Time{}, true}, dateTimeValue{dt1, false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() on null values got %v, want %v", vals, wantSet)
	}

	vals.Drop(0)
	wantDrop := &dateTimeValues{dateTimeValue{dt1, false}}
	if !reflect.DeepEqual(vals, wantDrop) {
		t.Errorf("Drop() got %v, want %v", vals, wantDrop)
	}

	vals.Insert(1, dt2)
	wantInsert := &dateTimeValues{dateTimeValue{dt1, false}, dateTimeValue{dt2, false}}
	if !reflect.DeepEqual(vals, wantInsert) {
		t.Errorf("Insert() got %v, want %v", vals, wantInsert)
	}
}

// TestSharedInterface tests shared interface{} interface methods
func TestSharedInterface(t *testing.T) {
	factory := newSliceInterface([]interface{}{false, true})
	wantContainer := Container{&interfaceValues{interfaceValue{false, false}, interfaceValue{true, false}}, options.Interface}
	if !reflect.DeepEqual(factory, wantContainer) {
		t.Errorf("newSlice got %v, want %v", factory, wantContainer)
	}

	vals := factory.Values
	gotLen := vals.Len()
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Len() got %v, want %v", gotLen, wantLen)
	}
	gotVals := vals.Vals()
	wantVals := []interface{}{false, true}
	if !reflect.DeepEqual(gotVals, wantVals) {
		t.Errorf("Vals() got %v, want %v", gotVals, wantVals)
	}
	gotValues := vals.Values()
	wantValues := []interface{}{false, true}
	if !reflect.DeepEqual(gotValues, wantValues) {
		t.Errorf("Values() got %v, want %v", gotValues, wantValues)
	}

	e := vals.Element(0)
	wantElem := Elem{false, false}
	if !reflect.DeepEqual(e, wantElem) {
		t.Errorf("Element() got %v, want %v", e, wantElem)
	}
	gotLess := vals.Less(0, 1)
	wantLess := true
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}
	gotLess = vals.Less(1, 0)
	wantLess = false
	if gotLess != wantLess {
		t.Errorf("Less() got %v, want %v", gotLess, wantLess)
	}

	v := vals.Copy()
	if reflect.ValueOf(v).Pointer() == reflect.ValueOf(vals).Pointer() {
		t.Errorf("Copy() retained reference to the original, want copy")
	}

	subset := vals.Subset([]int{0})
	wantSubset := &interfaceValues{interfaceValue{false, false}}
	if !reflect.DeepEqual(subset, wantSubset) {
		t.Errorf("Subset() got %v, want %v", subset, wantSubset)
	}

	vals.Swap(0, 1)
	wantSwap := &interfaceValues{interfaceValue{true, false}, interfaceValue{false, false}}
	if !reflect.DeepEqual(vals, wantSwap) {
		t.Errorf("Swap() got %v, want %v", vals, wantSwap)
	}

	vals.Set(0, false)
	wantSet := &interfaceValues{interfaceValue{false, false}, interfaceValue{false, false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() got %v, want %v", vals, wantSet)
	}

	vals.Set(0, "")
	wantSet = &interfaceValues{interfaceValue{"", true}, interfaceValue{false, false}}
	if !reflect.DeepEqual(vals, wantSet) {
		t.Errorf("Set() on null values got %v, want %v", vals, wantSet)
	}

	vals.Drop(0)
	wantDrop := &interfaceValues{interfaceValue{false, false}}
	if !reflect.DeepEqual(vals, wantDrop) {
		t.Errorf("Drop() got %v, want %v", vals, wantDrop)
	}

	vals.Insert(1, true)
	wantInsert := &interfaceValues{interfaceValue{false, false}, interfaceValue{true, false}}
	if !reflect.DeepEqual(vals, wantInsert) {
		t.Errorf("Insert() got %v, want %v", vals, wantInsert)
	}

}

// [START conversion tests]

func TestConvert(t *testing.T) {
	nullFloat := Container{&float64Values{float64Value{math.NaN(), true}}, options.Float64}
	nullInt := Container{&int64Values{int64Value{0, true}}, options.Int64}
	nullBool := Container{&boolValues{boolValue{false, true}}, options.Bool}
	nullDateTime := Container{&dateTimeValues{dateTimeValue{time.Time{}, true}}, options.DateTime}
	nullInterface := Container{&interfaceValues{interfaceValue{0, true}}, options.Interface}

	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	testDateInSeconds := 1556668800000000000
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	nan := math.NaN()
	var tests = []struct {
		name      string
		input     Container
		convertTo options.DataType
		wantVal   interface{}
		wantNull  bool
	}{
		// Float
		{"float(null)->float", nullFloat, options.Float64, nan, true},
		{"float->float", newSliceFloat64([]float64{1.5}), options.Float64, 1.5, false},

		{"float(null)->int", nullFloat, options.Int64, int64(0), true},
		{"float->int", newSliceFloat64([]float64{1.5}), options.Int64, int64(1), false},

		{"float(null)->string", nullFloat, options.String, "NaN", true},
		{"float->string", newSliceFloat64([]float64{1.5}), options.String, "1.5", false},

		{"float(null)->bool", nullFloat, options.Bool, false, true},
		{"float->bool0", newSliceFloat64([]float64{0}), options.Bool, false, false},
		{"float->bool1", newSliceFloat64([]float64{1.5}), options.Bool, true, false},

		{"float(null)->datetime", nullFloat, options.DateTime, time.Time{}, true},
		{"float->datetime", newSliceFloat64([]float64{float64(testDateInSeconds)}), options.DateTime, testDate, false},

		{"float(null)->float(null)interface", nullFloat, options.Interface, nan, true},
		{"float->interface", newSliceFloat64([]float64{1.5}), options.Interface, 1.5, false},

		// Int
		{"int(null)->float", nullInt, options.Float64, nan, true},
		{"int->float", newSliceInt64([]int64{1}), options.Float64, 1.0, false},

		{"int(null)->int", nullInt, options.Int64, int64(0), true},
		{"int->int", newSliceInt64([]int64{1}), options.Int64, int64(1), false},

		{"int(null)->string", nullInt, options.String, "NaN", true},
		{"int->string", newSliceInt64([]int64{1}), options.String, "1", false},

		{"int(null)->bool", nullInt, options.Bool, false, true},
		{"int->bool", newSliceInt64([]int64{0}), options.Bool, false, false},
		{"int->bool", newSliceInt64([]int64{1}), options.Bool, true, false},

		{"int(null)->datetime", nullInt, options.DateTime, time.Time{}, true},
		{"int->datetime", newSliceInt64([]int64{1}), options.DateTime, epochDate, false},
		{"int->datetime", newSliceInt64([]int64{int64(testDateInSeconds)}), options.DateTime, testDate, false},

		{"int(null)->interface", nullInt, options.Interface, int64(0), true},
		{"int->interface", newSliceInt64([]int64{1}), options.Interface, int64(1), false},

		// String
		{"string(null)->float", newSliceString([]string{""}), options.Float64, nan, true},
		{"string(null)->float", newSliceString([]string{"foo"}), options.Float64, nan, true},
		{"string->float", newSliceString([]string{"1.5"}), options.Float64, 1.5, false},

		{"string(null)->int", newSliceString([]string{""}), options.Int64, int64(0), true},
		{"string(null)->int", newSliceString([]string{"foo"}), options.Int64, int64(0), true},
		{"string(null)->int", newSliceString([]string{"1.5"}), options.Int64, int64(1), false},
		{"string(null)->int", newSliceString([]string{"1.0"}), options.Int64, int64(1), false},
		{"string(null)->int", newSliceString([]string{"1"}), options.Int64, int64(1), false},

		{"string(null)->string", newSliceString([]string{""}), options.String, "NaN", true},
		{"string(null)->string", newSliceString([]string{"NaN"}), options.String, "NaN", true},
		{"string(null)->string", newSliceString([]string{"n/a"}), options.String, "NaN", true},
		{"string(null)->string", newSliceString([]string{"N/A"}), options.String, "NaN", true},
		{"string->string", newSliceString([]string{"1.5"}), options.String, "1.5", false},
		{"string->string", newSliceString([]string{"foo"}), options.String, "foo", false},

		{"string(null)->bool", newSliceString([]string{""}), options.Bool, false, true},
		{"string->bool", newSliceString([]string{"foo"}), options.Bool, true, false},

		{"string(null)->datetime", newSliceString([]string{""}), options.DateTime, time.Time{}, true},
		{"string(null)->datetime", newSliceString([]string{"1 of May in year 2019"}), options.DateTime, time.Time{}, true},
		{"string->datetime", newSliceString([]string{"May 1, 2019"}), options.DateTime, testDate, false},
		{"string->datetime", newSliceString([]string{"5/1/2019"}), options.DateTime, testDate, false},
		{"string->datetime", newSliceString([]string{"2019-05-01"}), options.DateTime, testDate, false},

		{"string(null)->interface", newSliceString([]string{""}), options.Interface, "NaN", true},
		{"string->interface", newSliceString([]string{"foo"}), options.Interface, "foo", false},

		// Bool
		{"bool(null)->float", nullBool, options.Float64, nan, true},
		{"bool->float", newSliceBool([]bool{true}), options.Float64, 1.0, false},
		{"bool->float", newSliceBool([]bool{false}), options.Float64, 0.0, false},

		{"bool(null)->int", nullBool, options.Int64, int64(0), true},
		{"bool->int", newSliceBool([]bool{true}), options.Int64, int64(1), false},
		{"bool->int", newSliceBool([]bool{false}), options.Int64, int64(0), false},

		{"bool(null)->string", nullBool, options.String, "NaN", true},
		{"bool->string", newSliceBool([]bool{true}), options.String, "true", false},
		{"bool->string", newSliceBool([]bool{false}), options.String, "false", false},

		{"bool(null)->bool", nullBool, options.Bool, false, true},
		{"bool->bool", newSliceBool([]bool{true}), options.Bool, true, false},
		{"bool->bool", newSliceBool([]bool{false}), options.Bool, false, false},

		{"bool(null)->datetime", nullBool, options.DateTime, time.Time{}, true},
		{"bool->datetime", newSliceBool([]bool{true}), options.DateTime, epochDate, false},
		{"bool->datetime", newSliceBool([]bool{false}), options.DateTime, epochDate, false},

		{"bool(null)->interface", nullBool, options.Interface, false, true},
		{"bool->interface", newSliceBool([]bool{true}), options.Interface, true, false},

		// DateTime
		{"datetime(null)->float", nullDateTime, options.Float64, nan, true},
		{"datetime->float", newSliceDateTime([]time.Time{testDate}), options.Float64, float64(testDateInSeconds), false},

		{"datetime(null)->int", nullDateTime, options.Int64, int64(0), true},
		{"datetime->int", newSliceDateTime([]time.Time{testDate}), options.Int64, int64(testDateInSeconds), false},

		{"datetime(null)->string", nullDateTime, options.String, "NaN", true},
		{"datetime->string", newSliceDateTime([]time.Time{testDate}), options.String, "2019-05-01 00:00:00 +0000 UTC", false},

		{"datetime(null)->bool", nullDateTime, options.Bool, false, true},
		{"datetime->bool", newSliceDateTime([]time.Time{testDate}), options.Bool, true, false},

		{"datetime(null)->datetime", nullDateTime, options.DateTime, time.Time{}, true},
		{"datetime->datetime", newSliceDateTime([]time.Time{testDate}), options.DateTime, testDate, false},

		{"datetime(null)->interface", nullDateTime, options.Interface, time.Time{}, true},
		{"datetime->interface", newSliceDateTime([]time.Time{testDate}), options.Interface, testDate, false},

		// Interface
		{"interface(null)->float", nullInterface, options.Float64, nan, true},
		{"interface(null)->float", newSliceInterface([]interface{}{"foo"}), options.Float64, nan, true},
		{"interface(null)->float", newSliceInterface([]interface{}{complex64(1)}), options.Float64, nan, true},
		{"interfaceFloat->float", newSliceInterface([]interface{}{1.5}), options.Float64, 1.5, false},
		{"interfaceInt->float", newSliceInterface([]interface{}{1}), options.Float64, 1.0, false},
		{"interfaceUInt->float", newSliceInterface([]interface{}{uint(1)}), options.Float64, 1.0, false},
		{"interfaceString->float", newSliceInterface([]interface{}{"1"}), options.Float64, 1.0, false},
		{"interfaceBool->float", newSliceInterface([]interface{}{true}), options.Float64, 1.0, false},
		{"interfaceDateTime->float", newSliceInterface([]interface{}{testDate}), options.Float64, float64(testDateInSeconds), false},

		{"interface(null)->int", nullInterface, options.Int64, int64(0), true},
		{"interface(null)->int", newSliceInterface([]interface{}{complex64(1)}), options.Int64, int64(0), true},
		{"interfaceFloat->int", newSliceInterface([]interface{}{1.0}), options.Int64, int64(1), false},
		{"interfaceInt->int", newSliceInterface([]interface{}{1}), options.Int64, int64(1), false},
		{"interfaceUInt->int", newSliceInterface([]interface{}{uint(1)}), options.Int64, int64(1), false},
		{"interfaceString->int", newSliceInterface([]interface{}{"1"}), options.Int64, int64(1), false},
		{"interfaceBool->int", newSliceInterface([]interface{}{true}), options.Int64, int64(1), false},
		{"interfaceDateTime->int", newSliceInterface([]interface{}{testDate}), options.Int64, int64(testDateInSeconds), false},

		{"interface(null)->string", newSliceInterface([]interface{}{""}), options.String, "NaN", true},
		{"interface(null)->string", newSliceInterface([]interface{}{"NaN"}), options.String, "NaN", true},
		{"interface(null)->string", newSliceInterface([]interface{}{"n/a"}), options.String, "NaN", true},
		{"interface(null)->string", newSliceInterface([]interface{}{"N/A"}), options.String, "NaN", true},
		{"interfaceFloat->string", newSliceInterface([]interface{}{1.5}), options.String, "1.5", false},
		{"interfaceInt->string", newSliceInterface([]interface{}{1}), options.String, "1", false},
		{"interfaceUInt->string", newSliceInterface([]interface{}{uint(1)}), options.String, "1", false},
		{"interfaceString->string", newSliceInterface([]interface{}{"1.5"}), options.String, "1.5", false},
		{"interfaceBool->string", newSliceInterface([]interface{}{true}), options.String, "true", false},
		{"interfaceDateTime->string", newSliceInterface([]interface{}{testDate}), options.String, "2019-05-01 00:00:00 +0000 UTC", false},

		{"interface(null)->bool", nullInterface, options.Bool, false, true},
		{"interface(null)->bool", newSliceInterface([]interface{}{complex64(1)}), options.Bool, false, true},
		{"interfaceFloat->bool", newSliceInterface([]interface{}{1.5}), options.Bool, true, false},
		{"interfaceInt->bool", newSliceInterface([]interface{}{1}), options.Bool, true, false},
		{"interfaceUInt->bool", newSliceInterface([]interface{}{uint(1)}), options.Bool, true, false},
		{"interfaceString->bool", newSliceInterface([]interface{}{"1.5"}), options.Bool, true, false},
		{"interfaceBool->bool", newSliceInterface([]interface{}{true}), options.Bool, true, false},
		{"interfaceDateTime->bool", newSliceInterface([]interface{}{testDate}), options.Bool, true, false},

		{"interface(null)->datetime", nullInterface, options.DateTime, time.Time{}, true},
		{"interface(null)->datetime", newSliceInterface([]interface{}{complex64(1)}), options.DateTime, time.Time{}, true},
		{"interfaceFloat->datetime", newSliceInterface([]interface{}{float64(testDateInSeconds)}), options.DateTime, testDate, false},
		{"interfaceInt->datetime", newSliceInterface([]interface{}{testDateInSeconds}), options.DateTime, testDate, false},
		{"interfaceUInt->datetime", newSliceInterface([]interface{}{uint(testDateInSeconds)}), options.DateTime, testDate, false},
		{"interfaceString->datetime", newSliceInterface([]interface{}{"2019-05-01"}), options.DateTime, testDate, false},
		{"interfaceBool->datetime", newSliceInterface([]interface{}{true}), options.DateTime, epochDate, false},
		{"interfaceDateTime->datetime", newSliceInterface([]interface{}{time.Time{}}), options.DateTime, time.Time{}, true},

		{"interface(null)->interface", nullInterface, options.Interface, 0, true},
		{"interfaceFloat->interface", newSliceInterface([]interface{}{1.5}), options.Interface, 1.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converted, err := Convert(tt.input.Values, tt.convertTo)
			if err != nil {
				t.Errorf("Convert(): %v", err)
			}
			elem := converted.Element(0)
			val := elem.Value
			null := elem.Null
			if strings.Contains(tt.name, "(null)->float") {
				if !math.IsNaN(val.(float64)) {
					t.Errorf("Convert(): got %v, want NaN", val)
				}
			} else if val != tt.wantVal {
				t.Errorf("Convert() returned elem.Value %v, want %v", val, tt.wantVal)
			}
			if null != tt.wantNull {
				t.Errorf("Convert() returned elem.Null %v, want %v", null, tt.wantNull)
			}
		})
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

// [END conversion tests]
