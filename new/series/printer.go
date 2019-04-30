package series

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ptiger10/pd/new/options"
)

func (s Series) String() string {
	switch s.Kind {
	// case DateTime:
	// 	var printer string
	// 	vals := s.Values.(dateTimeValues)
	// 	for _, val := range vals {
	// 		printer += fmt.Sprintln(val.v.Format("01/02/2006"))
	// 	}
	// 	return printer
	default:
		return s.print()
	}
}

// expects to receive a slice of typed value structs (eg pd.floatValues)
// each struct must contain a boolean field called "v"
func (s Series) print() string {

	vals := reflect.ValueOf(s.Values)
	levels := reflect.ValueOf(s.Index.Levels)
	numLevels := len(s.Index.Levels)
	var printer string
	var header string
	// header row
	for j := 0; j < numLevels; j++ {
		name := s.Index.Levels[j].Name
		padding := s.Index.Levels[j].Longest
		header += fmt.Sprintf("%*v", padding, name)
		if j != numLevels-1 {
			// add buffer to all index levels except the last
			header += strings.Repeat(" ", options.DisplayIndexWhitespaceBuffer)
		}
	}
	if header != "" {
		header = fmt.Sprintln(header)
	}
	printer += header
	prior := make(map[int]string)
	for i := 0; i < s.Len(); i++ {

		var newLine string
		for j := 0; j < numLevels; j++ {
			var skip bool
			var buffer string
			padding := s.Index.Levels[j].Longest
			// s.Index.Levels[j].Labels[i]
			idxVal := levels.Index(j).FieldByName("Labels").Elem().Index(i)
			idxStr := fmt.Sprint(idxVal)
			if j != numLevels-1 {
				// add buffer to all index levels except the last
				buffer = strings.Repeat(" ", options.DisplayIndexWhitespaceBuffer)
				// skip repeated label values if this is not the last index level
				if prior[j] == idxStr {
					skip = true
					idxStr = ""
				}
			}

			printStr := fmt.Sprintf("%*v", padding, idxStr)
			// elide index string if longer than the max allowable width
			if padding == options.DisplayIndexMaxWidth {
				printStr = printStr[:len(printStr)-4] + "..."
			}

			newLine += printStr + buffer

			// set prior row value for each index level except the last
			if j != numLevels-1 && !skip {
				prior[j] = idxStr
			}
		}
		// s.Values[i].V
		valStr := fmt.Sprint(vals.Index(i).FieldByName("V"))
		// add buffer at beginning
		val := strings.Repeat(" ", options.DisplayValuesWhitespaceBuffer) + valStr
		// null string values must not return any trailing whitespace
		if valStr == "" {
			val = strings.TrimSpace(val)
		}
		newLine += val
		// Concatenate line onto printer string
		printer += fmt.Sprintln(newLine)
	}
	printer += fmt.Sprintf("kind: %s\n", s.Kind)
	if s.Name != "" {
		printer += fmt.Sprintf("name: %s\n", s.Name)
	}
	return printer
}
