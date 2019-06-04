// This file was automatically generated.
// Any changes will be lost if this file is regenerated.
// Run "make generate" to regenerate from template.

package values

import (
	"fmt"
	"time"

	"github.com/ptiger10/pd/opt"
)

// [START] float64Values

// float64Values is a slice of float64-typed value/null structs
type float64Values []float64Value

// float64Value is a float64-typed value/null struct
type float64Value struct {
	v    float64
	null bool
}

// float64Val constructs a float64Value value/null struct
func float64Val(v float64, null bool) float64Value {
	return float64Value{
		v:    v,
		null: null,
	}
}

// Len returns the number of value/null structs in the container
func (vals float64Values) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals float64Values) In(rowPositions []int) (Values, error) {
	var ret float64Values
	for _, position := range rowPositions {
		if position >= len(vals) {
			return nil, fmt.Errorf("%d is not a valid integer position (len: %v)", position, len(vals))
		}
		ret = append(ret, vals[position])
	}
	return ret, nil
}

// All returns only the Value fields for the collection of Value/Null structs as an interface slice.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals float64Values) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals float64Values) Vals() interface{} {
	var ret []float64
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals float64Values) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals float64Values) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a value element at an integer position in form
// []interface{val, null}
func (vals float64Values) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// ToString converts the values to stringValues
func (vals float64Values) ToString() Values {
	var ret stringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, stringVal(opt.GetDisplayStringNullFiller(), true))
		} else {
			ret = append(ret, stringVal(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// ToInterface converts the values to interfaceValues
func (vals float64Values) ToInterface() Values {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceVal(val.v, true))
		} else {
			ret = append(ret, interfaceVal(val.v, false))
		}
	}
	return ret
}

// [END] float64Values
// ---------------------------------------------------------------------------


// [START] int64Values

// int64Values is a slice of int64-typed value/null structs
type int64Values []int64Value

// int64Value is a int64-typed value/null struct
type int64Value struct {
	v    int64
	null bool
}

// int64Val constructs a int64Value value/null struct
func int64Val(v int64, null bool) int64Value {
	return int64Value{
		v:    v,
		null: null,
	}
}

// Len returns the number of value/null structs in the container
func (vals int64Values) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals int64Values) In(rowPositions []int) (Values, error) {
	var ret int64Values
	for _, position := range rowPositions {
		if position >= len(vals) {
			return nil, fmt.Errorf("%d is not a valid integer position (len: %v)", position, len(vals))
		}
		ret = append(ret, vals[position])
	}
	return ret, nil
}

// All returns only the Value fields for the collection of Value/Null structs as an interface slice.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals int64Values) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals int64Values) Vals() interface{} {
	var ret []int64
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals int64Values) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals int64Values) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a value element at an integer position in form
// []interface{val, null}
func (vals int64Values) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// ToString converts the values to stringValues
func (vals int64Values) ToString() Values {
	var ret stringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, stringVal(opt.GetDisplayStringNullFiller(), true))
		} else {
			ret = append(ret, stringVal(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// ToInterface converts the values to interfaceValues
func (vals int64Values) ToInterface() Values {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceVal(val.v, true))
		} else {
			ret = append(ret, interfaceVal(val.v, false))
		}
	}
	return ret
}

// [END] int64Values
// ---------------------------------------------------------------------------


// [START] stringValues

// stringValues is a slice of string-typed value/null structs
type stringValues []stringValue

// stringValue is a string-typed value/null struct
type stringValue struct {
	v    string
	null bool
}

// stringVal constructs a stringValue value/null struct
func stringVal(v string, null bool) stringValue {
	return stringValue{
		v:    v,
		null: null,
	}
}

// Len returns the number of value/null structs in the container
func (vals stringValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals stringValues) In(rowPositions []int) (Values, error) {
	var ret stringValues
	for _, position := range rowPositions {
		if position >= len(vals) {
			return nil, fmt.Errorf("%d is not a valid integer position (len: %v)", position, len(vals))
		}
		ret = append(ret, vals[position])
	}
	return ret, nil
}

// All returns only the Value fields for the collection of Value/Null structs as an interface slice.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals stringValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals stringValues) Vals() interface{} {
	var ret []string
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals stringValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals stringValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a value element at an integer position in form
// []interface{val, null}
func (vals stringValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// ToString converts the values to stringValues
func (vals stringValues) ToString() Values {
	var ret stringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, stringVal(opt.GetDisplayStringNullFiller(), true))
		} else {
			ret = append(ret, stringVal(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// ToInterface converts the values to interfaceValues
func (vals stringValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceVal(val.v, true))
		} else {
			ret = append(ret, interfaceVal(val.v, false))
		}
	}
	return ret
}

// [END] stringValues
// ---------------------------------------------------------------------------


