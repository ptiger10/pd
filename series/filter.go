package series

import (
	"fmt"
	"time"
)

// Float
func (vals floatValues) Filter(callbackFn func(float64) bool) floatValues {
	var ret floatValues
	valid, _ := vals.valid()
	for _, val := range valid {
		if callbackFn(val) {
			ret = append(ret, floatValue{v: val})
		}
	}
	return ret
}

func (s Series) FilterFloat(callbackFn func(float64) bool) (Series, error) {
	if s.Kind != Float {
		return s, fmt.Errorf("FilterFloat can be called only on Series with type Float, not %v", s.Kind)
	}
	vals := s.Values.(floatValues).Filter(callbackFn)

	return Series{
		Values: vals,
		Kind:   Float,
	}, nil
}

// Int
func (vals intValues) Filter(callbackFn func(float64) bool) intValues {
	var ret intValues
	valid, _ := vals.valid()
	for _, val := range valid {
		if callbackFn(float64(val)) {
			ret = append(ret, intValue{v: val})
		}
	}
	return ret
}

func (s Series) FilterInt(callbackFn func(float64) bool) (Series, error) {
	if s.Kind != Int {
		return s, fmt.Errorf("FilterInt can be called only on Series with type Int, not %v", s.Kind)
	}
	vals := s.Values.(intValues).Filter(callbackFn)

	return Series{
		Values: vals,
		Kind:   Int,
	}, nil
}

// String
func (vals stringValues) Filter(callbackFn func(string) bool) stringValues {
	var ret stringValues
	valid, _ := vals.valid()
	for _, val := range valid {
		if callbackFn(val) {
			ret = append(ret, stringValue{v: val})
		}
	}
	return ret
}

func (s Series) FilterString(callbackFn func(string) bool) (Series, error) {
	if s.Kind != String {
		return s, fmt.Errorf("FilterString can be called only on Series with type String, not %v", s.Kind)
	}
	vals := s.Values.(stringValues).Filter(callbackFn)

	return Series{
		Values: vals,
		Kind:   String,
	}, nil
}

// Datetime
func (vals dateTimeValues) Filter(callbackFn func(time.Time) bool) dateTimeValues {
	var ret dateTimeValues
	valid, _ := vals.valid()
	for _, val := range valid {
		if callbackFn(val) {
			ret = append(ret, dateTimeValue{v: val})
		}
	}
	return ret
}

func (s Series) FilterDateTime(callbackFn func(time.Time) bool) (Series, error) {
	if s.Kind != DateTime {
		return s, fmt.Errorf("FilterString can be called only on Series with type String, not %v", s.Kind)
	}
	vals := s.Values.(dateTimeValues).Filter(callbackFn)

	return Series{
		Values: vals,
		Kind:   DateTime,
	}, nil
}
