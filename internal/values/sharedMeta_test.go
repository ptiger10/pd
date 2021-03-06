package values

import (
	"reflect"
	"testing"
)

// Pro forma tests for generics
func TestMeta(t *testing.T) {
	newSlicevalueType([]valueType{newvalueType("foo")})

	val := newvalueType("foo")
	f := val.toFloat64()
	if vType := reflect.TypeOf(f); vType.Name() != "float64Value" {
		t.Errorf("%v", vType.Name())
	}
	i := val.toInt64()
	if vType := reflect.TypeOf(i); vType.Name() != "int64Value" {
		t.Errorf("%v", vType.Name())
	}
	s := val.toString()
	if vType := reflect.TypeOf(s); vType.Name() != "stringValue" {
		t.Errorf("%v", vType.Name())
	}
	b := val.toBool()
	if vType := reflect.TypeOf(b); vType.Name() != "boolValue" {
		t.Errorf("%v", vType.Name())
	}
	dt := val.toDateTime()
	if vType := reflect.TypeOf(dt); vType.Name() != "dateTimeValue" {
		t.Errorf("%v", vType.Name())
	}

	nullVal := valueTypeValue{"foo", true}
	nullVal.toString()
	nullVals := valueTypeValues{nullVal}
	nullVals.ToInterface()

	vals := valueTypeValues{val}
	vals.Len()
	vals.Swap(0, 0)
	vals.Less(0, 0)
	vals.Values()
	vals.Vals()
	vals.Copy()
	vals.Value(0)
	vals.Null(0)
	vals.ToFloat64()
	vals.ToInt64()
	vals.ToString()
	vals.ToBool()
	vals.ToDateTime()
	vals.ToInterface()

	vals.Subset([]int{0})
	vals.Set(0, "bar")
	vals.Set(0, "")
	vals.Drop(0)
	vals.Insert(0, "foo")

	v := interfaceValue{"foo", false}
	v.tovalueType()
}

// No easy way to Convert valueTypeValues, so expect panic
func TestPanic(t *testing.T) {
	val := newvalueType("foo")
	vals := valueTypeValues{val}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// The following is the code under test
	vals.Append(&vals)
}
