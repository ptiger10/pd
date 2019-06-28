package dataframe

import (
	"fmt"
	"strings"
	"time"

	"github.com/ptiger10/pd/options"
)

func (df *DataFrame) String() string {
	if Equal(df, newEmptyDataFrame()) {
		return "{Empty DataFrame}"
	}
	return df.print()
}

// printer for DataFrame index, values, and columns.
// Format (optional):
// (indexName)
// (multiLevelColumnName) (multiLevelColumns)
// (columnName) columns
// index value
// Syntax:
// i -> values
// j -> index or column levels
// k -> columns
func (df *DataFrame) print() string {
	numLevels := df.IndexLevels()
	var indexNameRow string
	var printer string
	// [START index name row]
	maxIndexWidths := df.index.MaxWidths()
	for j := 0; j < numLevels; j++ {
		name := df.index.Levels[j].Name
		padding := maxIndexWidths[j]
		indexNameRow += fmt.Sprintf("%*v", padding, name)
		if j != numLevels-1 {
			// add buffer to all index levels except the last
			indexNameRow += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
		}
	}
	if !df.index.Unnamed() {
		printer += indexNameRow + "\n"
	}
	// [END index name row]

	// [START column rows]
	var colLevelRows []string
	excl := df.makeExclusionsTable()
	maxColWidths := df.maxColWidths(excl)
	for j := 0; j < df.cols.NumLevels(); j++ {
		colLevelRow := strings.Repeat(" ", len(indexNameRow)+options.GetDisplayValuesWhitespaceBuffer())
		colName := df.cols.Levels[j].Name
		namePadding := df.cols.MaxNameWidth()
		colLevelRow += fmt.Sprintf("%*v", namePadding, colName)
		if colName != "" {
			colLevelRow += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
		}
		var prior string
		for k := 0; k < df.NumCols(); k++ {
			colLabel := fmt.Sprint(df.cols.Levels[j].Labels[k])
			if colLabel == prior && !options.GetDisplayRepeatedLabels() {
				colLabel = ""
				excl[j][k] = true
				maxColWidths = df.maxColWidths(excl)
			}
			colPadding := maxColWidths[k]
			colLevelRow += fmt.Sprintf("%*v", colPadding, colLabel)
			if colLabel != "" {
				prior = colLabel
			}
			if k != df.NumCols()-1 {
				// add buffer to all columns except the last
				colLevelRow += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
			} else {
				colLevelRow = strings.TrimRight(colLevelRow, " ")
			}
		}
		colLevelRows = append(colLevelRows, colLevelRow)
	}
	printer += strings.Join(colLevelRows, "\n") + "\n"
	// [END column rows]

	// [START rows]
	prior := make(map[int]string)
	for i := 0; i < df.Len(); i++ {
		idxElems := df.index.Elements(i)
		var newLine string

		// [START index printer]
		for j := 0; j < numLevels; j++ {
			var skip bool
			var buffer string
			padding := maxIndexWidths[j]
			idx := fmt.Sprint(idxElems.Labels[j])
			if j != numLevels-1 {
				// add buffer to all index levels except the last
				buffer = strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
				// skip repeated label values if this is not the last index level
				if prior[j] == idx && !options.GetDisplayRepeatedLabels() {
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
		newLine += strings.Repeat(" ", options.GetDisplayValuesWhitespaceBuffer()+df.cols.MaxNameWidth())
		if df.cols.MaxNameWidth() != 0 {
			newLine += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
		}

		var valStrs string
		for k := 0; k < df.NumCols(); k++ {
			valElem := df.s[k].Element(i)
			var valStr string
			padding := maxColWidths[k]
			if df.s[k].DataType() == string(options.DateTime) {
				valStr = valElem.Value.(time.Time).Format(options.GetDisplayTimeFormat())
			} else {
				valStr = fmt.Sprintf("%*v", padding, valElem.Value)
			}
			if padding == options.GetDisplayMaxWidth() {
				valStr = valStr[:len(valStr)-4] + "..."
			}
			if k != df.NumCols()-1 {
				// add buffer to all columns except the last
				valStr += strings.Repeat(" ", options.GetDisplayIndexWhitespaceBuffer())
			}
			valStrs += valStr
		}
		newLine += valStrs
		// null string values must not return any trailing whitespace
		newLine = strings.TrimRight(newLine, " ")
		// Concatenate line onto printer string
		printer += fmt.Sprintln(newLine)
	}
	// [END rows]
	if df.dataType() != "mixed" {
		printer += fmt.Sprintf("datatype: %s\n", df.dataType())
	}

	if df.name != "" {
		printer += fmt.Sprintf("name: %s\n", df.name)
	}
	return printer
}
