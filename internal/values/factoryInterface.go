package values

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/datatypes"
)

// InterfaceFactory converts interface{} to Values
func InterfaceFactory(data interface{}) (Factory, error) {
	var factory Factory
	var err error
	if data == nil {
		factory = Factory{Values: nil, DataType: datatypes.None}
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
			return Factory{}, fmt.Errorf("unable to create value from interface: type not supported: %T", data)
		}
	}
	return factory, err
}
