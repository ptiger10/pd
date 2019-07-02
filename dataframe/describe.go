package dataframe

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

func (df *DataFrame) String() string {
	if Equal(df, newEmptyDataFrame()) {
		return "{Empty DataFrame}"
	}
	return df.print()
}

// Len returns the number of values in each Series of the DataFrame.
func (df *DataFrame) Len() int {
	if df.vals == nil {
		return 0
	}
	return df.vals[0].Values.Len()
}

// Name returns the DataFrame's name.
func (df *DataFrame) Name() string {
	return df.name
}

// Rename the DataFrame.
func (df *DataFrame) Rename(name string) {
	df.name = name
}

// NumCols returns the number of columns in the DataFrame.
func (df *DataFrame) NumCols() int {
	if df.vals == nil {
		return 0
	}
	return len(df.vals)
}

// IndexLevels returns the number of index levels in the DataFrame.
func (df *DataFrame) IndexLevels() int {
	return df.index.NumLevels()
}

// ColLevels returns the number of column levels in the DataFrame.
func (df *DataFrame) ColLevels() int {
	return df.cols.NumLevels()
}

// printer for DataFrame index, values, and columns.
// Format (optional):
// (indexName)
// (multiLevelColumnName) (multiLevelColumns)
// (columnName) columns
// index value
// For loop syntax:
// i -> values
// j -> index or column levels
// m -> number of columns
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
			indexNameRow += strings.Repeat(" ", values.GetDisplayIndexWhitespaceBuffer())
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
		colLevelRow := strings.Repeat(" ", len(indexNameRow)+values.GetDisplayValuesWhitespaceBuffer())
		colName := df.cols.Levels[j].Name
		namePadding := df.cols.MaxNameWidth()
		colLevelRow += fmt.Sprintf("%*v", namePadding, colName)
		if colName != "" {
			colLevelRow += strings.Repeat(" ", values.GetDisplayIndexWhitespaceBuffer())
		}
		var prior string
		for m := 0; m < df.NumCols(); m++ {
			colLabel := fmt.Sprint(df.cols.Levels[j].Labels[m])
			if colLabel == prior && !options.GetDisplayRepeatedLabels() {
				colLabel = ""
				excl[j][m] = true
				maxColWidths = df.maxColWidths(excl)
			}
			colPadding := maxColWidths[m]
			colLevelRow += fmt.Sprintf("%*v", colPadding, colLabel)
			if colLabel != "" {
				prior = colLabel
			}
			if m != df.NumCols()-1 {
				// add buffer to all columns except the last
				colLevelRow += strings.Repeat(" ", values.GetDisplayIndexWhitespaceBuffer())
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
				buffer = strings.Repeat(" ", values.GetDisplayIndexWhitespaceBuffer())
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
		newLine += strings.Repeat(" ", values.GetDisplayValuesWhitespaceBuffer()+df.cols.MaxNameWidth())
		if df.cols.MaxNameWidth() != 0 {
			newLine += strings.Repeat(" ", values.GetDisplayIndexWhitespaceBuffer())
		}

		var valStrs string
		for m := 0; m < df.NumCols(); m++ {
			valElem := df.vals[m].Values.Element(i)
			var valStr string
			padding := maxColWidths[m]
			if df.vals[m].DataType == options.DateTime {
				valStr = valElem.Value.(time.Time).Format(options.GetDisplayTimeFormat())
			} else {
				valStr = fmt.Sprintf("%*v", padding, valElem.Value)
			}
			if padding == options.GetDisplayMaxWidth() {
				valStr = valStr[:len(valStr)-4] + "..."
			}
			if m != df.NumCols()-1 {
				// add buffer to all columns except the last
				valStr += strings.Repeat(" ", values.GetDisplayIndexWhitespaceBuffer())
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
	if df.dataType() != "mixed" || df.name != "" {
		printer += "\n"
	}
	if df.dataType() != "mixed" {
		printer += fmt.Sprintf("datatype: %s\n", df.dataType())
	}

	if df.name != "" {
		printer += fmt.Sprintf("name: %s\n", df.name)
	}
	return printer
}

// Equal returns true if two dataframes contain equivalent values.
func Equal(df, df2 *DataFrame) bool {
	if df.NumCols() != df2.NumCols() {
		return false
	}
	for i := 0; i < df.NumCols(); i++ {
		if !series.Equal(df.hydrateSeries(i), df2.hydrateSeries(i)) {
			return false
		}
	}
	if !reflect.DeepEqual(df.index, df2.index) {
		return false
	}
	if !reflect.DeepEqual(df.cols, df2.cols) {
		return false
	}
	if df.name != df2.name {
		return false
	}
	return true
}

// DataTypes returns the DataTypes of the Series in the DataFrame.
func (df *DataFrame) DataTypes() *series.Series {
	var types []interface{}
	for _, val := range df.vals {
		types = append(types, val.DataType)
	}
	vals := values.MustCreateValuesFromInterface(types)
	s := series.FromInternalComponents(vals.Values, df.index, options.String, "datatypes")
	return s
}

// dataType is the data type of the DataFrame's values. Mimics reflect.Type with the addition of time.Time as DateTime.
func (df *DataFrame) dataType() string {
	uniqueTypes := df.DataTypes().Unique()
	if len(uniqueTypes) == 1 {
		return string(df.vals[0].DataType)
	}
	return "mixed"
}

// maxWidths is the max characters in each values container of a dataframe.
func (df *DataFrame) maxWidths() []int {
	maxWidths := make([]int, df.NumCols())
	for m := 0; m < df.NumCols(); m++ {
		var max int
		vc := df.vals[m]
		for _, val := range vc.Values.Values() {
			var length int
			if vc.DataType == options.DateTime {
				length = len(val.(time.Time).Format(options.GetDisplayTimeFormat()))
			} else if vc.DataType == options.Float64 {
				length = len(fmt.Sprintf("%.*f", options.GetDisplayFloatPrecision(), val.(float64)))
			} else {
				length = len(fmt.Sprint(val))
			}
			if length > max {
				max = length
			}
		}
		maxWidths[m] = max
	}
	return maxWidths
}

// maxColWidths is the max characters in each column of a dataframe.
// exclusions should mimic the shape of the columns exactly
func (df *DataFrame) maxColWidths(exclusions [][]bool) []int {
	var maxColWidths []int
	if len(exclusions) != df.ColLevels() || len(exclusions) == 0 {
		return nil
	}
	if len(exclusions[0]) != df.NumCols() {
		return nil
	}
	maxValWidths := df.maxWidths()

	for m := 0; m < df.NumCols(); m++ {
		max := maxValWidths[m]
		for j := 0; j < df.ColLevels(); j++ {
			if !exclusions[j][m] {
				if length := len(fmt.Sprint(df.cols.Levels[j].Labels[m])); length > max {
					max = length
				}
			}
		}
		maxColWidths = append(maxColWidths, max)
	}
	return maxColWidths
}

// for use in  printing dataframe columns
func (df *DataFrame) makeExclusionsTable() [][]bool {
	table := make([][]bool, df.ColLevels())
	for row := range table {
		table[row] = make([]bool, df.NumCols())
	}
	return table
}
