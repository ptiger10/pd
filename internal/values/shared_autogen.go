// This file was automatically generated.
// Any changes will be lost if this file is regenerated.
// Run "make generate" to regenerate from template.

package values

import (
	"fmt"
	"time"

	"github.com/ptiger10/pd/options"
)

// [START] float64Values

// float64Values is a slice of float64-typed value/null structs.
type float64Values []float64Value

// float64Value is a float64-typed value/null struct.
type float64Value struct {
	v    float64
	null bool
}

// newSliceFloat64 converts []Float64 -> Factory with float64Values
func newSliceFloat64(vals []float64) Factory {
	var ret float64Values
	for _, val := range vals {
		ret = append(ret, newFloat64(val))
	}
	return Factory{&ret, options.Float64}
}

// Len returns the number of value/null structs in the container.
func (vals *float64Values) Len() int {
	return len(*vals)
}

func (vals *float64Values) Swap(i, j int) {
	(*vals)[i], (*vals)[j] = (*vals)[j], (*vals)[i]
}

// Subset returns the values located at specific index positions.
func (vals *float64Values) Subset(rowPositions []int) Values {
	var ret float64Values
	for _, position := range rowPositions {
		ret = append(ret, (*vals)[position])
	}
	return &ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *float64Values) Vals() interface{} {
	var ret []float64
	for _, val := range *vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals *float64Values) Element(position int) Elem {
	return Elem{(*vals)[position].v, (*vals)[position].null}
}

// Copy transfers every value from the current float64Values container into a new Values container
func (vals *float64Values) Copy() Values {
	newValues := float64Values{}
	for _, val := range *vals {
		newValues = append(newValues, val)
	}
	return &newValues
}

// Set overwrites a Value/Null pair at an integer position.
func (vals *float64Values) Set(position int, newVal interface{}) {
	var v interfaceValue
	if isNullInterface(newVal) {
		v = interfaceValue{newVal, true}
	} else {
		v = interfaceValue{newVal, false}
	}
	(*vals)[position] = v.toFloat64()
}

// Drop drops the Value/Null pair at an integer position.
func (vals *float64Values) Drop(pos int) {
	*vals = append((*vals)[:pos], (*vals)[pos+1:]...)
}

// Insert inserts a new Value/Null pair at an integer position.
func (vals *float64Values) Insert(pos int, val interface{}) {
	v := interfaceValue{val, false}
	*vals = append((*vals)[:pos], append([]float64Value{v.toFloat64()}, (*vals)[pos:]...)...)
}

// ToFloat converts float64Values to floatValues.
func (vals *float64Values) ToFloat64() Values {
	var ret float64Values
	for _, val := range *vals {
		ret = append(ret, val.toFloat64())
	}
	return &ret
}

// ToInt converts float64Values to intValues.
func (vals *float64Values) ToInt64() Values {
	var ret int64Values
	for _, val := range *vals {
		ret = append(ret, val.toInt64())
	}
	return &ret
}

