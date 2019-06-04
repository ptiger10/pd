package values

import (
	"fmt"
	"reflect"
	"time"
)

// ScalarFactory creates an Elem from an interface{} that has been determined elsewhere to be a scalar.
// Be careful modifying this constructor and its children,
// as reflection is inherently unsafe and each function expects to receive specific types only.
func ScalarFactory(data interface{}) (Elem, error) {
	var ret Elem
	switch data.(type) {
	case float32, float64:
		val := scalarFloatToFloat64(data)
		f := newFloat(val)
		ret = Elem{f.v, f.null}

	case int, int8, int16, int32, int64:
		val := scalarIntToInt64(data)
		i := newInt(val)
		ret = Elem{i.v, i.null}

	case uint, uint8, uint16, uint32, uint64:
		// converts into int64 so it can be passed into SliceInt
		val := scalarUIntToInt64(data)
		i := newInt(val)
		ret = Elem{i.v, i.null}

	case string:
		val := newString(data.(string))
		ret = Elem{val.v, val.null}

	case bool:
		val := newBool(data.(bool))
		ret = Elem{val.v, val.null}

	case time.Time:
		val := newDateTime(data.(time.Time))
		ret = Elem{val.v, val.null}

	case interface{}:
		val := newInterface(data.(interface{}))
		ret = Elem{val.v, val.null}

	default:
		ret = Elem{}
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
