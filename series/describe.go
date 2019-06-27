package series

import (
	"fmt"
	"log"
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
		return "Series{}"
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

// Element returns information about the value and index labels at this position but panics if an out-of-range position is provided.
func (s *Series) Element(position int) Element {
	elem := s.values.Element(position)
	idxElems := s.index.Elements(position)
	return Element{elem.Value, elem.Null, idxElems.Labels, idxElems.DataTypes}
}

// DataType is the data type of the Series' values. Mimics reflect.Type with the addition of time.Time as DateTime.
func (s *Series) DataType() string {
	return fmt.Sprint(s.datatype)
}

// selectByRows copies a Series then subsets it to include only index items and values at the positions supplied
func (s *Series) selectByRows(positions []int) (*Series, error) {
	if err := s.ensureAlignment(); err != nil {
		return newEmptySeries(), fmt.Errorf("series internal alignment error: %v", err)
	}
	if positions == nil {
		return newEmptySeries(), nil
	}
	if err := s.ensureRowPositions(positions); err != nil {
		return newEmptySeries(), fmt.Errorf("s.selectByRows(): %v", err)
	}

	s = s.Copy()
	s.values = s.values.Subset(positions)
	s.index = s.index.Subset(positions)
	return s, nil
}

func (s *Series) mustSelectRows(positions []int) *Series {
	s, err := s.selectByRows(positions)
	if err != nil {
		log.Printf("Internal error: %v\n", err)
		return newEmptySeries()
	}
	return s
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

// Subset returns a subset of a Series based on the supplied integer positions.
func (s *Series) Subset(rowPositions []int) (*Series, error) {
	if len(rowPositions) == 0 {
		return newEmptySeries(), fmt.Errorf("series.Subset(): no valid rows provided")
	}

	sub, err := s.selectByRows(rowPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("series.Subset(): %v", err)
	}
	return sub, nil
}

// At returns the value at a single integer position, bur returns nil if the position is out of range.
func (s *Series) At(position int) interface{} {
	if position >= s.Len() {
		if options.GetLogWarnings() {
			log.Printf("s.Index.At(): invalid position: %d (max: %v)", position, s.Len()-1)
		}
		return nil
	}
	elem := s.Element(position)
	return elem.Value
}

// [END Series methods]

// [START Selection]

// XS returns a new Series with only the rows and index levels at the specified positions.
func (s *Series) XS(rowPositions []int, levelPositions []int) (*Series, error) {
	var err error
	s, err = s.Subset(rowPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.XS() rows: %v", err)
	}
	s, err = s.Index.Subset(levelPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.XS() index levels: %v", err)
	}
	return s, nil
}

// SelectLabel returns the integer location of the first row with the specified index label at the specified level, or -1 if the label does not exist.
func (s *Series) SelectLabel(label string, level int) int {
	val, ok := s.index.Levels[0].LabelMap[label]
	if !ok {
		if options.GetLogWarnings() {
			log.Printf("Series.SelectLabel(): %v not in label map", label)
		}
		return -1
	}
	return val[0]
}

// SelectLabels returns the integer locations of all rows with the specified index label at the specified level.
func (s *Series) SelectLabels(labels []string, level int) []int {
	empty := make([]int, 0)
	include := make([]int, 0)
	for _, label := range labels {
		val, ok := s.index.Levels[0].LabelMap[label]
		if !ok {
			if options.GetLogWarnings() {
				log.Printf("s.Index.ByLabels(): %v not in label map", label)
			}
			return empty
		}
		include = append(include, val...)
	}
	return include
}

// [START string/datetime description methods]

// ValueCounts returns a map of non-null value labels to number of occurrences in the Series.
//
// Applies to: All
func (s *Series) ValueCounts() map[string]int {
	valid, _ := s.selectByRows(s.valid())
	if valid == nil {
		return nil
	}
	vals := valid.all()
	counter := make(map[string]int)
	for _, val := range vals {
		counter[fmt.Sprint(val)]++
	}
	return counter
}

// UniqueVals returns a de-duplicated list of all non-null values (as []string) that appear in the Series.
//
// Applies to: All
func (s *Series) UniqueVals() []string {
	var ret []string
	valid, _ := s.Subset(s.valid())
	counter := valid.ValueCounts()
	if counter == nil {
		return []string{}
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
