package dataframe

import (
	"math"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// A DataFrame is a 2D collection of one or more Series.
type DataFrame struct {
	Name     string
	DataType options.DataType
	s        []*series.Series
	cols     []string
	index    index.Index
}

// Sum all numerical or boolean columns.
func (df *DataFrame) Sum() *series.Series {
	ret := new(series.Series)
	for _, s := range df.s {
		if sum := s.Sum(); !math.IsNaN(sum) {
			newSum := series.MustNew(sum, series.Idx(s.Name))
			ret.Join(newSum)
		}
	}
	return ret
}

// func (df *DataFrame) String() string {
// 	if df == nil {
// 		return "DataFrame{}"
// 	}
// 	return df.s[0].
// }

// // expects to receive a slice of typed value structs (eg values.float64Values)
// func (df *DataFrame) print() string {
// 	numLevels := len(df.s[0].values.Levels)
// 	var header string
// 	var printer string
// 	// [START header row]
// 	for j := 0; j < numLevels; j++ {
// 		name := s.Index.Levels[j].Name
// 		padding := s.Index.Levels[j].Longest
// 		header += fmt.Sprintf("%*v", padding, name)
// 		if j != numLevels-1 {
// 			// add buffer to all index levels except the last
// 			header += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
// 		}
// 	}
// 	// omit header line if empty
// 	if strings.TrimSpace((header)) != "" {
// 		printer += header + "\n"
// 	}

// 	// [END header row]

// 	// [START rows]
// 	prior := make(map[int]string)
// 	for i := 0; i < s.Len(); i++ {
// 		elem := s.Element(i)
// 		var newLine string

// 		// [START index printer]
// 		for j := 0; j < numLevels; j++ {
// 			var skip bool
// 			var buffer string
// 			padding := s.Index.Levels[j].Longest
// 			idx := fmt.Sprint(elem.Labels[j])
// 			if j != numLevels-1 {
// 				// add buffer to all index levels except the last
// 				buffer = strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
// 				// skip repeated label values if this is not the last index level
// 				if prior[j] == idx {
// 					skip = true
// 					idx = ""
// 				}
// 			}

// 			printStr := fmt.Sprintf("%*v", padding, idx)
// 			// elide index string if longer than the max allowable width
// 			if padding == options.GetDisplayIndexMaxWidth() {
// 				printStr = printStr[:len(printStr)-4] + "..."
// 			}

// 			newLine += printStr + buffer

// 			// set prior row value for each index level except the last
// 			if j != numLevels-1 && !skip {
// 				prior[j] = idx
// 			}
// 		}

// 		// [END index printer]

// 		// [START value printer]
// 		var valStr string
// 		if s.datatype == options.DateTime {
// 			valStr = elem.Value.(time.Time).Format(options.GetDisplayTimeFormat())
// 		} else {
// 			valStr = fmt.Sprint(elem.Value)
// 		}

// 		// add buffer at beginning
// 		val := strings.Repeat(" ", options.GetDisplayValuesWhitespaceBuffer()) + valStr
// 		// null string values must not return any trailing whitespace
// 		if valStr == "" {
// 			val = strings.TrimSpace(val)
// 		}
// 		newLine += val
// 		// Concatenate line onto printer string
// 		printer += fmt.Sprintln(newLine)
// 	}
// 	printer += fmt.Sprintf("datatype: %s\n", s.datatype)
// 	// [END rows]

// 	if s.Name != "" {
// 		printer += fmt.Sprintf("name: %s\n", s.Name)
// 	}
// 	return printer
// }
