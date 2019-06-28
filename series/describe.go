package series

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ptiger10/pd/options"
)

// Describe the key details of the Series.
func (s *Series) Describe() {
	var values interface{}
	var idx interface{}
	// dt preserves the datatype of the original Series
	dt := s.datatype

	// shared data
	l := s.Len()
	valids := len(s.valid())
	nulls := len(s.null())
	length := fmt.Sprint(l)
	valid := fmt.Sprint(valids)
	null := fmt.Sprint(nulls)
	// type-specific data
	switch s.datatype {
	case options.Float64, options.Int64:
		precision := options.GetDisplayFloatPrecision()
		mean := fmt.Sprintf("%.*f", precision, s.Mean())
		min := fmt.Sprintf("%.*f", precision, s.Min())
		q1 := fmt.Sprintf("%.*f", precision, s.Quartile(1))
		q2 := fmt.Sprintf("%.*f", precision, s.Quartile(2))
		q3 := fmt.Sprintf("%.*f", precision, s.Quartile(3))
		max := fmt.Sprintf("%.*f", precision, s.Max())

		values = []string{length, valid, null, mean, min, q1, q2, q3, max}
		idx = []string{"len", "valid", "null", "mean", "min", "25%", "50%", "75%", "max"}

	case options.String:
		unique := fmt.Sprint(len(s.UniqueVals()))
		values = []string{length, valid, null, unique}
		idx = []string{"len", "valid", "null", "unique"}

	case options.Bool:
		precision := options.GetDisplayFloatPrecision()
		sum := fmt.Sprintf("%.*f", precision, s.Sum())
		mean := fmt.Sprintf("%.*f", precision, s.Mean())
		values = []string{length, valid, null, sum, mean}
		idx = []string{"len", "valid", "null", "sum", "mean"}

	case options.DateTime:
		unique := fmt.Sprint(len(s.UniqueVals()))
		earliest := fmt.Sprint(s.Earliest())
		latest := fmt.Sprint(s.Latest())
		values = []string{length, valid, null, unique, earliest, latest}
		idx = []string{"len", "valid", "null", "unique", "earliest", "latest"}

	// Interface or None
	default:
		values = []string{length, valid, null}
		idx = []string{"len", "valid", "null"}
	}

	// duck errors because constructor called internally
	s = MustNew(values, Config{Name: s.name, Index: idx})
	s.datatype = dt
	fmt.Println(s)
	return

}

func (s *Series) String() string {
	if Equal(s, newEmptySeries()) {
		return "{Empty Series}"
	}
	return s.print()
}

// expects to receive a slice of typed value structs (eg values.float64Values)
func (s *Series) print() string {
	numLevels := len(s.index.Levels)
	var header string
	var printer string
	// [START header row]
	maxIndexWidths := s.index.MaxWidths()
	for j := 0; j < numLevels; j++ {
		name := s.index.Levels[j].Name
		padding := maxIndexWidths[j]
		header += fmt.Sprintf("%*v", padding, name)
		if j != numLevels-1 {
			// add buffer to all index levels except the last
			header += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
		}
	}
	// omit header line if empty
	if strings.TrimSpace((header)) != "" {
		printer += header + "\n"
	}

	// [END header row]

	// [START rows]
	prior := make(map[int]string)
	for i := 0; i < s.Len(); i++ {
		elem := s.Element(i)
		var newLine string

		// [START index printer]
		for j := 0; j < numLevels; j++ {
			var skip bool
			var buffer string
			padding := maxIndexWidths[j]
			idx := fmt.Sprint(elem.Labels[j])
			if j != numLevels-1 {
				// add buffer to all index levels except the last
				buffer = strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
				// skip repeated label values if this is not the last index level
				if prior[j] == idx {
					skip = true
					idx = ""
				}
			}

			printStr := fmt.Sprintf("%*v", padding, idx)
			// elide index string if longer than the max allowable width
			if padding == options.GetDisplayMaxWidth() {
				printStr = printStr[:len(printStr)-4] + "..."
			}

			newLine += printStr + buffer

			// set prior row value for each index level except the last
			if j != numLevels-1 && !skip {
				prior[j] = idx
			}
		}

		// [END index printer]

		// [START value printer]
		var valStr string
		if s.datatype == options.DateTime {
			valStr = elem.Value.(time.Time).Format(options.GetDisplayTimeFormat())
		} else {
			valStr = fmt.Sprint(elem.Value)
		}

		// add buffer at beginning
		val := strings.Repeat(" ", options.GetDisplayValuesWhitespaceBuffer()) + valStr
		// null string values must not return any trailing whitespace
		if valStr == "" {
			val = strings.TrimSpace(val)
		}
		newLine += val
		// Concatenate line onto printer string
		printer += fmt.Sprintln(newLine)
	}
	if s.datatype != options.None {
		printer += fmt.Sprintf("datatype: %s\n", s.datatype)
	}
	// [END rows]

	if s.name != "" {
		printer += fmt.Sprintf("name: %s\n", s.name)
	}
	return printer
}