// [START] boolValues

// boolValues is a slice of bool-typed value/null structs
type boolValues []boolValue

// boolValue is a bool-typed value/null struct
type boolValue struct {
	v    bool
	null bool
}

// boolVal constructs a boolValue value/null struct
func boolVal(v bool, null bool) boolValue {
	return boolValue{
		v:    v,
		null: null,
	}
}

// Len returns the number of value/null structs in the container
func (vals boolValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals boolValues) In(rowPositions []int) (Values, error) {
	var ret boolValues
	for _, position := range rowPositions {
		if position >= len(vals) {
			return nil, fmt.Errorf("%d is not a valid integer position (len: %v)", position, len(vals))
		}
		ret = append(ret, vals[position])
	}
	return ret, nil
}

// All returns only the Value fields for the collection of Value/Null structs as an interface slice.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals boolValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals boolValues) Vals() interface{} {
	var ret []bool
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals boolValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals boolValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a value element at an integer position in form
// []interface{val, null}
func (vals boolValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// ToString converts the values to stringValues
func (vals boolValues) ToString() Values {
	var ret stringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, stringVal(opt.GetDisplayStringNullFiller(), true))
		} else {
			ret = append(ret, stringVal(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// ToInterface converts the values to interfaceValues
func (vals boolValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceVal(val.v, true))
		} else {
			ret = append(ret, interfaceVal(val.v, false))
		}
	}
	return ret
}

// [END] boolValues
// ---------------------------------------------------------------------------


// [START] dateTimeValues

// dateTimeValues is a slice of dateTime-typed value/null structs
type dateTimeValues []dateTimeValue

// dateTimeValue is a dateTime-typed value/null struct
type dateTimeValue struct {
	v    time.Time
	null bool
}

// dateTimeVal constructs a dateTimeValue value/null struct
func dateTimeVal(v time.Time, null bool) dateTimeValue {
	return dateTimeValue{
		v:    v,
		null: null,
	}
}

// Len returns the number of value/null structs in the container
func (vals dateTimeValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals dateTimeValues) In(rowPositions []int) (Values, error) {
	var ret dateTimeValues
	for _, position := range rowPositions {
		if position >= len(vals) {
			return nil, fmt.Errorf("%d is not a valid integer position (len: %v)", position, len(vals))
		}
		ret = append(ret, vals[position])
	}
	return ret, nil
}

// All returns only the Value fields for the collection of Value/Null structs as an interface slice.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals dateTimeValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals dateTimeValues) Vals() interface{} {
	var ret []time.Time
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals dateTimeValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals dateTimeValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a value element at an integer position in form
// []interface{val, null}
func (vals dateTimeValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// ToString converts the values to stringValues
func (vals dateTimeValues) ToString() Values {
	var ret stringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, stringVal(opt.GetDisplayStringNullFiller(), true))
		} else {
			ret = append(ret, stringVal(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// ToInterface converts the values to interfaceValues
func (vals dateTimeValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceVal(val.v, true))
		} else {
			ret = append(ret, interfaceVal(val.v, false))
		}
	}
	return ret
}

// [END] dateTimeValues
// ---------------------------------------------------------------------------


// [START] interfaceValues

// interfaceValues is a slice of interface-typed value/null structs
type interfaceValues []interfaceValue

// interfaceValue is a interface-typed value/null struct
type interfaceValue struct {
	v    interface{}
	null bool
}

// interfaceVal constructs a interfaceValue value/null struct
func interfaceVal(v interface{}, null bool) interfaceValue {
	return interfaceValue{
		v:    v,
		null: null,
	}
}

// Len returns the number of value/null structs in the container
func (vals interfaceValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals interfaceValues) In(rowPositions []int) (Values, error) {
	var ret interfaceValues
	for _, position := range rowPositions {
		if position >= len(vals) {
			return nil, fmt.Errorf("%d is not a valid integer position (len: %v)", position, len(vals))
		}
		ret = append(ret, vals[position])
	}
	return ret, nil
}

// All returns only the Value fields for the collection of Value/Null structs as an interface slice.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals interfaceValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals interfaceValues) Vals() interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals interfaceValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals interfaceValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a value element at an integer position in form
// []interface{val, null}
func (vals interfaceValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// ToString converts the values to stringValues
func (vals interfaceValues) ToString() Values {
	var ret stringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, stringVal(opt.GetDisplayStringNullFiller(), true))
		} else {
			ret = append(ret, stringVal(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// ToInterface converts the values to interfaceValues
func (vals interfaceValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceVal(val.v, true))
		} else {
			ret = append(ret, interfaceVal(val.v, false))
		}
	}
	return ret
}

// [END] interfaceValues
// ---------------------------------------------------------------------------

