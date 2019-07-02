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

func (df *DataFrame) selectByRows(rowPositions []int) (*DataFrame, error) {
	if err := df.ensureAlignment(); err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe internal alignment error: %v", err)
	}
	idx := df.index.Subset(rowPositions)
	var subsetVals []values.Container
	for i := 0; i < df.NumCols(); i++ {
		vals := df.vals[i].Values.Subset(rowPositions)
		subsetVals = append(subsetVals, values.Container{Values: vals, DataType: df.vals[i].DataType})
	}
	df = newFromComponents(subsetVals, idx, df.cols, df.name)
	return df, nil
}

func (df *DataFrame) selectByCols(colPositions []int) (*DataFrame, error) {
	if err := df.ensureAlignment(); err != nil {
		return df, fmt.Errorf("dataframe internal alignment error: %v", err)
	}
	var subsetVals []values.Container
	for _, pos := range colPositions {
		if pos >= df.NumCols() {
			return nil, fmt.Errorf("dataframe.selectByCols(): invalid col position %d (max: %d)", pos, df.NumCols()-1)
		}
		subsetVals = append(subsetVals, df.vals[pos])
	}
	columnsSlice, err := df.cols.Subset(colPositions)
	if err != nil {
		return nil, fmt.Errorf("dataframe.selectByCols(): %v", err)
	}

	df = newFromComponents(subsetVals, df.index, columnsSlice, df.name)
	return df, nil
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
	df, _ = df.selectByCols(colPos)
	return df.hydrateSeries(0)
}

// Subset a DataFrame to include only the rows at supplied integer positions.
func (df *DataFrame) Subset(rowPositions []int) (*DataFrame, error) {
	if len(rowPositions) == 0 {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.Subset(): no valid rows provided")
	}

	sub, err := df.selectByRows(rowPositions)
	if err != nil {
		return newEmptyDataFrame(), fmt.Errorf("dataframe.Subset(): %v", err)
	}
	return sub, nil
}

// subsetRows subsets a DataFrame to include only index items and values at the row positions supplied and modifies the DataFrame in place.
func (ip InPlace) subsetRows(positions []int) {
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values = ip.df.vals[m].Values.Subset(positions)
	}

	ip.df.index = ip.df.index.Subset(positions)
}

// subsetRows subsets a DataFrame to include only index items and values at the row positions supplied and returns a new DataFrame.
func (df *DataFrame) subsetRows(positions []int) *DataFrame {
	df = df.Copy()
	df.InPlace.subsetRows(positions)
	return df
}
