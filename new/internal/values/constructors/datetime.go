package constructors

import (
	"time"

	"github.com/ptiger10/pd/new/internal/values"
)

// [START Constructor Functions]

// SliceDateTime converts []time.Time{} -> values.DateTimeValues
func SliceDateTime(data interface{}) values.DateTimeValues {
	var vals values.DateTimeValues
	d := data.([]time.Time)
	for i := 0; i < len(d); i++ {
		val := d[i]
		if (time.Time{}) == val {
			vals = append(vals, values.DateTime(val, true))
		} else {
			vals = append(vals, values.DateTime(val, false))
		}

	}
	return vals
}

// [END Constructor Functions]
