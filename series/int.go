package series

import (
	"math"
	"reflect"
	"strconv"
)

// Data Types
// ------------------------------------------------
type intValues []intValue
type intValue struct {
	v    int64
	null bool
}

// Convenience Functions
// ------------------------------------------------
func (vals intValues) toFloat() floatValues {
	var ret floatValues
	for _, val := range vals {
		ret = append(ret, floatValue{null: val.null, v: float64(val.v)})
	}
	return ret
}

func (vals intValues) valid() ([]int64, []int) {
	var valid []int64
	var nullMap []int
	for i, val := range vals {
		if !val.null {
			valid = append(valid, int64(val.v))
		} else {
			nullMap = append(nullMap, i)
		}
	}
	return valid, nullMap
}

// Methods
// ------------------------------------------------
func (vals intValues) count() int {
	floatVals := vals.toFloat()
	return floatVals.count()
}

func (vals intValues) sum() float64 {
	floatVals := vals.toFloat()
	return floatVals.sum()
}

func (vals intValues) mean() float64 {
	floatVals := vals.toFloat()
	return floatVals.mean()
}

func (vals intValues) median() float64 {
	floatVals := vals.toFloat()
	return floatVals.median()
}

func (vals intValues) min() float64 {
	floatVals := vals.toFloat()
	return floatVals.min()
}

func (vals intValues) max() float64 {
	floatVals := vals.toFloat()
	return floatVals.max()
}

func (vals intValues) describe() string {
	floatVals := vals.toFloat()
	return floatVals.describe()
}

// Constructor Functions
// ------------------------------------------------
func intToIntValues(data interface{}) intValues {
	var vals []intValue
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		vals = append(vals, intValue{v: d.Index(i).Int()})
	}
	return intValues(vals)
}

func uIntToIntValues(data interface{}) intValues {
	var vals []intValue
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		val := int64(d.Index(i).Uint())
		vals = append(vals, intValue{v: val})
	}
	return intValues(vals)
}

func interfaceToIntValues(d reflect.Value) intValues {
	var vals []intValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, intValue{null: true})
		case reflect.String:
			val, err := strconv.ParseFloat(v.String(), 64)
			if err != nil || math.IsNaN(val) {
				vals = append(vals, intValue{null: true})
			} else {
				vals = append(vals, intValue{v: int64(val)})
			}
		case reflect.Float32, reflect.Float64:
			val := v.Float()
			if math.IsNaN(val) {
				vals = append(vals, intValue{null: true})
			} else {
				vals = append(vals, intValue{v: int64(val)})
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, intValue{v: int64(v.Int())})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, intValue{v: int64(v.Uint())})
		default:
			vals = append(vals, intValue{null: true})
		}
	}
	return intValues(vals)
}
