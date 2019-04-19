package series

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"

	"github.com/gonum/floats"       // uses optimized gonum/floats methods where available
	"github.com/montanaflynn/stats" // and stats package otherwise
)

// Data Type
// ------------------------------------------------

type floatValues []floatValue
type floatValue struct {
	v    float64
	null bool
}

// Convenience Methods
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

// Methods
// ------------------------------------------------
func (vals floatValues) sum() float64 {
	valid, _ := vals.valid()
	return floats.Sum(valid)
}

func (vals floatValues) count() int {
	valid, _ := vals.valid()
	return len(valid)
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

func (vals floatValues) median() float64 {
	valid, _ := vals.valid()
	if len(valid) == 0 {
		return math.NaN()
	}
	sort.Float64s(valid)
	mNumber := len(valid) / 2
	if len(valid)%2 != 0 { // checks if sequence has odd number of elements
		return valid[mNumber]
	}
	return (valid[mNumber-1] + valid[mNumber]) / 2
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

func (vals floatValues) min() float64 {
	valid, _ := vals.valid()
	if len(valid) == 0 {
		return math.NaN()
	}
	return floats.Min(valid)
}

func (vals floatValues) max() float64 {
	valid, _ := vals.valid()
	if len(valid) == 0 {
		return math.NaN()
	}
	return floats.Max(valid)
}

func (vals floatValues) quartiles() stats.Quartiles {
	valid, _ := vals.valid()
	val, err := stats.Quartile(valid)
	if err != nil {
		return stats.Quartiles{Q1: math.NaN(), Q2: math.NaN(), Q3: math.NaN()}
	}
	return val
}

func (vals floatValues) describe() string {
	offset := 7
	precision := 4
	l := len(vals)
	v := vals.count()
	quartiles := vals.quartiles()
	len := fmt.Sprintf("%-*s %d\n", offset, "len", l)
	valid := fmt.Sprintf("%-*s %d\n", offset, "valid", v)
	null := fmt.Sprintf("%-*s %d\n", offset, "null", l-v)
	mean := fmt.Sprintf("%-*s %.*f\n", offset, "mean", precision, vals.mean())
	min := fmt.Sprintf("%-*s %.*f\n", offset, "min", precision, vals.min())
	firstQ := fmt.Sprintf("%-*s %.*f\n", offset, "25%", precision, quartiles.Q1)
	secondQ := fmt.Sprintf("%-*s %.*f\n", offset, "50%", precision, quartiles.Q2)
	thirdQ := fmt.Sprintf("%-*s %.*f\n", offset, "75%", precision, quartiles.Q3)
	max := fmt.Sprintf("%-*s %.*f\n", offset, "max", precision, vals.max())
	return fmt.Sprint(len, valid, null, mean, min, firstQ, secondQ, thirdQ, max)
}

// Constructor Functions
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
