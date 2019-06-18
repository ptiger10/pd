package dataframe

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/series"
)

// A DataFrame is a 2D collection of one or more Series with a shared index and associated columns.
type DataFrame struct {
	name  string
	s     []*series.Series
	cols  index.Columns
	index index.Index
}

// Len returns the number of values in each Series of the DataFrame.
func (df DataFrame) Len() int {
	if df.s == nil {
		return 0
	}
	return df.s[0].Len()
}

// Name returns the DataFrame's name.
func (df DataFrame) Name() string {
	return df.name
}

// Rename the DataFrame.
func (df *DataFrame) Rename(name string) {
	df.name = name
}

// Cols returns the number of columsn in the DataFrame.
func (df *DataFrame) Cols() int {
	if df.s == nil {
		return 0
	}
	return len(df.s)
}

// IndexLevels returns the number of index levels in the DataFrame.
func (df *DataFrame) IndexLevels() int {
	return df.index.Len()
}

// Copy creates a new deep copy of a Series.
func (df *DataFrame) Copy() *DataFrame {
	var sCopy []*series.Series
	for i := 0; i < len(df.s); i++ {
		sCopy = append(sCopy, df.s[i].Copy())
	}
	idxCopy := df.index.Copy()
	colsCopy := df.cols.Copy()
	dfCopy := &DataFrame{
		s:     sCopy,
		index: idxCopy,
		cols:  colsCopy,
		name:  df.name,
	}
	// dfCopy.Apply = Apply{s: copyS}
	// dfCopy.Filter = Filter{s: copyS}
	// dfCopy.Index = Index{s: copyS}
	// dfCopy.InPlace = InPlace{s: copyS}
	return dfCopy
}

// // in copies a Series then subsets it to include only index items and values at the integer positions supplied
// func (df *DataFrame) in(rowPositions []int, colPositions []int) (*DataFrame, error) {
// 	if err := df.ensureAlignment(); err != nil {
// 		return df, fmt.Errorf("dataframe internal alignment error: %v", err)
// 	}
// 	if rowPositions == nil && colPositions == nil {
// 		return nil, nil
// 	}

// 	df = df.Copy()
// 	values, err := df.values.In(positions)
// 	if err != nil {
// 		return nil, fmt.Errorf("Series.in() values: %v", err)
// 	}
// 	df.values = values
// 	for i, level := range df.index.Levels {
// 		df.index.Levels[i].Labels, err = level.Labels.In(positions)
// 		if err != nil {
// 			return nil, fmt.Errorf("Series.in() index: %v", err)
// 		}
// 	}
// 	df.index.Refresh()
// 	return df, nil
// }

// func (df *DataFrame) Col(label string) *Series {

// }

// DataTypes returns the DataTypes of the Series in the DataFrame.
func (df *DataFrame) DataTypes() *series.Series {
	var vals []interface{}
	var idx []interface{}
	for _, s := range df.s {
		vals = append(vals, s.DataType())
		idx = append(idx, s.Name())
	}
	s, err := newSingleIndexFromSeries(vals, idx, "datatypes")
	if err != nil {
		log.Printf("DataTypes(): %v", err)
		return nil
	}
	return s
}

func newSingleIndexFromSeries(values []interface{}, idx []interface{}, name string) (*series.Series, error) {
	ret, err := series.New(nil)
	if err != nil {
		return nil, fmt.Errorf("internal error: newFromSeries(): %v", err)
	}
	if len(values) != len(idx) {
		return nil, fmt.Errorf("internal error: newFromSeries(): values must have same length as index: %d != %d", len(values), len(idx))
	}
	for i := 0; i < len(values); i++ {
		n, err := series.New(values[i], series.Config{Index: idx[i], Name: name})
		if err != nil {
			return nil, fmt.Errorf("internal error: newFromSeries(): %v", err)
		}
		ret.InPlace.Join(n)
	}
	return ret, nil
}

// dataType is the data type of the DataFrame's values. Mimics reflect.Type with the addition of time.Time as DateTime.
func (df *DataFrame) dataType() string {
	uniqueTypes := df.DataTypes().UniqueVals()
	if len(uniqueTypes) == 1 {
		return df.s[0].DataType()
	}
	return "mixed"
}
