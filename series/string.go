package series

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

// data type
// ------------------------------------------------
type stringValues []stringValue
type stringValue struct {
	v    string
	null bool
}

// convenience methods
// ------------------------------------------------
func isNullString(s string) bool {
	nullStrings := []string{"nan", "n/a", ""}
	for _, ns := range nullStrings {
		if strings.TrimSpace(strings.ToLower(s)) == ns {
			return true
		}
	}
	return false
}

// constructor functions
// ------------------------------------------------
func stringToStringValues(data interface{}) stringValues {
	var vals []stringValue
	d := data.([]string)
	for i := 0; i < len(d); i++ {
		val := d[i]
		if isNullString(val) {
			vals = append(vals, stringValue{null: true})
		} else {
			vals = append(vals, stringValue{v: val})
		}
	}
	return stringValues(vals)
}

func interfaceToStringValues(d reflect.Value) stringValues {
	var vals []stringValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, stringValue{null: true})
		case reflect.String:
			val := v.String()
			if isNullString(val) {
				vals = append(vals, stringValue{null: true})
			} else {
				vals = append(vals, stringValue{v: val})
			}
		case reflect.Float32, reflect.Float64:
			vals = append(vals, floatToString(v.Float()))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, stringValue{v: fmt.Sprint(v.Int())})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, stringValue{v: fmt.Sprint(v.Uint())})
		default:
			vals = append(vals, stringValue{null: true})
		}
	}
	return stringValues(vals)
}

func floatToString(v float64) stringValue {
	if math.IsNaN(v) {
		return stringValue{null: true}
	}
	return stringValue{v: fmt.Sprint(v)}
}
