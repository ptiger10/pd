package series

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/ptiger10/pd/internal/values"
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
		unique := fmt.Sprint(len(s.Unique()))
		values = []string{length, valid, null, unique}
		idx = []string{"len", "valid", "null", "unique"}

	case options.Bool:
		precision := options.GetDisplayFloatPrecision()
		sum := fmt.Sprintf("%.*f", precision, s.Sum())
		mean := fmt.Sprintf("%.*f", precision, s.Mean())
		values = []string{length, valid, null, sum, mean}
		idx = []string{"len", "valid", "null", "sum", "mean"}

	case options.DateTime:
		unique := fmt.Sprint(len(s.Unique()))
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
			header += strings.Repeat(" ", values.GetDisplayIndexWhitespaceBuffer())
		}
	}
	// omit header line if empty
	if strings.TrimSpace((header)) != "" {
		printer += header + "\n"
	}

	// [END header row]

	// [START rows]
	prior := make(map[int]string)
	var excludeRows []int
	if s.Len() >= options.GetDisplayMaxRows() {
		half := (options.GetDisplayMaxRows() / 2)
		if options.GetDisplayMaxRows()%2 != 0 {
			excludeRows = values.MakeIntRange(half+1, s.Len()-half)
		} else {
			excludeRows = values.MakeIntRange(half, s.Len()-half)
		}
	}
	var counter int
	for i := 0; i < s.Len(); i++ {
		if excludeRows != nil && counter < len(excludeRows) && i == excludeRows[counter] {
			if counter == 0 {
				printer += "...\n"
			}
			counter++
			continue
		}
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
				buffer = strings.Repeat(" ", values.GetDisplayIndexWhitespaceBuffer())
			}
			// skip repeated label values
			if prior[j] == idx && !options.GetDisplayRepeatedLabels() {
				skip = true
				idx = ""
			}

			// elide index string if longer than the max allowable width
			if padding >= options.GetDisplayMaxWidth() {
				padding = options.GetDisplayMaxWidth()
			}
			if len(idx) >= options.GetDisplayMaxWidth() {
				idx = idx[:options.GetDisplayMaxWidth()-3] + "..."
			}
			printStr := fmt.Sprintf("%*v", padding, idx)

			newLine += printStr + buffer

			// set prior row value
			if !skip {
				prior[j] = idx
			}
		}

		// [END index printer]

		// [START value printer]
		// add buffer at beginning
		newLine += strings.Repeat(" ", values.GetDisplayValuesWhitespaceBuffer())

		var valStr string
		if s.datatype == options.DateTime {
			v, ok := elem.Value.(time.Time)
			if ok {
				valStr = v.Format(options.GetDisplayTimeFormat())
			} else {
				valStr = fmt.Sprint(elem.Value)
			}
		} else if s.datatype == options.Float64 {
			v, ok := elem.Value.(float64)
			if ok {
				valStr = fmt.Sprintf("%.*f", options.GetDisplayFloatPrecision(), v)
			} else {
				valStr = fmt.Sprint(elem.Value)
			}
		} else {
			valStr = fmt.Sprint(elem.Value)
		}

		padding := s.MaxWidth()
		if padding >= options.GetDisplayMaxWidth() {
			padding = options.GetDisplayMaxWidth()
		}
		if len(valStr) >= options.GetDisplayMaxWidth() {
			valStr = valStr[:options.GetDisplayMaxWidth()-3] + "..."
		}
		newLine += fmt.Sprintf("%*v", padding, valStr)
		// Concatenate line onto printer string
		printer += fmt.Sprintln(newLine)
	}
	// [END rows]
	if s.datatype != options.None || s.name != "" {
		printer += "\n"
	}
	if s.datatype != options.None {
		printer += fmt.Sprintf("datatype: %s\n", s.datatype)
	}
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
		printStr += fmt.Sprintf("%10v:%v%v\n", pair[0], strings.Repeat(" ", values.GetDisplayElementWhitespaceBuffer()), pair[1])
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

// Values returns all the values (including null values) in the Series as an interface slice.
func (s *Series) Values() []interface{} {
	var ret []interface{}
	for i := 0; i < s.Len(); i++ {
		ret = append(ret, s.values.Element(i).Value)
	}
	return ret
}

