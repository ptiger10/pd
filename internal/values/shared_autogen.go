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

// float64Values is a slice of float64-typed value/null structs.
type float64Values []float64Value

// float64Value is a float64-typed value/null struct.
type float64Value struct {
	v    float64
	null bool
}

// Len returns the number of value/null structs in the container.
func (vals float64Values) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions.
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

// Valid returns the integer position of all valid (i.e., non-null) values in the collection.
func (vals float64Values) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection.
func (vals float64Values) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals float64Values) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// Set overwrites a Value/Null pair at an integer position.
func (vals float64Values) Set(position int, val interface{}) error {
	v := interfaceValue{val, false}
	if position >= vals.Len() {
		return fmt.Errorf("unable to set value at position %v: index out of range", position)
	}
	vals[position] = v.toFloat64()
	return nil
}

// ToFloat converts float64Values to floatValues.
func (vals float64Values) ToFloat() float64Values {
	var ret float64Values
	for _, val := range vals {
		ret = append(ret, val.toFloat64())
	}
	return ret
}

// ToInt converts float64Values to intValues.
func (vals float64Values) ToInt() int64Values {
	var ret int64Values
	for _, val := range vals {
		ret = append(ret, val.toInt64())
	}
	return ret
}

func (val float64Value) toString() stringValue {
	if val.null {
		return stringValue{opt.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts float64Values to stringValues.
func (vals float64Values) ToString() stringValues {
	var ret stringValues
	for _, val := range vals {
		ret = append(ret, val.toString())
	}
	return ret
}

// ToBool converts float64Values to boolValues.
func (vals float64Values) ToBool() boolValues {
	var ret boolValues
	for _, val := range vals {
		ret = append(ret, val.toBool())
	}
	return ret
}

// ToBool converts float64Values to dateTimeValues.
func (vals float64Values) ToDateTime() dateTimeValues {
	var ret dateTimeValues
	for _, val := range vals {
		ret = append(ret, val.toDateTime())
	}
	return ret
}

// ToInterface converts float64Values to interfaceValues.
func (vals float64Values) ToInterface() interfaceValues {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return ret
}

// [END] float64Values
// ---------------------------------------------------------------------------


// [START] int64Values

// int64Values is a slice of int64-typed value/null structs.
type int64Values []int64Value

// int64Value is a int64-typed value/null struct.
type int64Value struct {
	v    int64
	null bool
}

// Len returns the number of value/null structs in the container.
func (vals int64Values) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions.
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

// Valid returns the integer position of all valid (i.e., non-null) values in the collection.
func (vals int64Values) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection.
func (vals int64Values) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals int64Values) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// Set overwrites a Value/Null pair at an integer position.
func (vals int64Values) Set(position int, val interface{}) error {
	v := interfaceValue{val, false}
	if position >= vals.Len() {
		return fmt.Errorf("unable to set value at position %v: index out of range", position)
	}
	vals[position] = v.toInt64()
	return nil
}

// ToFloat converts int64Values to floatValues.
func (vals int64Values) ToFloat() float64Values {
	var ret float64Values
	for _, val := range vals {
		ret = append(ret, val.toFloat64())
	}
	return ret
}

// ToInt converts int64Values to intValues.
func (vals int64Values) ToInt() int64Values {
	var ret int64Values
	for _, val := range vals {
		ret = append(ret, val.toInt64())
	}
	return ret
}

func (val int64Value) toString() stringValue {
	if val.null {
		return stringValue{opt.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts int64Values to stringValues.
func (vals int64Values) ToString() stringValues {
	var ret stringValues
	for _, val := range vals {
		ret = append(ret, val.toString())
	}
	return ret
}

// ToBool converts int64Values to boolValues.
func (vals int64Values) ToBool() boolValues {
	var ret boolValues
	for _, val := range vals {
		ret = append(ret, val.toBool())
	}
	return ret
}

// ToBool converts int64Values to dateTimeValues.
func (vals int64Values) ToDateTime() dateTimeValues {
	var ret dateTimeValues
	for _, val := range vals {
		ret = append(ret, val.toDateTime())
	}
	return ret
}

// ToInterface converts int64Values to interfaceValues.
func (vals int64Values) ToInterface() interfaceValues {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return ret
}

// [END] int64Values
// ---------------------------------------------------------------------------


// [START] stringValues

// stringValues is a slice of string-typed value/null structs.
type stringValues []stringValue

// stringValue is a string-typed value/null struct.
type stringValue struct {
	v    string
	null bool
}

// Len returns the number of value/null structs in the container.
func (vals stringValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions.
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

// Valid returns the integer position of all valid (i.e., non-null) values in the collection.
func (vals stringValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection.
func (vals stringValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals stringValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// Set overwrites a Value/Null pair at an integer position.
func (vals stringValues) Set(position int, val interface{}) error {
	v := interfaceValue{val, false}
	if position >= vals.Len() {
		return fmt.Errorf("unable to set value at position %v: index out of range", position)
	}
	vals[position] = v.toString()
	return nil
}

// ToFloat converts stringValues to floatValues.
func (vals stringValues) ToFloat() float64Values {
	var ret float64Values
	for _, val := range vals {
		ret = append(ret, val.toFloat64())
	}
	return ret
}

// ToInt converts stringValues to intValues.
func (vals stringValues) ToInt() int64Values {
	var ret int64Values
	for _, val := range vals {
		ret = append(ret, val.toInt64())
	}
	return ret
}

func (val stringValue) toString() stringValue {
	if val.null {
		return stringValue{opt.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts stringValues to stringValues.
func (vals stringValues) ToString() stringValues {
	var ret stringValues
	for _, val := range vals {
		ret = append(ret, val.toString())
	}
	return ret
}

// ToBool converts stringValues to boolValues.
func (vals stringValues) ToBool() boolValues {
	var ret boolValues
	for _, val := range vals {
		ret = append(ret, val.toBool())
	}
	return ret
}

// ToBool converts stringValues to dateTimeValues.
func (vals stringValues) ToDateTime() dateTimeValues {
	var ret dateTimeValues
	for _, val := range vals {
		ret = append(ret, val.toDateTime())
	}
	return ret
}

// ToInterface converts stringValues to interfaceValues.
func (vals stringValues) ToInterface() interfaceValues {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return ret
}

// [END] stringValues
// ---------------------------------------------------------------------------


// [START] boolValues

// boolValues is a slice of bool-typed value/null structs.
type boolValues []boolValue

// boolValue is a bool-typed value/null struct.
type boolValue struct {
	v    bool
	null bool
}

// Len returns the number of value/null structs in the container.
func (vals boolValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions.
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

// Valid returns the integer position of all valid (i.e., non-null) values in the collection.
func (vals boolValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection.
func (vals boolValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals boolValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// Set overwrites a Value/Null pair at an integer position.
func (vals boolValues) Set(position int, val interface{}) error {
	v := interfaceValue{val, false}
	if position >= vals.Len() {
		return fmt.Errorf("unable to set value at position %v: index out of range", position)
	}
	vals[position] = v.toBool()
	return nil
}

// ToFloat converts boolValues to floatValues.
func (vals boolValues) ToFloat() float64Values {
	var ret float64Values
	for _, val := range vals {
		ret = append(ret, val.toFloat64())
	}
	return ret
}

// ToInt converts boolValues to intValues.
func (vals boolValues) ToInt() int64Values {
	var ret int64Values
	for _, val := range vals {
		ret = append(ret, val.toInt64())
	}
	return ret
}

func (val boolValue) toString() stringValue {
	if val.null {
		return stringValue{opt.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts boolValues to stringValues.
func (vals boolValues) ToString() stringValues {
	var ret stringValues
	for _, val := range vals {
		ret = append(ret, val.toString())
	}
	return ret
}

// ToBool converts boolValues to boolValues.
func (vals boolValues) ToBool() boolValues {
	var ret boolValues
	for _, val := range vals {
		ret = append(ret, val.toBool())
	}
	return ret
}

// ToBool converts boolValues to dateTimeValues.
func (vals boolValues) ToDateTime() dateTimeValues {
	var ret dateTimeValues
	for _, val := range vals {
		ret = append(ret, val.toDateTime())
	}
	return ret
}

// ToInterface converts boolValues to interfaceValues.
func (vals boolValues) ToInterface() interfaceValues {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return ret
}

// [END] boolValues
// ---------------------------------------------------------------------------


// [START] dateTimeValues

// dateTimeValues is a slice of dateTime-typed value/null structs.
type dateTimeValues []dateTimeValue

// dateTimeValue is a dateTime-typed value/null struct.
type dateTimeValue struct {
	v    time.Time
	null bool
}

// Len returns the number of value/null structs in the container.
func (vals dateTimeValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions.
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

// Valid returns the integer position of all valid (i.e., non-null) values in the collection.
func (vals dateTimeValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection.
func (vals dateTimeValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals dateTimeValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// Set overwrites a Value/Null pair at an integer position.
func (vals dateTimeValues) Set(position int, val interface{}) error {
	v := interfaceValue{val, false}
	if position >= vals.Len() {
		return fmt.Errorf("unable to set value at position %v: index out of range", position)
	}
	vals[position] = v.toDateTime()
	return nil
}

// ToFloat converts dateTimeValues to floatValues.
func (vals dateTimeValues) ToFloat() float64Values {
	var ret float64Values
	for _, val := range vals {
		ret = append(ret, val.toFloat64())
	}
	return ret
}

// ToInt converts dateTimeValues to intValues.
func (vals dateTimeValues) ToInt() int64Values {
	var ret int64Values
	for _, val := range vals {
		ret = append(ret, val.toInt64())
	}
	return ret
}

func (val dateTimeValue) toString() stringValue {
	if val.null {
		return stringValue{opt.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts dateTimeValues to stringValues.
func (vals dateTimeValues) ToString() stringValues {
	var ret stringValues
	for _, val := range vals {
		ret = append(ret, val.toString())
	}
	return ret
}

// ToBool converts dateTimeValues to boolValues.
func (vals dateTimeValues) ToBool() boolValues {
	var ret boolValues
	for _, val := range vals {
		ret = append(ret, val.toBool())
	}
	return ret
}

// ToBool converts dateTimeValues to dateTimeValues.
func (vals dateTimeValues) ToDateTime() dateTimeValues {
	var ret dateTimeValues
	for _, val := range vals {
		ret = append(ret, val.toDateTime())
	}
	return ret
}

// ToInterface converts dateTimeValues to interfaceValues.
func (vals dateTimeValues) ToInterface() interfaceValues {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return ret
}

// [END] dateTimeValues
// ---------------------------------------------------------------------------


// [START] interfaceValues

// interfaceValues is a slice of interface-typed value/null structs.
type interfaceValues []interfaceValue

// interfaceValue is a interface-typed value/null struct.
type interfaceValue struct {
	v    interface{}
	null bool
}

// Len returns the number of value/null structs in the container.
func (vals interfaceValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions.
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

// Valid returns the integer position of all valid (i.e., non-null) values in the collection.
func (vals interfaceValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection.
func (vals interfaceValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals interfaceValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// Set overwrites a Value/Null pair at an integer position.
func (vals interfaceValues) Set(position int, val interface{}) error {
	v := interfaceValue{val, false}
	if position >= vals.Len() {
		return fmt.Errorf("unable to set value at position %v: index out of range", position)
	}
	vals[position] = v.toInterface()
	return nil
}

// ToFloat converts interfaceValues to floatValues.
func (vals interfaceValues) ToFloat() float64Values {
	var ret float64Values
	for _, val := range vals {
		ret = append(ret, val.toFloat64())
	}
	return ret
}

// ToInt converts interfaceValues to intValues.
func (vals interfaceValues) ToInt() int64Values {
	var ret int64Values
	for _, val := range vals {
		ret = append(ret, val.toInt64())
	}
	return ret
}

func (val interfaceValue) toString() stringValue {
	if val.null {
		return stringValue{opt.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts interfaceValues to stringValues.
func (vals interfaceValues) ToString() stringValues {
	var ret stringValues
	for _, val := range vals {
		ret = append(ret, val.toString())
	}
	return ret
}

// ToBool converts interfaceValues to boolValues.
func (vals interfaceValues) ToBool() boolValues {
	var ret boolValues
	for _, val := range vals {
		ret = append(ret, val.toBool())
	}
	return ret
}

// ToBool converts interfaceValues to dateTimeValues.
func (vals interfaceValues) ToDateTime() dateTimeValues {
	var ret dateTimeValues
	for _, val := range vals {
		ret = append(ret, val.toDateTime())
	}
	return ret
}

// ToInterface converts interfaceValues to interfaceValues.
func (vals interfaceValues) ToInterface() interfaceValues {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return ret
}

// [END] interfaceValues
// ---------------------------------------------------------------------------

