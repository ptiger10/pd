package series

import (
	"math"
	"reflect"
	"sort"
	"strconv"

	"github.com/gonum/floats"
)

// data type
// ------------------------------------------------

type floatValues []floatValue
type floatValue struct {
	v    float64
	null bool
}

// convenience functions
// ------------------------------------------------

func (vals floatValues) valid() ([]float64, []int) {
	var valid []float64
	var nullMap []int
	for i, val := range vals {
		if !val.null {
			valid = append(valid, float64(val.v))
		} else {
			nullMap = append(nullMap, i)
		}
	}
	return valid, nullMap
}

// methods
// ------------------------------------------------
func (vals floatValues) sum() float64 {
	var sum float64
	for _, val := range vals {
		if !val.null {
			sum += val.v
		}
	}
	return sum
}

func (vals floatValues) median() float64 {
	valid, _ := vals.valid()
	sort.Float64s(valid)
	mNumber := len(valid) / 2
	if len(valid)%2 != 0 { // checks if sequence has odd number of elements
		return valid[mNumber]
	}
	return (valid[mNumber-1] + valid[mNumber]) / 2
}

func (vals floatValues) mean() float64 {
	var sum float64
	var count int
	for _, val := range vals {
		if !val.null {
			sum += val.v
			count++
		}
	}
	return sum / float64(count)
}

// transcribe copies []float64 into a new floatValues object but inserts null values wherever they existed in the original
func (vals floatValues) transcribe(valid []float64, nullMap []int) floatValues {
	var fv floatValues
	var nullCounter int
	for i := 0; i < len(vals); i++ {
		validCounter := i - nullCounter
		if nullCounter == len(nullMap) { // in this case, there are no more nil values to transcribe
			fv = append(fv, floatValue{v: valid[validCounter]})
		} else if nullMap[nullCounter] == i {
			fv = append(fv, floatValue{null: true, v: math.NaN()})
			nullCounter++
		} else {
			fv = append(fv, floatValue{v: valid[validCounter]})
		}
	}
	return fv
}

func (vals floatValues) addConst(c float64) Series {
	valid, nullMap := vals.valid()
	floats.AddConst(c, valid)
	fv := vals.transcribe(valid, nullMap)
	return Series{
		Values: fv,
		Kind:   Float,
	}
}

// constructor functions
// ------------------------------------------------
func floatToFloatValues(data interface{}) floatValues {
	var vals []floatValue
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		val := d.Index(i).Float()
		vals = append(vals, floatValue{v: val})
		if math.IsNaN(val) {
			vals[i].null = true
		}
	}
	return floatValues(vals)
}

func interfaceToFloatValues(d reflect.Value) floatValues {
	var vals []floatValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, floatValue{null: true})
		case reflect.String:
			vals = append(vals, stringToFloat(v.String()))
		case reflect.Float32, reflect.Float64:
			val := v.Float()
			vals = append(vals, floatValue{v: val})
			if math.IsNaN(val) {
				vals[i].null = true
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, floatValue{v: float64(v.Int())})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, floatValue{v: float64(v.Uint())})
		default:
			vals = append(vals, floatValue{null: true})
		}

	}
	return floatValues(vals)
}

func stringToFloat(v string) floatValue {
	val, err := strconv.ParseFloat(v, 64)
	if err != nil || math.IsNaN(val) {
		return floatValue{null: true, v: math.NaN()}
	}
	return floatValue{v: val}
}
