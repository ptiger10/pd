package constructors

import (
	"reflect"

	"github.com/ptiger10/pd/new/internal/values"
)

// [START Utilities]

// [END Utilities]

// [START Constructor Functions]

// SliceInt converts []int (of any variety) -> values.IntValues
func SliceInt(data interface{}) values.IntValues {
	var vals values.IntValues
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		val := d.Index(i).Int()
		vals = append(vals, values.Int(val, false))
	}
	return vals
}

// SliceUint converts []uint (of any variety) -> values.IntValues
func SliceUint(data interface{}) values.IntValues {
	var vals values.IntValues
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		val := d.Index(i).Uint()
		vals = append(vals, values.Int(int64(val), false))
	}
	return vals
}

// [END Constructor Functions]
