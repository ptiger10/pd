package constructors

import (
	"time"

	"github.com/ptiger10/pd/internal/values"
)

// [START Constructor Functions]

// SliceDateTime converts []time.Time{} -> values.DateTimeValues
func SliceDateTime(vals []time.Time) values.DateTimeValues {
	var v values.DateTimeValues
	if len(vals) == 0 {
		return nil
	}
	for _, val := range vals {
		if (time.Time{}) == val {
			v = append(v, values.DateTime(val, true))
		} else {
			v = append(v, values.DateTime(val, false))
		}

	}
	return v
}

// [END Constructor Functions]
