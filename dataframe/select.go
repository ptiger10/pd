package dataframe

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

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