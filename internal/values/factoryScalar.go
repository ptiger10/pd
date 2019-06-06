package values

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ptiger10/pd/kinds"
)

// ScalarFactory creates a Factory from an interface{} that has been determined elsewhere to be a scalar.
// Be careful modifying this constructor and its children,
// as reflection is inherently unsafe and each function expects to receive specific types only.
func ScalarFactory(data interface{}) (Factory, error) {
	var ret Factory
	switch data.(type) {
	case float32, float64:
		val := scalarFloatToFloat64(data)
		f := newFloat64(val)
		ret = Factory{&float64Values{f}, kinds.Float64}

	case int, int8, int16, int32, int64:
		val := scalarIntToInt64(data)
		i := newInt64(val)
		ret = Factory{&int64Values{i}, kinds.Int64}

	case uint, uint8, uint16, uint32, uint64:
		val := scalarUIntToInt64(data)
		i := newInt64(val)
		ret = Factory{&int64Values{i}, kinds.Int64}

	case string:
		val := newString(data.(string))
		ret = Factory{&stringValues{val}, kinds.String}

	case bool:
		val := newBool(data.(bool))
		ret = Factory{&boolValues{val}, kinds.Bool}

	case time.Time:
		val := newDateTime(data.(time.Time))
		ret = Factory{&dateTimeValues{val}, kinds.DateTime}

	case interface{}:
		val := newInterface(data.(interface{}))
		ret = Factory{&interfaceValues{val}, kinds.Interface}

	default:
		ret = Factory{}
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
