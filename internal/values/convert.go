package values

import (
	"fmt"

	"github.com/ptiger10/pd/datatypes"
)

// Convert a collection of values from one type to another, and coerce to null if a value cannot be converted sensibly
func Convert(currentVals Values, kind datatypes.DataType) (Values, error) {
	var vals Values
	switch kind {
	case datatypes.None:
		return nil, fmt.Errorf("unable to convert values: must supply a valid Kind")
	case datatypes.Float64:
		vals = currentVals.ToFloat()
	case datatypes.Int64:
		vals = currentVals.ToInt()
	case datatypes.String:
		vals = currentVals.ToString()
	case datatypes.Bool:
		vals = currentVals.ToBool()
	case datatypes.DateTime:
		vals = currentVals.ToDateTime()
	case datatypes.Interface:
		vals = currentVals.ToInterface()
	default:
		return nil, fmt.Errorf("unable to convert values: kind not supported: %v", kind)
	}
	return vals, nil
}
