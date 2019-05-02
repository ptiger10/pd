package values

import (
	"fmt"
	"log"
)

// [START FloatValues]

func (vals FloatValues) count() int {
	var count int
	for _, val := range vals {
		if !val.null {
			count++
		}
	}
	return count
}

// Len returns the total number of values in the collection (including null values)
func (vals FloatValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals FloatValues) In(positions []int) Values {
	var ret FloatValues
	for _, position := range positions {
		if position >= len(vals) {
			log.Panicf("Unable to get value: index out of range: %d", position)
		}
		ret = append(ret, vals[position])
	}
	return ret
}

// All returns all of the value elements in the collection (including null values)
// as an interface that supports type assertions
func (vals FloatValues) All() interface{} {
	return vals
}

// Vals returns only the Value fields for the collection of Value/Null structs.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals FloatValues) Vals() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals FloatValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals FloatValues) Null() []int {
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
func (vals FloatValues) Element(position int) []interface{} {
	return []interface{}{vals[position].v, vals[position].null}
}

// ToString converts the values to String values
func (vals FloatValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, String("", true))
		} else {
			ret = append(ret, String(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// [END FloatValues]

// [START IntValues]

func (vals IntValues) count() int {
	var count int
	for _, val := range vals {
		if !val.null {
			count++
		}
	}
	return count
}

// Len returns the total number of values in the collection (including null values)
func (vals IntValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals IntValues) In(positions []int) Values {
	var ret IntValues
	for _, position := range positions {
		if position >= len(vals) {
			log.Panicf("Unable to get value: index out of range: %d", position)
		}
		ret = append(ret, vals[position])
	}
	return ret
}

// All returns all of the value elementsin the collection (including null values)
// as an interface that supports type assertions
func (vals IntValues) All() interface{} {
	return vals
}

// Vals returns only the Value fields for the collection of Value/Null structs.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals IntValues) Vals() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals IntValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals IntValues) Null() []int {
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
func (vals IntValues) Element(position int) []interface{} {
	return []interface{}{vals[position].v, vals[position].null}
}

// ToString converts the values to String values
func (vals IntValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, String("", true))
		} else {
			ret = append(ret, String(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// [END IntValues]

// [START StringValues]

func (vals StringValues) count() int {
	var count int
	for _, val := range vals {
		if !val.null {
			count++
		}
	}
	return count
}

// Len returns the total number of values in the collection (including null values)
func (vals StringValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals StringValues) In(positions []int) Values {
	var ret StringValues
	for _, position := range positions {
		if position >= len(vals) {
			log.Panicf("Unable to get value: index out of range: %d", position)
		}
		ret = append(ret, vals[position])
	}
	return ret
}

// All returns all of the value elements in the collection (including null values)
// as an interface that supports type assertions
func (vals StringValues) All() interface{} {
	return vals
}

// Vals returns only the Value fields for the collection of Value/Null structs.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals StringValues) Vals() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals StringValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals StringValues) Null() []int {
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
func (vals StringValues) Element(position int) []interface{} {
	return []interface{}{vals[position].v, vals[position].null}
}

// ToString converts the values to String values
func (vals StringValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, String("", true))
		} else {
			ret = append(ret, String(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// [END StringValues]

// [START BoolValues]

func (vals BoolValues) count() int {
	var count int
	for _, val := range vals {
		if !val.null {
			count++
		}
	}
	return count
}

// Len returns the total number of values in the collection (including null values)
func (vals BoolValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals BoolValues) In(positions []int) Values {
	var ret BoolValues
	for _, position := range positions {
		if position >= len(vals) {
			log.Panicf("Unable to get value: index out of range: %d", position)
		}
		ret = append(ret, vals[position])
	}
	return ret
}

// All returns all of the value elements in the collection (including null values)
// as an
func (vals BoolValues) All() interface{} {
	return vals
}

// Vals returns only the Value fields for the collection of Value/Null structs.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals BoolValues) Vals() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals BoolValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals BoolValues) Null() []int {
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
func (vals BoolValues) Element(position int) []interface{} {
	return []interface{}{vals[position].v, vals[position].null}
}

// ToString converts the values to String values
func (vals BoolValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, String("", true))
		} else {
			ret = append(ret, String(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// [END BoolValues]

// [START DateTimeValues]

func (vals DateTimeValues) count() int {
	var count int
	for _, val := range vals {
		if !val.null {
			count++
		}
	}
	return count
}

// Len returns the total number of values in the collection (including null values)
func (vals DateTimeValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals DateTimeValues) In(positions []int) Values {
	var ret DateTimeValues
	for _, position := range positions {
		if position >= len(vals) {
			log.Panicf("Unable to get value: index out of range: %d", position)
		}
		ret = append(ret, vals[position])
	}
	return ret
}

// All returns all of the value elements in the collection (including null values)
// as an interface that supports type assertions
func (vals DateTimeValues) All() interface{} {
	return vals
}

// Vals returns only the Value fields for the collection of Value/Null structs.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals DateTimeValues) Vals() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals DateTimeValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals DateTimeValues) Null() []int {
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
func (vals DateTimeValues) Element(position int) []interface{} {
	return []interface{}{vals[position].v, vals[position].null}
}

// ToString converts the values to String values
func (vals DateTimeValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, String("", true))
		} else {
			ret = append(ret, String(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// [END DateTimeValues]

// [START InterfaceValues]

func (vals InterfaceValues) count() int {
	var count int
	for _, val := range vals {
		if !val.null {
			count++
		}
	}
	return count
}

// Len returns the total number of values in the collection (including null values)
func (vals InterfaceValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals InterfaceValues) In(positions []int) Values {
	var ret InterfaceValues
	for _, position := range positions {
		if position >= len(vals) {
			log.Panicf("Unable to get value: index out of range: %d", position)
		}
		ret = append(ret, vals[position])
	}
	return ret
}

// All returns all of the value elements in the collection (including null values)
// as an interface that supports type assertions
func (vals InterfaceValues) All() interface{} {
	return vals
}

// Vals returns only the Value fields for the collection of Value/Null structs.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals InterfaceValues) Vals() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals InterfaceValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals InterfaceValues) Null() []int {
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
func (vals InterfaceValues) Element(position int) []interface{} {
	return []interface{}{vals[position].v, vals[position].null}
}

// ToString converts the values to String values
func (vals InterfaceValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, String("", true))
		} else {
			ret = append(ret, String(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// [END InterfaceValues]
