package constructors

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/new/internal/values"
	"github.com/ptiger10/pd/new/kinds"
)

// Convert converts a collection of values from one type to another, if possible
func Convert(currentVals values.Values, kind reflect.Kind) (values.Values, error) {
	var vals values.Values
	switch kind {
	case kinds.None: // this checks for the pseduo-nil type
		return nil, fmt.Errorf("Unable to convert values: must supply a valid Kind")
	case kinds.Float:
		vals = currentVals.ToFloat()
	case kinds.Int:
		vals = currentVals.ToInt()
	case kinds.String:
		vals = currentVals.ToString()
	case kinds.Bool:
		vals = currentVals.ToBool()
	case kinds.DateTime:
		vals = currentVals.ToDateTime()
	default:
		return nil, fmt.Errorf("Unable to convert values: kind not supported: %v", kind)
	}
	return vals, nil
}
