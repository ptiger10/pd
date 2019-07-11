package values

import (
	"github.com/cheekybits/genny/generic"
	"github.com/ptiger10/pd/options"
)

//go:generate genny -in=$GOFILE -out=shared_autogen.go gen "valueType=float64,int64,string,bool,time.Time,interface{}"

// [START] valueTypeValues

// valueType is the generic ValueType that will be replaced by specific types on `make generate`
type valueType generic.Type

// valueTypeValues is a slice of valueType-typed value/null structs.
type valueTypeValues []valueTypeValue

// valueTypeValue is a valueType-typed value/null struct.
type valueTypeValue struct {
	v    valueType
	null bool
}

// newSlicevalueType converts []valueType -> Container with valueTypeValues
func newSlicevalueType(vals []valueType) Container {
	ret := make(valueTypeValues, len(vals))
	for i := 0; i < len(vals); i++ {
		ret[i] = newvalueType(vals[i])
	}
	return Container{&ret, options.PlaceholdervalueType}
}

// Len returns the number of value/null structs in the container.
func (vals *valueTypeValues) Len() int {
	return len(*vals)
}

func (vals *valueTypeValues) Swap(i, j int) {
	(*vals)[i], (*vals)[j] = (*vals)[j], (*vals)[i]
}

// Subset returns the values located at specific index positions.
func (vals *valueTypeValues) Subset(rowPositions []int) Values {
	ret := make(valueTypeValues, len(rowPositions))
	for i := 0; i < len(rowPositions); i++ {
		ret[i] = (*vals)[rowPositions[i]]
	}
	return &ret
}

// Append converts vals2 to valueTypeValues and extends the original valueTypeValues.
func (vals *valueTypeValues) Append(vals2 Values) {
	convertedVals, _ := Convert(vals2, options.PlaceholdervalueType)
	newVals := convertedVals.(*valueTypeValues)
	*vals = append(*vals, *newVals...)
}

// Values returns only the Value fields for the collection of Value/Null structs as an interface slice.
func (vals *valueTypeValues) Values() []interface{} {
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
func (vals *valueTypeValues) Vals() interface{} {
	v := *vals
	ret := make([]valueType, len(v))
	for i := 0; i < len(v); i++ {
		ret[i] = v[i].v
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals *valueTypeValues) Element(position int) Elem {
	v := (*vals)[position]
	return Elem{v.v, v.null}
}

// Value returns the Value field at the specified integer position.
func (vals *valueTypeValues) Value(position int) interface{} {
	return (*vals)[position].v
}

// Value returns the Null field at the specified integer position.
func (vals *valueTypeValues) Null(position int) bool {
	return (*vals)[position].null
}

// Copy transfers every value from the current valueTypeValues container into a new Values container
func (vals *valueTypeValues) Copy() Values {
	v := *vals
	newValues := make(valueTypeValues, len(v))
	for i := 0; i < len(v); i++ {
		newValues[i] = v[i]
	}
	return &newValues
}

// Set overwrites a Value/Null pair at an integer position.
func (vals *valueTypeValues) Set(position int, newVal interface{}) {
	var v interfaceValue
	if isNullInterface(newVal) {
		v = interfaceValue{newVal, true}
	} else {
		v = interfaceValue{newVal, false}
	}
	(*vals)[position] = v.tovalueType()
}

// Drop drops the Value/Null pair at an integer position.
func (vals *valueTypeValues) Drop(pos int) {
	*vals = append((*vals)[:pos], (*vals)[pos+1:]...)
}

// Insert inserts a new Value/Null pair at an integer position.
func (vals *valueTypeValues) Insert(pos int, val interface{}) {
	v := interfaceValue{val, false}
	*vals = append((*vals)[:pos], append([]valueTypeValue{v.tovalueType()}, (*vals)[pos:]...)...)
}

// ToFloat converts valueTypeValues to floatValues.
func (vals *valueTypeValues) ToFloat64() Values {
	ret := make(float64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toFloat64()
	}
	return &ret
}

// ToInt converts valueTypeValues to intValues.
func (vals *valueTypeValues) ToInt64() Values {
	ret := make(int64Values, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toInt64()
	}
	return &ret
}

// ToString converts valueTypeValues to stringValues.
func (vals *valueTypeValues) ToString() Values {
	ret := make(stringValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toString()
	}
	return &ret
}

// ToBool converts valueTypeValues to boolValues.
func (vals *valueTypeValues) ToBool() Values {
	ret := make(boolValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toBool()
	}
	return &ret
}

// ToBool converts valueTypeValues to dateTimeValues.
func (vals *valueTypeValues) ToDateTime() Values {
	ret := make(dateTimeValues, len(*vals))
	for i := 0; i < len(*vals); i++ {
		ret[i] = (*vals)[i].toDateTime()
	}
	return &ret
}

// ToInterface converts valueTypeValues to interfaceValues.
func (vals *valueTypeValues) ToInterface() Values {
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

// [END] valueTypeValues
// ---------------------------------------------------------------------------
var placeholder = true

// the placeholder and this comment are overwritten on `make generate`, but are included so that the [END] comment survives
