// Package values is an internal package that powers the values within pd/Series and pd/DataFrame.
// This package defines the Values interface and multiple concrete implementations of the interface.
package values

import (
	"fmt"

	"github.com/ptiger10/pd/options"
)

// The Values interface is the primary means of handling a collection of values.
// The same interface and value types are used for both Series values and Index labels
type Values interface {
	Len() int                // number of Value/Null structs
	Vals() interface{}       // an interface of values, ready for type assertion into a slice of their native type
	Values() []interface{}   // an interface slice of values, for handling values as a predictable slice
	Subset([]int) Values     // a new Values object comprised of the Value/Null pairs at one or more integer positions
	Value(int) interface{}   // the value field at an integer position
	Null(int) bool           // the null field at an integer position
	Set(int, interface{})    // overwrite the value/null struct at an integer position
	Copy() Values            // clone the Values
	Insert(int, interface{}) // insert a Value/Null pair at an integer position
	Append(Values)           // append Values together
	Drop(int)                // drop a Value/Null pair at an integer position
	Swap(i, j int)           // swap two values - necessary for sorting
	Less(i, j int) bool      // compare two values and return the lesser - required for sorting

	ToFloat64() Values
	ToInt64() Values
	ToString() Values
	ToBool() Values
	ToDateTime() Values
	ToInterface() Values
}

// Container contains Values (a list of Value/Null pairs satisfying the Values interface) and Kind.
type Container struct {
	Values   Values
	DataType options.DataType
}

// Convert a collection of values from one type to another, and coerce to null if a value cannot be converted sensibly
func Convert(currentVals Values, dataType options.DataType) (Values, error) {
	var vals Values
	switch dataType {
	case options.None:
		return nil, fmt.Errorf("unable to convert values: must supply a valid Kind")
	case options.Float64:
		vals = currentVals.ToFloat64()
	case options.Int64:
		vals = currentVals.ToInt64()
	case options.String:
		vals = currentVals.ToString()
	case options.Bool:
		vals = currentVals.ToBool()
	case options.DateTime:
		vals = currentVals.ToDateTime()
	case options.Interface:
		vals = currentVals.ToInterface()
	default:
		return nil, fmt.Errorf("unable to convert values: kind not supported: %v", dataType)
	}
	return vals, nil
}
