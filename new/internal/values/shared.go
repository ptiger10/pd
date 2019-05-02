package values

import (
	"fmt"
	"log"
)

// [START FloatValues]

func (vals FloatValues) count() int {
	var count int
	for _, val := range vals {
		if !val.Null {
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

// All returns all of the values in the collection (including null values)
// as an interface slice
func (vals FloatValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.V)
	}
	return ret
}

// ToString converts the values to String values
func (vals FloatValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		ret = append(ret, String(fmt.Sprint(val.V), val.Null))
	}
	return ret
}

// [END FloatValues]

// [START IntValues]

func (vals IntValues) count() int {
	var count int
	for _, val := range vals {
		if !val.Null {
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

// All returns all of the values in the collection (including null values)
// as an interface slice
func (vals IntValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.V)
	}
	return ret
}

// ToString converts the values to String values
func (vals IntValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		ret = append(ret, String(fmt.Sprint(val.V), val.Null))
	}
	return ret
}

// [END IntValues]

// [START StringValues]

func (vals StringValues) count() int {
	var count int
	for _, val := range vals {
		if !val.Null {
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

// All returns all of the values in the collection (including null values)
// as an interface slice
func (vals StringValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.V)
	}
	return ret
}

// ToString converts the values to String values
func (vals StringValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		ret = append(ret, String(fmt.Sprint(val.V), val.Null))
	}
	return ret
}

// [END StringValues]

// [START BoolValues]

func (vals BoolValues) count() int {
	var count int
	for _, val := range vals {
		if !val.Null {
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

// All returns all of the values in the collection (including null values)
// as an interface slice
func (vals BoolValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.V)
	}
	return ret
}

// ToString converts the values to String values
func (vals BoolValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		ret = append(ret, String(fmt.Sprint(val.V), val.Null))
	}
	return ret
}

// [END BoolValues]

// [START DateTimeValues]

func (vals DateTimeValues) count() int {
	var count int
	for _, val := range vals {
		if !val.Null {
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

// All returns all of the values in the collection (including null values)
// as an interface slice
func (vals DateTimeValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.V)
	}
	return ret
}

// ToString converts the values to String values
func (vals DateTimeValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		ret = append(ret, String(fmt.Sprint(val.V), val.Null))
	}
	return ret
}

// [END DateTimeValues]

// [START InterfaceValues]

func (vals InterfaceValues) count() int {
	var count int
	for _, val := range vals {
		if !val.Null {
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

// All returns all of the values in the collection (including null values)
// as an interface slice
func (vals InterfaceValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.V)
	}
	return ret
}

// ToString converts the values to String values
func (vals InterfaceValues) ToString() Values {
	var ret StringValues
	for _, val := range vals {
		ret = append(ret, String(fmt.Sprint(val.V), val.Null))
	}
	return ret
}

// [END InterfaceValues]
