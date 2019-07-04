package dataframe

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// Row returns information about the values and index labels in this row but panics if an out-of-range position is provided.
func (df *DataFrame) Row(position int) Row {
	vals := make([]interface{}, df.NumCols())
	nulls := make([]bool, df.NumCols())
	types := make([]options.DataType, df.NumCols())
	for m := 0; m < df.NumCols(); m++ {
		elem := df.vals[m].Values.Element(position)
		vals[m] = elem.Value
		nulls[m] = elem.Null
		types[m] = df.vals[m].DataType
	}
	idxElems := df.index.Elements(position)
	return Row{Values: vals, Nulls: nulls, ValueTypes: types, Labels: idxElems.Labels, LabelTypes: idxElems.DataTypes}
}

// SelectLabel returns the integer location of the first row in index level 0 with the supplied label, or -1 if the label does not exist.
func (df *DataFrame) SelectLabel(label string) int {
	if df.IndexLevels() == 0 {
		if options.GetLogWarnings() {
			log.Println("DataFrame.SelectLabel(): index has no levels")
		}
		return -1
	}
	val, ok := df.index.Levels[0].LabelMap[label]
	if !ok {
		if options.GetLogWarnings() {
			log.Printf("DataFrame.SelectLabel(): %v not in label map\n", label)
		}
		return -1
	}
	return val[0]
}

// SelectLabels returns the integer locations of all rows with the supplied labels within the supplied level.
// If an error is encountered, returns a new slice of 0 length.
func (df *DataFrame) SelectLabels(labels []string, level int) []int {
	empty := make([]int, 0)
	err := df.ensureIndexLevelPositions([]int{level})
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("DataFrame.SelectLabels(): %v", err)
		}
		return empty
	}
	include := make([]int, 0)
	for _, label := range labels {
		val, ok := df.index.Levels[level].LabelMap[label]
		if !ok {
			if options.GetLogWarnings() {
				log.Printf("DataFrame.SelectLabels(): %v not in label map", label)
			}
			return empty
		}
		include = append(include, val...)
	}
	return include
}

// Col returns the first Series with the specified column label at column level 0.
func (df *DataFrame) Col(label string) *series.Series {
	colPos, ok := df.cols.Levels[0].LabelMap[label]
	if !ok {
		if options.GetLogWarnings() {
			log.Printf("df.Col(): invalid column label: %v not in labels", label)
		}
		s, _ := series.New(nil)
		return s
	}
	return df.hydrateSeries(colPos[0])
}

// subsetRows subsets a DataFrame to include only index items and values at the row positions supplied and modifies the DataFrame in place.
func (ip InPlace) subsetRows(positions []int) {
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values = ip.df.vals[m].Values.Subset(positions)
	}

	ip.df.index.Subset(positions)
}

// subsetRows subsets a DataFrame to include only index items and values at the row positions supplied and modifies the DataFrame in place.
// For use in internal functions that do not expect en error, such as GroupBy.
func (df *DataFrame) subsetRows(positions []int) *DataFrame {
	df = df.Copy()
	df.InPlace.subsetRows(positions)
	return df
}

// SubsetRows subsets a DataFrame to include only the rows at supplied integer positions and modifies the DataFrame in place.
func (ip InPlace) SubsetRows(rowPositions []int) error {
	if len(rowPositions) == 0 {
		return fmt.Errorf("dataframe.SubsetRows(): no valid rows provided")
	}
	if err := ip.df.ensureRowPositions(rowPositions); err != nil {
		return fmt.Errorf("dataframe.SubsetRows(): %v", err)
	}

	ip.subsetRows(rowPositions)
	return nil
}

// SubsetRows subsets a DataFrame to include only the rows at supplied integer positions and returns a new DataFrame.
func (df *DataFrame) SubsetRows(rowPositions []int) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.SubsetRows(rowPositions)
	return df, err
}

// subsetCols subsets a DataFrame to include only columns at the column positions supplied and modifies the DataFrame in place.
func (ip InPlace) subsetCols(positions []int) {
	vals := make([]values.Container, len(positions))
	for i, pos := range positions {
		vals[i] = ip.df.vals[pos]
	}
	ip.df.vals = vals
	ip.df.cols.Subset(positions)
}

// subsetCols subsets a DataFrame to include only index items and values at the columns positions supplied and returns a copy of the DataFrame.
// For use in internal functions that do not expect en error.
func (df *DataFrame) subsetCols(positions []int) *DataFrame {
	df = df.Copy()
	df.InPlace.subsetCols(positions)
	return df
}

// SubsetColumns subsets a DataFrame to include only the columns at supplied integer positions and modifies the DataFrame in place.
func (ip InPlace) SubsetColumns(columnPositions []int) error {
	if len(columnPositions) == 0 {
		return fmt.Errorf("dataframe.SubsetColumns(): no valid columns provided")
	}

	if err := ip.df.ensureColumnPositions(columnPositions); err != nil {
		return fmt.Errorf("dataframe.SubsetColumns(): %v", err)
	}

	ip.subsetCols(columnPositions)

	return nil
}

// SubsetColumns subsets a DataFrame to include only the columns at supplied integer positions and returns a new DataFrame.
func (df *DataFrame) SubsetColumns(columnPositions []int) (*DataFrame, error) {
	df = df.Copy()
	err := df.InPlace.SubsetColumns(columnPositions)
	return df, err
}
