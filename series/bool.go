package series

import (
	"math"
	"reflect"
)

// data type
// ------------------------------------------------
type boolValues []boolValue
type boolValue struct {
	v    bool
	null bool
}

// methods
// ------------------------------------------------
func (vals boolValues) sum() float64 {
	var sum float64
	for _, val := range vals {
		if val.v {
			sum++
		}
	}
	return sum
}

// constructor functions
// ------------------------------------------------
func boolToBoolValues(data interface{}) boolValues {
	var vals []boolValue
	d := data.([]bool)
	for i := 0; i < len(d); i++ {
		vals = append(vals, boolValue{v: d[i]})
	}
	return boolValues(vals)
}

func interfaceToBoolValues(d reflect.Value) boolValues {
	var vals []boolValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, boolValue{null: true})
		case reflect.String:
			if isNullString(v.String()) {
				vals = append(vals, boolValue{null: true})
			} else {
				vals = append(vals, boolValue{v: true})
			}
		case reflect.Float32, reflect.Float64:
			val := v.Float()
			if math.IsNaN(val) {
				vals = append(vals, boolValue{null: true})
			} else {
				vals = append(vals, boolValue{v: true})
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, boolValue{v: true})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, boolValue{v: true})
		default:
			vals = append(vals, boolValue{null: true})
		}
	}
	return boolValues(vals)
}