func (val *float64Value) toString() stringValue {
	if val.null {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts float64Values to stringValues.
func (vals *float64Values) ToString() Values {
	var ret stringValues
	for _, val := range *vals {
		ret = append(ret, val.toString())
	}
	return &ret
}

// ToBool converts float64Values to boolValues.
func (vals *float64Values) ToBool() Values {
	var ret boolValues
	for _, val := range *vals {
		ret = append(ret, val.toBool())
	}
	return &ret
}

// ToBool converts float64Values to dateTimeValues.
func (vals *float64Values) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range *vals {
		ret = append(ret, val.toDateTime())
	}
	return &ret
}

// ToInterface converts float64Values to interfaceValues.
func (vals *float64Values) ToInterface() Values {
	var ret interfaceValues
	for _, val := range *vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return &ret
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

// newSliceInt64 converts []Int64 -> Factory with int64Values
func newSliceInt64(vals []int64) Factory {
	var ret int64Values
	for _, val := range vals {
		ret = append(ret, newInt64(val))
	}
	return Factory{&ret, options.Int64}
}

// Len returns the number of value/null structs in the container.
func (vals *int64Values) Len() int {
	return len(*vals)
}

func (vals *int64Values) Swap(i, j int) {
	(*vals)[i], (*vals)[j] = (*vals)[j], (*vals)[i]
}

// Subset returns the values located at specific index positions.
func (vals *int64Values) Subset(rowPositions []int) Values {
	var ret int64Values
	for _, position := range rowPositions {
		ret = append(ret, (*vals)[position])
	}
	return &ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *int64Values) Vals() interface{} {
	var ret []int64
	for _, val := range *vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals *int64Values) Element(position int) Elem {
	return Elem{(*vals)[position].v, (*vals)[position].null}
}

// Copy transfers every value from the current int64Values container into a new Values container
func (vals *int64Values) Copy() Values {
	newValues := int64Values{}
	for _, val := range *vals {
		newValues = append(newValues, val)
	}
	return &newValues
}

// Set overwrites a Value/Null pair at an integer position.
func (vals *int64Values) Set(position int, newVal interface{}) {
	var v interfaceValue
	if isNullInterface(newVal) {
		v = interfaceValue{newVal, true}
	} else {
		v = interfaceValue{newVal, false}
	}
	(*vals)[position] = v.toInt64()
}

// Drop drops the Value/Null pair at an integer position.
func (vals *int64Values) Drop(pos int) {
	*vals = append((*vals)[:pos], (*vals)[pos+1:]...)
}

// Insert inserts a new Value/Null pair at an integer position.
func (vals *int64Values) Insert(pos int, val interface{}) {
	v := interfaceValue{val, false}
	*vals = append((*vals)[:pos], append([]int64Value{v.toInt64()}, (*vals)[pos:]...)...)
}

// ToFloat converts int64Values to floatValues.
func (vals *int64Values) ToFloat64() Values {
	var ret float64Values
	for _, val := range *vals {
		ret = append(ret, val.toFloat64())
	}
	return &ret
}

// ToInt converts int64Values to intValues.
func (vals *int64Values) ToInt64() Values {
	var ret int64Values
	for _, val := range *vals {
		ret = append(ret, val.toInt64())
	}
	return &ret
}

func (val *int64Value) toString() stringValue {
	if val.null {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts int64Values to stringValues.
func (vals *int64Values) ToString() Values {
	var ret stringValues
	for _, val := range *vals {
		ret = append(ret, val.toString())
	}
	return &ret
}

// ToBool converts int64Values to boolValues.
func (vals *int64Values) ToBool() Values {
	var ret boolValues
	for _, val := range *vals {
		ret = append(ret, val.toBool())
	}
	return &ret
}

// ToBool converts int64Values to dateTimeValues.
func (vals *int64Values) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range *vals {
		ret = append(ret, val.toDateTime())
	}
	return &ret
}

// ToInterface converts int64Values to interfaceValues.
func (vals *int64Values) ToInterface() Values {
	var ret interfaceValues
	for _, val := range *vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return &ret
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

// newSliceString converts []String -> Factory with stringValues
func newSliceString(vals []string) Factory {
	var ret stringValues
	for _, val := range vals {
		ret = append(ret, newString(val))
	}
	return Factory{&ret, options.String}
}

// Len returns the number of value/null structs in the container.
func (vals *stringValues) Len() int {
	return len(*vals)
}

func (vals *stringValues) Swap(i, j int) {
	(*vals)[i], (*vals)[j] = (*vals)[j], (*vals)[i]
}

// Subset returns the values located at specific index positions.
func (vals *stringValues) Subset(rowPositions []int) Values {
	var ret stringValues
	for _, position := range rowPositions {
		ret = append(ret, (*vals)[position])
	}
	return &ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *stringValues) Vals() interface{} {
	var ret []string
	for _, val := range *vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals *stringValues) Element(position int) Elem {
	return Elem{(*vals)[position].v, (*vals)[position].null}
}

// Copy transfers every value from the current stringValues container into a new Values container
func (vals *stringValues) Copy() Values {
	newValues := stringValues{}
	for _, val := range *vals {
		newValues = append(newValues, val)
	}
	return &newValues
}

// Set overwrites a Value/Null pair at an integer position.
func (vals *stringValues) Set(position int, newVal interface{}) {
	var v interfaceValue
	if isNullInterface(newVal) {
		v = interfaceValue{newVal, true}
	} else {
		v = interfaceValue{newVal, false}
	}
	(*vals)[position] = v.toString()
}

// Drop drops the Value/Null pair at an integer position.
func (vals *stringValues) Drop(pos int) {
	*vals = append((*vals)[:pos], (*vals)[pos+1:]...)
}

// Insert inserts a new Value/Null pair at an integer position.
func (vals *stringValues) Insert(pos int, val interface{}) {
	v := interfaceValue{val, false}
	*vals = append((*vals)[:pos], append([]stringValue{v.toString()}, (*vals)[pos:]...)...)
}

// ToFloat converts stringValues to floatValues.
func (vals *stringValues) ToFloat64() Values {
	var ret float64Values
	for _, val := range *vals {
		ret = append(ret, val.toFloat64())
	}
	return &ret
}

// ToInt converts stringValues to intValues.
func (vals *stringValues) ToInt64() Values {
	var ret int64Values
	for _, val := range *vals {
		ret = append(ret, val.toInt64())
	}
	return &ret
}

func (val *stringValue) toString() stringValue {
	if val.null {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts stringValues to stringValues.
func (vals *stringValues) ToString() Values {
	var ret stringValues
	for _, val := range *vals {
		ret = append(ret, val.toString())
	}
	return &ret
}

// ToBool converts stringValues to boolValues.
func (vals *stringValues) ToBool() Values {
	var ret boolValues
	for _, val := range *vals {
		ret = append(ret, val.toBool())
	}
	return &ret
}

// ToBool converts stringValues to dateTimeValues.
func (vals *stringValues) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range *vals {
		ret = append(ret, val.toDateTime())
	}
	return &ret
}

// ToInterface converts stringValues to interfaceValues.
func (vals *stringValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range *vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return &ret
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

// newSliceBool converts []Bool -> Factory with boolValues
func newSliceBool(vals []bool) Factory {
	var ret boolValues
	for _, val := range vals {
		ret = append(ret, newBool(val))
	}
	return Factory{&ret, options.Bool}
}

// Len returns the number of value/null structs in the container.
func (vals *boolValues) Len() int {
	return len(*vals)
}

func (vals *boolValues) Swap(i, j int) {
	(*vals)[i], (*vals)[j] = (*vals)[j], (*vals)[i]
}

// Subset returns the values located at specific index positions.
func (vals *boolValues) Subset(rowPositions []int) Values {
	var ret boolValues
	for _, position := range rowPositions {
		ret = append(ret, (*vals)[position])
	}
	return &ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *boolValues) Vals() interface{} {
	var ret []bool
	for _, val := range *vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals *boolValues) Element(position int) Elem {
	return Elem{(*vals)[position].v, (*vals)[position].null}
}

// Copy transfers every value from the current boolValues container into a new Values container
func (vals *boolValues) Copy() Values {
	newValues := boolValues{}
	for _, val := range *vals {
		newValues = append(newValues, val)
	}
	return &newValues
}

// Set overwrites a Value/Null pair at an integer position.
func (vals *boolValues) Set(position int, newVal interface{}) {
	var v interfaceValue
	if isNullInterface(newVal) {
		v = interfaceValue{newVal, true}
	} else {
		v = interfaceValue{newVal, false}
	}
	(*vals)[position] = v.toBool()
}

// Drop drops the Value/Null pair at an integer position.
func (vals *boolValues) Drop(pos int) {
	*vals = append((*vals)[:pos], (*vals)[pos+1:]...)
}

// Insert inserts a new Value/Null pair at an integer position.
func (vals *boolValues) Insert(pos int, val interface{}) {
	v := interfaceValue{val, false}
	*vals = append((*vals)[:pos], append([]boolValue{v.toBool()}, (*vals)[pos:]...)...)
}

// ToFloat converts boolValues to floatValues.
func (vals *boolValues) ToFloat64() Values {
	var ret float64Values
	for _, val := range *vals {
		ret = append(ret, val.toFloat64())
	}
	return &ret
}

// ToInt converts boolValues to intValues.
func (vals *boolValues) ToInt64() Values {
	var ret int64Values
	for _, val := range *vals {
		ret = append(ret, val.toInt64())
	}
	return &ret
}

func (val *boolValue) toString() stringValue {
	if val.null {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts boolValues to stringValues.
func (vals *boolValues) ToString() Values {
	var ret stringValues
	for _, val := range *vals {
		ret = append(ret, val.toString())
	}
	return &ret
}

// ToBool converts boolValues to boolValues.
func (vals *boolValues) ToBool() Values {
	var ret boolValues
	for _, val := range *vals {
		ret = append(ret, val.toBool())
	}
	return &ret
}

// ToBool converts boolValues to dateTimeValues.
func (vals *boolValues) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range *vals {
		ret = append(ret, val.toDateTime())
	}
	return &ret
}

// ToInterface converts boolValues to interfaceValues.
func (vals *boolValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range *vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return &ret
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

// newSliceDateTime converts []DateTime -> Factory with dateTimeValues
func newSliceDateTime(vals []time.Time) Factory {
	var ret dateTimeValues
	for _, val := range vals {
		ret = append(ret, newDateTime(val))
	}
	return Factory{&ret, options.DateTime}
}

// Len returns the number of value/null structs in the container.
func (vals *dateTimeValues) Len() int {
	return len(*vals)
}

func (vals *dateTimeValues) Swap(i, j int) {
	(*vals)[i], (*vals)[j] = (*vals)[j], (*vals)[i]
}

// Subset returns the values located at specific index positions.
func (vals *dateTimeValues) Subset(rowPositions []int) Values {
	var ret dateTimeValues
	for _, position := range rowPositions {
		ret = append(ret, (*vals)[position])
	}
	return &ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *dateTimeValues) Vals() interface{} {
	var ret []time.Time
	for _, val := range *vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals *dateTimeValues) Element(position int) Elem {
	return Elem{(*vals)[position].v, (*vals)[position].null}
}

// Copy transfers every value from the current dateTimeValues container into a new Values container
func (vals *dateTimeValues) Copy() Values {
	newValues := dateTimeValues{}
	for _, val := range *vals {
		newValues = append(newValues, val)
	}
	return &newValues
}

// Set overwrites a Value/Null pair at an integer position.
func (vals *dateTimeValues) Set(position int, newVal interface{}) {
	var v interfaceValue
	if isNullInterface(newVal) {
		v = interfaceValue{newVal, true}
	} else {
		v = interfaceValue{newVal, false}
	}
	(*vals)[position] = v.toDateTime()
}

// Drop drops the Value/Null pair at an integer position.
func (vals *dateTimeValues) Drop(pos int) {
	*vals = append((*vals)[:pos], (*vals)[pos+1:]...)
}

// Insert inserts a new Value/Null pair at an integer position.
func (vals *dateTimeValues) Insert(pos int, val interface{}) {
	v := interfaceValue{val, false}
	*vals = append((*vals)[:pos], append([]dateTimeValue{v.toDateTime()}, (*vals)[pos:]...)...)
}

// ToFloat converts dateTimeValues to floatValues.
func (vals *dateTimeValues) ToFloat64() Values {
	var ret float64Values
	for _, val := range *vals {
		ret = append(ret, val.toFloat64())
	}
	return &ret
}

// ToInt converts dateTimeValues to intValues.
func (vals *dateTimeValues) ToInt64() Values {
	var ret int64Values
	for _, val := range *vals {
		ret = append(ret, val.toInt64())
	}
	return &ret
}

func (val *dateTimeValue) toString() stringValue {
	if val.null {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts dateTimeValues to stringValues.
func (vals *dateTimeValues) ToString() Values {
	var ret stringValues
	for _, val := range *vals {
		ret = append(ret, val.toString())
	}
	return &ret
}

// ToBool converts dateTimeValues to boolValues.
func (vals *dateTimeValues) ToBool() Values {
	var ret boolValues
	for _, val := range *vals {
		ret = append(ret, val.toBool())
	}
	return &ret
}

// ToBool converts dateTimeValues to dateTimeValues.
func (vals *dateTimeValues) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range *vals {
		ret = append(ret, val.toDateTime())
	}
	return &ret
}

// ToInterface converts dateTimeValues to interfaceValues.
func (vals *dateTimeValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range *vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return &ret
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

// newSliceInterface converts []Interface -> Factory with interfaceValues
func newSliceInterface(vals []interface{}) Factory {
	var ret interfaceValues
	for _, val := range vals {
		ret = append(ret, newInterface(val))
	}
	return Factory{&ret, options.Interface}
}

// Len returns the number of value/null structs in the container.
func (vals *interfaceValues) Len() int {
	return len(*vals)
}

func (vals *interfaceValues) Swap(i, j int) {
	(*vals)[i], (*vals)[j] = (*vals)[j], (*vals)[i]
}

// Subset returns the values located at specific index positions.
func (vals *interfaceValues) Subset(rowPositions []int) Values {
	var ret interfaceValues
	for _, position := range rowPositions {
		ret = append(ret, (*vals)[position])
	}
	return &ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *interfaceValues) Vals() interface{} {
	var ret []interface{}
	for _, val := range *vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals *interfaceValues) Element(position int) Elem {
	return Elem{(*vals)[position].v, (*vals)[position].null}
}

// Copy transfers every value from the current interfaceValues container into a new Values container
func (vals *interfaceValues) Copy() Values {
	newValues := interfaceValues{}
	for _, val := range *vals {
		newValues = append(newValues, val)
	}
	return &newValues
}

// Set overwrites a Value/Null pair at an integer position.
func (vals *interfaceValues) Set(position int, newVal interface{}) {
	var v interfaceValue
	if isNullInterface(newVal) {
		v = interfaceValue{newVal, true}
	} else {
		v = interfaceValue{newVal, false}
	}
	(*vals)[position] = v.toInterface()
}

// Drop drops the Value/Null pair at an integer position.
func (vals *interfaceValues) Drop(pos int) {
	*vals = append((*vals)[:pos], (*vals)[pos+1:]...)
}

// Insert inserts a new Value/Null pair at an integer position.
func (vals *interfaceValues) Insert(pos int, val interface{}) {
	v := interfaceValue{val, false}
	*vals = append((*vals)[:pos], append([]interfaceValue{v.toInterface()}, (*vals)[pos:]...)...)
}

// ToFloat converts interfaceValues to floatValues.
func (vals *interfaceValues) ToFloat64() Values {
	var ret float64Values
	for _, val := range *vals {
		ret = append(ret, val.toFloat64())
	}
	return &ret
}

// ToInt converts interfaceValues to intValues.
func (vals *interfaceValues) ToInt64() Values {
	var ret int64Values
	for _, val := range *vals {
		ret = append(ret, val.toInt64())
	}
	return &ret
}

func (val *interfaceValue) toString() stringValue {
	if val.null {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts interfaceValues to stringValues.
func (vals *interfaceValues) ToString() Values {
	var ret stringValues
	for _, val := range *vals {
		ret = append(ret, val.toString())
	}
	return &ret
}

// ToBool converts interfaceValues to boolValues.
func (vals *interfaceValues) ToBool() Values {
	var ret boolValues
	for _, val := range *vals {
		ret = append(ret, val.toBool())
	}
	return &ret
}

// ToBool converts interfaceValues to dateTimeValues.
func (vals *interfaceValues) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range *vals {
		ret = append(ret, val.toDateTime())
	}
	return &ret
}

// ToInterface converts interfaceValues to interfaceValues.
func (vals *interfaceValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range *vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return &ret
}

// [END] interfaceValues
// ---------------------------------------------------------------------------
