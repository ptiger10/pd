package constructors

import (
	"time"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// [START Constructor Functions]

// SliceDateTime converts []time.Time{} -> ValuesFactory with values.DateTimeValues
func SliceDateTime(vals []time.Time) ValuesFactory {
	var v values.DateTimeValues
	if len(vals) == 0 {
		return ValuesFactory{nil, kinds.DateTime}
	}
	for _, val := range vals {
		if (time.Time{}) == val {
			v = append(v, values.DateTime(val, true))
		} else {
			v = append(v, values.DateTime(val, false))
		}

	}
	return ValuesFactory{v, kinds.DateTime}
}

// [END Constructor Functions]