// MaxWidth returns the max number of characters in any value in the Series.
// For use in printing Series.
func (s *Series) MaxWidth() int {
	var max int
	for _, v := range s.Values() {
		var length int
		if s.datatype == options.DateTime {
			if val, ok := v.(time.Time); ok {
				length = len(val.Format(options.GetDisplayTimeFormat()))
			} else {
				length = len(fmt.Sprint(v))
			}
		} else if s.datatype == options.Float64 {
			if val, ok := v.(float64); ok {
				length = len(fmt.Sprintf("%.*f", options.GetDisplayFloatPrecision(), val))
			} else {
				length = len(fmt.Sprint(v))
			}
		} else {
			length = len(fmt.Sprint(v))
		}
		if length > max {
			max = length
		}
	}
	return max
}

// Name returns the Series' name.
func (s Series) Name() string {
	return s.name
}

// Unique returns a de-duplicated list of all non-null values (as []string) that appear in the Series.
func (s *Series) Unique() []string {
	ret := make([]string, 0)
	valid, _ := s.Subset(s.valid())
	counter := valid.ValueCounts()
	for k := range counter {
		ret = append(ret, k)
	}
	return ret
}

// ValueCounts returns a map of non-null value labels to number of occurrences in the Series.
func (s *Series) ValueCounts() map[string]int {
	vals := s.DropNull().Values()
	counter := make(map[string]int)
	for _, val := range vals {
		counter[fmt.Sprint(val)]++
	}
	return counter
}

// Earliest returns the earliest non-null time.Time{} in the Series.
// If applied to anything other than dateTime, return time.Time{}.
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

// Latest returns the latest non-null time.Time{} in the Series.
// If applied to anything other than dateTime, return time.Time{}.
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

// [START ensure methods]

// appropriate for numeric data only ([]float64 or []int64)
func ensureFloatFromNumerics(vals interface{}) []float64 {
	var data []float64
	if ints, ok := vals.([]int64); ok {
		data = convertIntToFloat(ints)
	} else if floats, ok := vals.([]float64); ok {
		data = floats
	} else {
		log.Printf("Internal error: ensureFloatFromNumerics has received an unallowable value: %v", vals)
		return nil
	}
	return data
}

func convertIntToFloat(vals []int64) []float64 {
	var ret []float64
	for _, val := range vals {
		ret = append(ret, float64(val))
	}
	return ret
}

func ensureBools(vals interface{}) []bool {
	if bools, ok := vals.([]bool); ok {
		return bools
	}
	log.Printf("Internal error: ensureBools has received an unallowable value: %v", vals)
	return nil
}

func ensureDateTime(vals interface{}) []time.Time {
	if datetime, ok := vals.([]time.Time); ok {
		return datetime
	}
	log.Printf("Internal error: ensureDateTime has received an unallowable value: %v", vals)
	return nil
}

// returns an error if any index levels have different lengths
// or if there is a mismatch between the number of values and index items
func (s *Series) ensureAlignment() error {
	if err := s.index.Aligned(); err != nil {
		return fmt.Errorf("out of alignment: %v", err)
	}
	if labels := s.index.Levels[0].Len(); s.Len() != labels {
		return fmt.Errorf("out of alignment: series must have same number of values as index labels (%d != %d)", s.Len(), labels)
	}
	return nil
}

// returns an error if any row position does not exist
func (s *Series) ensureRowPositions(positions []int) error {
	if len(positions) == 0 {
		return fmt.Errorf("no valid rows")
	}

	len := s.Len()
	for _, pos := range positions {
		if pos >= len {
			return fmt.Errorf("invalid position: %d (max %v)", pos, len-1)
		}
	}
	return nil
}

// returns an error if any level position does not exist
func (s *Series) ensureLevelPositions(positions []int) error {
	for _, pos := range positions {
		len := s.index.NumLevels()
		if pos >= len {
			return fmt.Errorf("invalid position: %d (max %v)", pos, len-1)
		}
	}
	return nil
}

// [END ensure methods]
