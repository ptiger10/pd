package dataframe

import (
	"fmt"
	"strings"
	"time"

	"github.com/ptiger10/pd/options"
)

func (df *DataFrame) String() string {
	if df == nil {
		return "DataFrame{}"
	}
	return df.print()
}

// printer for DataFrame index, values, and columns
func (df *DataFrame) print() string {
	numLevels := df.Levels()
	var header string
	var printer string
	// [START header row]
	for j := 0; j < numLevels; j++ {
		name := df.index.Levels[j].Name
		padding := df.index.Levels[j].MaxWidth()
		header += fmt.Sprintf("%*v", padding, name)
		if j != numLevels-1 {
			// add buffer to all index levels except the last
			header += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
		}
	}
	// add buffer at beginning
	header += strings.Repeat(" ", options.GetDisplayValuesWhitespaceBuffer())
	// handle columns
	for k := 0; k < df.Cols(); k++ {
		colStr := df.cols[k]
		padding := df.s[k].MaxWidth()
		if padding == options.GetDisplayMaxWidth() {
			colStr = colStr[:len(colStr)-4] + "..."
		}
		header += fmt.Sprintf("%*v", padding, colStr)
		if k != df.Cols()-1 {
			// add buffer to all columns except the last
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
	for i := 0; i < df.Len(); i++ {
		idxElems := df.index.Elements(i)
		var newLine string

		// [START index printer]
		for j := 0; j < numLevels; j++ {
			var skip bool
			var buffer string
			padding := df.index.Levels[j].MaxWidth()
			idx := fmt.Sprint(idxElems.Labels[j])
			if j != numLevels-1 {
				// add buffer to all index levels except the last
				buffer = strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
				// skip repeated label values if this is not the last index level
				if prior[j] == idx {
					skip = true
					idx = ""
				}
			}

			idxStr := fmt.Sprintf("%*v", padding, idx)
			// elide index string if longer than the max allowable width
			if padding == options.GetDisplayMaxWidth() {
				idxStr = idxStr[:len(idxStr)-4] + "..."
			}

			newLine += idxStr + buffer

			// set prior row value for each index level except the last
			if j != numLevels-1 && !skip {
				prior[j] = idx
			}
		}

		// [END index printer]

		// [START value printer]
		var valStrs string
		for k := 0; k < df.Cols(); k++ {
			valElem := df.s[k].Element(i)
			var valStr string
			padding := df.s[k].MaxWidth()
			if df.s[k].DataType() == string(options.DateTime) {
				valStr = valElem.Value.(time.Time).Format(options.GetDisplayTimeFormat())
			} else {
				valStr = fmt.Sprintf("%*v", padding, valElem.Value)
			}
			if padding == options.GetDisplayMaxWidth() {
				valStr = valStr[:len(valStr)-4] + "..."
			}
			if k != df.Cols()-1 {
				// add buffer to all columns except the last
				valStr += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
			}
			valStrs += valStr
		}

		// add buffer at beginning
		val := strings.Repeat(" ", options.GetDisplayValuesWhitespaceBuffer()) + valStrs
		// null string values must not return any trailing whitespace
		if valStrs == "" {
			val = strings.TrimSpace(val)
		}
		newLine += val
		// Concatenate line onto printer string
		printer += fmt.Sprintln(newLine)
	}
	// printer += fmt.Sprintf("datatype: %s\n", df.DataType())
	// [END rows]

	if df.Name != "" {
		printer += fmt.Sprintf("name: %s\n", df.Name)
	}
	return printer
}
