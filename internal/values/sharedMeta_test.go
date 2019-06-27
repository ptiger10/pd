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
	vals.Vals()
	vals.Element(0)
	vals.Copy()
	vals.ToFloat64()
	vals.ToInt64()
	vals.ToString()
	vals.ToBool()
	vals.ToDateTime()
	vals.ToInterface()

	_, err := vals.Subset([]int{0})
	if err != nil {
		t.Errorf("valueTypeValues.Subset() %v", err)
	}
	err = vals.Set(0, "bar")
	if err != nil {
		t.Errorf("valueTypeValues.Set() %v", err)
	}
	err = vals.Drop(0)
	if err != nil {
		t.Errorf("valueTypeValues.Drop() %v", err)
	}
	err = vals.Insert(0, "foo")
	if err != nil {
		t.Errorf("valueTypeValues.Insert() %v", err)
	}

	// Should return errors
	_, err = vals.Subset([]int{1})
	if err == nil {
		t.Errorf("valueTypeValues.Subset() error = nil, want err")
	}
	err = vals.Set(1, "bar")
	if err == nil {
		t.Errorf("valueTypeValues.Set() %v", err)
	}
	err = vals.Drop(1)
	if err == nil {
		t.Errorf("valueTypeValues.Drop() %v", err)
	}
	err = vals.Insert(2, "bar")
	if err == nil {
		t.Errorf("valueTypeValues.Insert() %v", err)
	}

	v := interfaceValue{"foo", false}
	v.tovalueType()
}
