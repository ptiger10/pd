package values

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/ptiger10/pd/options"
)

// InterfaceFactory converts interface{} to Values
func InterfaceFactory(data interface{}) (Container, error) {
	var factory Container
	var err error
	if data == nil {
		factory = Container{Values: emptyValues(), DataType: options.None}
	} else {
		switch reflect.ValueOf(data).Kind() {
		case reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.String,
			reflect.Bool,
			reflect.Struct:
			factory, err = ScalarFactory(data)

		case reflect.Slice:
			factory, err = SliceFactory(data)

		default:
			return Container{Values: emptyValues()}, fmt.Errorf("internal.values.InterfaceFactory(): type not supported: %T", data)
		}
	}
	return factory, err
}

// MustCreateValuesFromInterface returns a container that satisfies the Values interface or panics.
func MustCreateValuesFromInterface(data interface{}) Container {
	container, err := InterfaceFactory(data)
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("MustCreateValuesFromInterface(): %v", err)
		}
	}
	return container
}

// ScalarFactory creates a Container from an interface{} that has been determined elsewhere to be a scalar.
// Be careful modifying this constructor and its children,
// as reflection is inherently unsafe and each function expects to receive specific types only.
func ScalarFactory(data interface{}) (Container, error) {
	var ret Container
	switch data.(type) {
	case float32, float64:
		val := scalarFloatToFloat64(data)
		f := newFloat64(val)
		ret = Container{&float64Values{f}, options.Float64}

	case int, int8, int16, int32, int64:
		val := scalarIntToInt64(data)
		i := newInt64(val)
		ret = Container{&int64Values{i}, options.Int64}

	case uint, uint8, uint16, uint32, uint64:
		val := scalarUIntToInt64(data)
		i := newInt64(val)
		ret = Container{&int64Values{i}, options.Int64}

	case string:
		val := newString(data.(string))
		ret = Container{&stringValues{val}, options.String}

	case bool:
		val := newBool(data.(bool))
		ret = Container{&boolValues{val}, options.Bool}

	case time.Time:
		val := newDateTime(data.(time.Time))
		ret = Container{&dateTimeValues{val}, options.DateTime}

	default:
		ret = Container{}
		return ret, fmt.Errorf("Type %T not supported", data)
	}

	return ret, nil
}

// scalarFloatToFloat64 converts a known float interface{} -> float64
func scalarFloatToFloat64(data interface{}) float64 {
	d := reflect.ValueOf(data)
	return d.Float()
}

// scalarIntToInt64 converts a known int interface{} -> int64
func scalarIntToInt64(data interface{}) int64 {
	d := reflect.ValueOf(data)
	return d.Int()
}

// scalarUIntToInt64 converts a known uint interface{} -> int64
func scalarUIntToInt64(data interface{}) int64 {
	d := reflect.ValueOf(data)
	return int64(d.Uint())
}

// SliceFactory creates a Container from an interface{} that has been determined elsewhere to be a Slice.
// Be careful modifying this constructor and its children,
// as reflection is inherently unsafe and each function expects to receive specific types only.
func SliceFactory(data interface{}) (Container, error) {
	var ret Container

	switch data.(type) {
	case []float32:
		vals := sliceFloatToSliceFloat64(data)
		ret = newSliceFloat64(vals)

	case []float64:
		vals := data.([]float64)
		ret = newSliceFloat64(vals)

	case []int, []int8, []int16, []int32:
		vals := sliceIntToSliceInt64(data)
		ret = newSliceInt64(vals)
	case []int64:
		vals := data.([]int64)
		ret = newSliceInt64(vals)

	case []uint, []uint8, []uint16, []uint32, []uint64:
		// converts into []int64 so it can be passed into SliceInt
		vals := sliceUIntToSliceInt64(data)
		ret = newSliceInt64(vals)

	case []string:
		vals := data.([]string)
		ret = newSliceString(vals)

	case []bool:
		vals := data.([]bool)
		ret = newSliceBool(vals)

	case []time.Time:
		vals := data.([]time.Time)
		ret = newSliceDateTime(vals)

	case []interface{}:
		vals := data.([]interface{})
		ret = newSliceInterface(vals)

	default:
		ret = Container{nil, options.None}
		return ret, fmt.Errorf("Type %T not supported", data)
	}

	return ret, nil
}

// [START interface converters]

// sliceFloatToSliceFloat64 converts known []float interface{} -> []float64
func sliceFloatToSliceFloat64(data interface{}) []float64 {
	var vals []float64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Float()
		vals = append(vals, v)
	}
	return vals
}

// sliceIntToSliceInt64 converts known []int interface{} -> []int64
func sliceIntToSliceInt64(data interface{}) []int64 {
	var vals []int64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Int()
		vals = append(vals, v)
	}
	return vals
}

// sliceUIntToSliceInt64 converts knonw []uint interface{} -> []int64
func sliceUIntToSliceInt64(data interface{}) []int64 {
	var vals []int64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Uint()
		vals = append(vals, int64(v))
	}
	return vals
}

// [END interface converters]

// [START utility slices]

// MakeIntRange returns a sequential series of numbers, for use in making default index labels.
// Includes min and excludes max.
func MakeIntRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// MakeIntRangeInclusive returns a sequential series of numbers, for use in making default index labels.
// Includes start and end.
func MakeIntRangeInclusive(start, end int) []int {
	var ret []int
	if start <= end {
		ret = make([]int, (end-start)+1)
		for i := range ret {
			ret[i] = start + i
		}
	} else {
		ret = make([]int, (start-end)+1)
		for i := range ret {
			ret[i] = start - i
		}
	}

	return ret
}

// MakeStringRange returns a sequential series of numbers as string values, for use in making default column labels.
func MakeStringRange(min, max int) []string {
	a := make([]string, max-min)
	for i := range a {
		a[i] = strconv.Itoa(min + i)
	}
	return a
}

// [END utility slices]

// Copy copies a Values Container
func (vc Container) Copy() Container {
	return Container{
		Values:   vc.Values.Copy(),
		DataType: vc.DataType,
	}
}

// AllValues returns all the values (including null values) in the Container as an interface slice.
func (vc Container) AllValues() []interface{} {
	var ret []interface{}
	for i := 0; i < vc.Values.Len(); i++ {
		ret = append(ret, vc.Values.Element(i).Value)
	}
	return ret
}

// MaxWidth returns the max width of any value in the container. For use in printing dataframes.
func (vc Container) MaxWidth() int {
	var max int
	for _, v := range vc.AllValues() {
		var length int
		if vc.DataType == options.DateTime {
			if val, ok := v.(time.Time); ok {
				length = len(val.Format(options.GetDisplayTimeFormat()))
			} else {
				length = len(fmt.Sprint(v))
			}
		} else if vc.DataType == options.Float64 {
			if val, ok := v.(float64); ok {
				length = len(fmt.Sprintf("%.*f", options.GetDisplayFloatPrecision(), val))
			} else {
				length = len(fmt.Sprint(v))
			}
		} else {
			length = len(fmt.Sprint(v))
		}
		if length > max {
			max = length
		}
	}
	return max
}
