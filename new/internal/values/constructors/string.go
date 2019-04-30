package constructors

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/ptiger10/pd/new/internal/values"
)

// **[START Utilities]**

func isNullString(s string) bool {
	nullStrings := []string{"nan", "n/a", ""}
	for _, ns := range nullStrings {
		if strings.TrimSpace(strings.ToLower(s)) == ns {
			return true
		}
	}
	return false
}

// **[END Utilities]**

// Constructor Functions

// SliceString converts []string -> StringValues
func SliceString(data interface{}) values.StringValues {
	var vals []values.StringValue
	d := data.([]string)
	for i := 0; i < len(d); i++ {
		val := d[i]
		if isNullString(val) {
			vals = append(vals, values.StringValue{Null: true})
		} else {
			vals = append(vals, values.StringValue{V: val})
		}
	}
	return values.StringValues(vals)
}

// InterfaceToString converts []interface -> StringValues
func InterfaceToString(d reflect.Value) values.StringValues {
	var vals []values.StringValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, values.StringValue{Null: true})
		case reflect.String:
			val := v.String()
			if isNullString(val) {
				vals = append(vals, values.StringValue{Null: true})
			} else {
				vals = append(vals, values.StringValue{V: val})
			}
		case reflect.Float32, reflect.Float64:
			vals = append(vals, floatToString(v.Float()))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, values.StringValue{V: fmt.Sprint(v.Int())})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, values.StringValue{V: fmt.Sprint(v.Uint())})
		default:
			vals = append(vals, values.StringValue{Null: true})
		}
	}
	return values.StringValues(vals)
}

func floatToString(v float64) values.StringValue {
	if math.IsNaN(v) {
		return values.StringValue{Null: true}
	}
	return values.StringValue{V: fmt.Sprint(v)}
}
