// This file was automatically generated.
// Any changes will be lost if this file is regenerated.
// Run "make generate" to regenerate from template.

package values

import (
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

// newSliceFloat64 converts []Float64 -> Container with float64Values
func newSliceFloat64(vals []float64) Container {
	ret := make(float64Values, len(vals))
	for i := 0; i < len(vals); i++ {
		ret[i] = newFloat64(vals[i])
	}
	return Container{&ret, options.Float64}
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
	ret := make(float64Values, len(rowPositions))
	for i := 0; i < len(rowPositions); i++ {
		ret[i] = (*vals)[rowPositions[i]]
	}
	return &ret
}

// Append converts vals2 to float64Values and extends the original float64Values.
func (vals *float64Values) Append(vals2 Values) {
	convertedVals, _ := Convert(vals2, options.Float64)
	newVals := convertedVals.(*float64Values)
	*vals = append(*vals, *newVals...)
}

// Values returns only the Value fields for the collection of Value/Null structs as an interface slice.
func (vals *float64Values) Values() []interface{} {
	v := *vals
	ret := make([]interface{}, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *float64Values) Vals() interface{} {
	v := *vals
	ret := make([]float64, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Value returns the Value field at the specified integer position.
func (vals *float64Values) Value(position int) interface{} {
	return (*vals)[position].v
}

// Value returns the Null field at the specified integer position.
func (vals *float64Values) Null(position int) bool {
	return (*vals)[position].null
}

// Copy transfers every value from the current float64Values container into a new Values container
func (vals *float64Values) Copy() Values {
	v := *vals
	newValues := make(float64Values, len(v))
	for i := 0; i < len(v); i++ {
		newValues[i] = v[i]
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
	ret := make(float64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toFloat64()
	}
	return &ret
}

// ToInt converts float64Values to intValues.
func (vals *float64Values) ToInt64() Values {
	ret := make(int64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toInt64()
	}
	return &ret
}

// ToString converts float64Values to stringValues.
func (vals *float64Values) ToString() Values {
	ret := make(stringValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toString()
	}
	return &ret
}

// ToBool converts float64Values to boolValues.
func (vals *float64Values) ToBool() Values {
	ret := make(boolValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toBool()
	}
	return &ret
}

// ToBool converts float64Values to dateTimeValues.
func (vals *float64Values) ToDateTime() Values {
	ret := make(dateTimeValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toDateTime()
	}
	return &ret
}

// ToInterface converts float64Values to interfaceValues.
func (vals *float64Values) ToInterface() Values {
	ret := make(interfaceValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		if (*vals)[i].null {
			ret[i] = interfaceValue{(*vals)[i].v, true}
		} else {
			ret[i] = interfaceValue{(*vals)[i].v, false}
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

// newSliceInt64 converts []Int64 -> Container with int64Values
func newSliceInt64(vals []int64) Container {
	ret := make(int64Values, len(vals))
	for i := 0; i < len(vals); i++ {
		ret[i] = newInt64(vals[i])
	}
	return Container{&ret, options.Int64}
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
	ret := make(int64Values, len(rowPositions))
	for i := 0; i < len(rowPositions); i++ {
		ret[i] = (*vals)[rowPositions[i]]
	}
	return &ret
}

// Append converts vals2 to int64Values and extends the original int64Values.
func (vals *int64Values) Append(vals2 Values) {
	convertedVals, _ := Convert(vals2, options.Int64)
	newVals := convertedVals.(*int64Values)
	*vals = append(*vals, *newVals...)
}

// Values returns only the Value fields for the collection of Value/Null structs as an interface slice.
func (vals *int64Values) Values() []interface{} {
	v := *vals
	ret := make([]interface{}, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *int64Values) Vals() interface{} {
	v := *vals
	ret := make([]int64, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Value returns the Value field at the specified integer position.
func (vals *int64Values) Value(position int) interface{} {
	return (*vals)[position].v
}

// Value returns the Null field at the specified integer position.
func (vals *int64Values) Null(position int) bool {
	return (*vals)[position].null
}

// Copy transfers every value from the current int64Values container into a new Values container
func (vals *int64Values) Copy() Values {
	v := *vals
	newValues := make(int64Values, len(v))
	for i := 0; i < len(v); i++ {
		newValues[i] = v[i]
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
	ret := make(float64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toFloat64()
	}
	return &ret
}

// ToInt converts int64Values to intValues.
func (vals *int64Values) ToInt64() Values {
	ret := make(int64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toInt64()
	}
	return &ret
}

// ToString converts int64Values to stringValues.
func (vals *int64Values) ToString() Values {
	ret := make(stringValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toString()
	}
	return &ret
}

// ToBool converts int64Values to boolValues.
func (vals *int64Values) ToBool() Values {
	ret := make(boolValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toBool()
	}
	return &ret
}

// ToBool converts int64Values to dateTimeValues.
func (vals *int64Values) ToDateTime() Values {
	ret := make(dateTimeValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toDateTime()
	}
	return &ret
}

// ToInterface converts int64Values to interfaceValues.
func (vals *int64Values) ToInterface() Values {
	ret := make(interfaceValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		if (*vals)[i].null {
			ret[i] = interfaceValue{(*vals)[i].v, true}
		} else {
			ret[i] = interfaceValue{(*vals)[i].v, false}
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

// newSliceString converts []String -> Container with stringValues
func newSliceString(vals []string) Container {
	ret := make(stringValues, len(vals))
	for i := 0; i < len(vals); i++ {
		ret[i] = newString(vals[i])
	}
	return Container{&ret, options.String}
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
	ret := make(stringValues, len(rowPositions))
	for i := 0; i < len(rowPositions); i++ {
		ret[i] = (*vals)[rowPositions[i]]
	}
	return &ret
}

// Append converts vals2 to stringValues and extends the original stringValues.
func (vals *stringValues) Append(vals2 Values) {
	convertedVals, _ := Convert(vals2, options.String)
	newVals := convertedVals.(*stringValues)
	*vals = append(*vals, *newVals...)
}

// Values returns only the Value fields for the collection of Value/Null structs as an interface slice.
func (vals *stringValues) Values() []interface{} {
	v := *vals
	ret := make([]interface{}, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *stringValues) Vals() interface{} {
	v := *vals
	ret := make([]string, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Value returns the Value field at the specified integer position.
func (vals *stringValues) Value(position int) interface{} {
	return (*vals)[position].v
}

// Value returns the Null field at the specified integer position.
func (vals *stringValues) Null(position int) bool {
	return (*vals)[position].null
}

// Copy transfers every value from the current stringValues container into a new Values container
func (vals *stringValues) Copy() Values {
	v := *vals
	newValues := make(stringValues, len(v))
	for i := 0; i < len(v); i++ {
		newValues[i] = v[i]
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
	ret := make(float64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toFloat64()
	}
	return &ret
}

// ToInt converts stringValues to intValues.
func (vals *stringValues) ToInt64() Values {
	ret := make(int64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toInt64()
	}
	return &ret
}

// ToString converts stringValues to stringValues.
func (vals *stringValues) ToString() Values {
	ret := make(stringValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toString()
	}
	return &ret
}

// ToBool converts stringValues to boolValues.
func (vals *stringValues) ToBool() Values {
	ret := make(boolValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toBool()
	}
	return &ret
}

// ToBool converts stringValues to dateTimeValues.
func (vals *stringValues) ToDateTime() Values {
	ret := make(dateTimeValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toDateTime()
	}
	return &ret
}

// ToInterface converts stringValues to interfaceValues.
func (vals *stringValues) ToInterface() Values {
	ret := make(interfaceValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		if (*vals)[i].null {
			ret[i] = interfaceValue{(*vals)[i].v, true}
		} else {
			ret[i] = interfaceValue{(*vals)[i].v, false}
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

// newSliceBool converts []Bool -> Container with boolValues
func newSliceBool(vals []bool) Container {
	ret := make(boolValues, len(vals))
	for i := 0; i < len(vals); i++ {
		ret[i] = newBool(vals[i])
	}
	return Container{&ret, options.Bool}
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
	ret := make(boolValues, len(rowPositions))
	for i := 0; i < len(rowPositions); i++ {
		ret[i] = (*vals)[rowPositions[i]]
	}
	return &ret
}

// Append converts vals2 to boolValues and extends the original boolValues.
func (vals *boolValues) Append(vals2 Values) {
	convertedVals, _ := Convert(vals2, options.Bool)
	newVals := convertedVals.(*boolValues)
	*vals = append(*vals, *newVals...)
}

// Values returns only the Value fields for the collection of Value/Null structs as an interface slice.
func (vals *boolValues) Values() []interface{} {
	v := *vals
	ret := make([]interface{}, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *boolValues) Vals() interface{} {
	v := *vals
	ret := make([]bool, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Value returns the Value field at the specified integer position.
func (vals *boolValues) Value(position int) interface{} {
	return (*vals)[position].v
}

// Value returns the Null field at the specified integer position.
func (vals *boolValues) Null(position int) bool {
	return (*vals)[position].null
}

// Copy transfers every value from the current boolValues container into a new Values container
func (vals *boolValues) Copy() Values {
	v := *vals
	newValues := make(boolValues, len(v))
	for i := 0; i < len(v); i++ {
		newValues[i] = v[i]
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
	ret := make(float64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toFloat64()
	}
	return &ret
}

// ToInt converts boolValues to intValues.
func (vals *boolValues) ToInt64() Values {
	ret := make(int64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toInt64()
	}
	return &ret
}

// ToString converts boolValues to stringValues.
func (vals *boolValues) ToString() Values {
	ret := make(stringValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toString()
	}
	return &ret
}

// ToBool converts boolValues to boolValues.
func (vals *boolValues) ToBool() Values {
	ret := make(boolValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toBool()
	}
	return &ret
}

// ToBool converts boolValues to dateTimeValues.
func (vals *boolValues) ToDateTime() Values {
	ret := make(dateTimeValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toDateTime()
	}
	return &ret
}

// ToInterface converts boolValues to interfaceValues.
func (vals *boolValues) ToInterface() Values {
	ret := make(interfaceValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		if (*vals)[i].null {
			ret[i] = interfaceValue{(*vals)[i].v, true}
		} else {
			ret[i] = interfaceValue{(*vals)[i].v, false}
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

// newSliceDateTime converts []DateTime -> Container with dateTimeValues
func newSliceDateTime(vals []time.Time) Container {
	ret := make(dateTimeValues, len(vals))
	for i := 0; i < len(vals); i++ {
		ret[i] = newDateTime(vals[i])
	}
	return Container{&ret, options.DateTime}
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
	ret := make(dateTimeValues, len(rowPositions))
	for i := 0; i < len(rowPositions); i++ {
		ret[i] = (*vals)[rowPositions[i]]
	}
	return &ret
}

// Append converts vals2 to dateTimeValues and extends the original dateTimeValues.
func (vals *dateTimeValues) Append(vals2 Values) {
	convertedVals, _ := Convert(vals2, options.DateTime)
	newVals := convertedVals.(*dateTimeValues)
	*vals = append(*vals, *newVals...)
}

// Values returns only the Value fields for the collection of Value/Null structs as an interface slice.
func (vals *dateTimeValues) Values() []interface{} {
	v := *vals
	ret := make([]interface{}, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *dateTimeValues) Vals() interface{} {
	v := *vals
	ret := make([]time.Time, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Value returns the Value field at the specified integer position.
func (vals *dateTimeValues) Value(position int) interface{} {
	return (*vals)[position].v
}

// Value returns the Null field at the specified integer position.
func (vals *dateTimeValues) Null(position int) bool {
	return (*vals)[position].null
}

// Copy transfers every value from the current dateTimeValues container into a new Values container
func (vals *dateTimeValues) Copy() Values {
	v := *vals
	newValues := make(dateTimeValues, len(v))
	for i := 0; i < len(v); i++ {
		newValues[i] = v[i]
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
	ret := make(float64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toFloat64()
	}
	return &ret
}

// ToInt converts dateTimeValues to intValues.
func (vals *dateTimeValues) ToInt64() Values {
	ret := make(int64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toInt64()
	}
	return &ret
}

// ToString converts dateTimeValues to stringValues.
func (vals *dateTimeValues) ToString() Values {
	ret := make(stringValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toString()
	}
	return &ret
}

// ToBool converts dateTimeValues to boolValues.
func (vals *dateTimeValues) ToBool() Values {
	ret := make(boolValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toBool()
	}
	return &ret
}

// ToBool converts dateTimeValues to dateTimeValues.
func (vals *dateTimeValues) ToDateTime() Values {
	ret := make(dateTimeValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toDateTime()
	}
	return &ret
}

// ToInterface converts dateTimeValues to interfaceValues.
func (vals *dateTimeValues) ToInterface() Values {
	ret := make(interfaceValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		if (*vals)[i].null {
			ret[i] = interfaceValue{(*vals)[i].v, true}
		} else {
			ret[i] = interfaceValue{(*vals)[i].v, false}
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

// newSliceInterface converts []Interface -> Container with interfaceValues
func newSliceInterface(vals []interface{}) Container {
	ret := make(interfaceValues, len(vals))
	for i := 0; i < len(vals); i++ {
		ret[i] = newInterface(vals[i])
	}
	return Container{&ret, options.Interface}
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
	ret := make(interfaceValues, len(rowPositions))
	for i := 0; i < len(rowPositions); i++ {
		ret[i] = (*vals)[rowPositions[i]]
	}
	return &ret
}

// Append converts vals2 to interfaceValues and extends the original interfaceValues.
func (vals *interfaceValues) Append(vals2 Values) {
	convertedVals, _ := Convert(vals2, options.Interface)
	newVals := convertedVals.(*interfaceValues)
	*vals = append(*vals, *newVals...)
}

// Values returns only the Value fields for the collection of Value/Null structs as an interface slice.
func (vals *interfaceValues) Values() []interface{} {
	v := *vals
	ret := make([]interface{}, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *interfaceValues) Vals() interface{} {
	v := *vals
	ret := make([]interface{}, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Value returns the Value field at the specified integer position.
func (vals *interfaceValues) Value(position int) interface{} {
	return (*vals)[position].v
}

// Value returns the Null field at the specified integer position.
func (vals *interfaceValues) Null(position int) bool {
	return (*vals)[position].null
}

// Copy transfers every value from the current interfaceValues container into a new Values container
func (vals *interfaceValues) Copy() Values {
	v := *vals
	newValues := make(interfaceValues, len(v))
	for i := 0; i < len(v); i++ {
		newValues[i] = v[i]
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
	ret := make(float64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toFloat64()
	}
	return &ret
}

// ToInt converts interfaceValues to intValues.
func (vals *interfaceValues) ToInt64() Values {
	ret := make(int64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toInt64()
	}
	return &ret
}

// ToString converts interfaceValues to stringValues.
func (vals *interfaceValues) ToString() Values {
	ret := make(stringValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toString()
	}
	return &ret
}

// ToBool converts interfaceValues to boolValues.
func (vals *interfaceValues) ToBool() Values {
	ret := make(boolValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toBool()
	}
	return &ret
}

// ToBool converts interfaceValues to dateTimeValues.
func (vals *interfaceValues) ToDateTime() Values {
	ret := make(dateTimeValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toDateTime()
	}
	return &ret
}

// ToInterface converts interfaceValues to interfaceValues.
func (vals *interfaceValues) ToInterface() Values {
	ret := make(interfaceValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		if (*vals)[i].null {
			ret[i] = interfaceValue{(*vals)[i].v, true}
		} else {
			ret[i] = interfaceValue{(*vals)[i].v, false}
		}
	}
	return &ret
}

// [END] interfaceValues
// ---------------------------------------------------------------------------
