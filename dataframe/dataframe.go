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

func (df *DataFrame) rowsIn(rowPositions []int) (*DataFrame, error) {
	if err := df.ensureAlignment(); err != nil {
		return df, fmt.Errorf("dataframe internal alignment error: %v", err)
	}
	var seriesSlice []*series.Series
	for i := 0; i < df.Cols(); i++ {
		s, err := df.s[i].SelectRows(rowPositions).Series()
		if err != nil {
			return nil, fmt.Errorf("dataframe.rowsIn() selecting rows within series (position %v): %v", i, err)
		}
		seriesSlice = append(seriesSlice, s)
	}
	idx, err := df.index.In(rowPositions)
	if err != nil {
		return nil, fmt.Errorf("dataframe.rowsIn() selecting index labels: %v", err)
	}
	df = newFromComponents(seriesSlice, idx, df.cols, df.name)
	return df, nil
}

func (df *DataFrame) colsIn(colPositions []int) (*DataFrame, error) {
	if err := df.ensureAlignment(); err != nil {
		return df, fmt.Errorf("dataframe internal alignment error: %v", err)
	}
	var seriesSlice []*series.Series
	for _, pos := range colPositions {
		if pos > df.Cols() {
			return nil, fmt.Errorf("dataframe.colsIn(): invalid col position %d (max: %d)", pos, df.Cols()-1)
		}
		seriesSlice = append(seriesSlice, df.s[pos])
	}
	columnsSlice, err := df.cols.In(colPositions)
	if err != nil {
		return nil, fmt.Errorf("dataframe.colsIn(): %v", err)
	}

	df = newFromComponents(seriesSlice, df.index, columnsSlice, df.name)
	return df, nil
}

func newFromComponents(s []*series.Series, idx index.Index, cols index.Columns, name string) *DataFrame {
	if s == nil {
		df, _ := New(nil)
		return df
	}
	return &DataFrame{
		s:     s,
		index: idx,
		cols:  cols,
		name:  name,
	}
}

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
	s, err := newSingleIndexSeries(vals, idx, "datatypes")
	if err != nil {
		log.Printf("DataTypes(): %v", err)
		return nil
	}
	return s
}

// newSingleIndexSeries constructs a Series with a single-level index from values and index slices. Used to convert DataFrames to Series.
func newSingleIndexSeries(values []interface{}, idx []interface{}, name string) (*series.Series, error) {
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