func (el Element) String() string {
	var printStr string
	for _, pair := range [][]interface{}{
		[]interface{}{"Value", el.Value},
		[]interface{}{"Null", el.Null},
		[]interface{}{"Labels", el.Labels},
		[]interface{}{"LabelTypes", el.LabelTypes},
	} {
		// LabelTypes is 10 characters wide, so left padding set to 10
		printStr += fmt.Sprintf("%10v:%v%v\n", pair[0], strings.Repeat(" ", options.GetDisplayElementWhitespaceBuffer()), pair[1])
	}
	return printStr
}

// DataType is the data type of the Series' values. Mimics reflect.Type with the addition of time.Time as DateTime.
func (s *Series) DataType() string {
	return fmt.Sprint(s.datatype)
}

// Equal compares whether two series are equivalent.
func Equal(s1, s2 *Series) bool {
	if !reflect.DeepEqual(s1.values, s2.values) {
		return false
	}
	if !reflect.DeepEqual(s1.index, s2.index) {
		return false
	}
	if s1.name != s2.name {
		return false
	}
	if s1.datatype != s2.datatype {
		return false
	}
	return true
}

// Len returns the number of Elements (i.e., Value/Null pairs) in the Series.
func (s *Series) Len() int {
	if s.values == nil {
		return 0
	}
	return s.values.Len()
}

// valid returns integer positions of valid (i.e., non-null) values in the series.
func (s *Series) valid() []int {
	var ret []int
	for i := 0; i < s.Len(); i++ {
		if !s.values.Element(i).Null {
			ret = append(ret, i)
		}
	}
	return ret
}

// null returns the integer position of all null values in the collection.
func (s *Series) null() []int {
	var ret []int
	for i := 0; i < s.Len(); i++ {
		if s.values.Element(i).Null {
			ret = append(ret, i)
		}
	}
	return ret
}

// all returns only the Value fields for the collection of Value/Null structs as an interface slice.
//
// Caution: This operation excludes the Null field but retains any null values.
func (s *Series) all() []interface{} {
	var ret []interface{}
	for i := 0; i < s.Len(); i++ {
		ret = append(ret, s.values.Element(i).Value)
	}
	return ret
}

// MaxWidth returns the max number of characters in any value in the Series.
// For use in printing Series and DataFrames.
func (s *Series) MaxWidth() int {
	var max int
	for _, v := range s.all() {
		if length := len(fmt.Sprint(v)); length > max {
			max = length
		}
	}
	return max
}

// Name returns the Series' name.
func (s Series) Name() string {
	return s.name
}

// UniqueVals returns a de-duplicated list of all non-null values (as []string) that appear in the Series.
//
// Applies to: All
func (s *Series) UniqueVals() []string {
	ret := make([]string, 0)
	valid, _ := s.Subset(s.valid())
	counter := valid.ValueCounts()
	if counter == nil {
		return ret
	}
	for k := range counter {
		ret = append(ret, k)
	}
	return ret
}

// Earliest returns the earliest non-null time.Time{} in the Series
//
// Applies to: time.Time. If inapplicable, defaults to time.Time{}.
func (s *Series) Earliest() time.Time {
	earliest := time.Time{}
	vals := s.validVals()
	switch s.datatype {
	case options.DateTime:
		data := ensureDateTime(vals)
		for _, t := range data {
			if earliest == (time.Time{}) || t.Before(earliest) {
				earliest = t
			}
		}
		return earliest
	default:
		return earliest

	}
}

// Latest returns the latest non-null time.Time{} in the Series
//
// Applies to: time.Time. If inapplicable, defaults to time.Time{}.
func (s *Series) Latest() time.Time {
	latest := time.Time{}
	vals := s.validVals()
	switch s.datatype {
	case options.DateTime:
		data := ensureDateTime(vals)
		for _, t := range data {
			if latest == (time.Time{}) || t.After(latest) {
				latest = t
			}
		}
		return latest
	default:
		return latest

	}
}

// [END string/datetime description methods]
