package dataframe

import (
	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/series"
)

// A DataFrame is a 2D collection of one or more Series with a shared index and associated columns.
type DataFrame struct {
	Name  string
	s     []*series.Series
	cols  index.Columns
	index index.Index
}

// Len returns the number of values in each Series of the DataFrame.
func (df *DataFrame) Len() int {
	if df.s == nil {
		return 0
	}
	return df.s[0].Len()
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
		Name:  df.Name,
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

// DT returns the DataTypes of the Series in the DataFrame.
func (df *DataFrame) DT() *series.Series {
	ret, _ := series.New(nil)
	for _, s := range df.s {
		dt := series.MustNew(s.DataType(), series.Config{Index: s.Name, Name: "datatypes"})
		ret.InPlace.Join(dt)
	}
	return ret
}

// dataType is the data type of the DataFrame's values. Mimics reflect.Type with the addition of time.Time as DateTime.
func (df *DataFrame) dataType() string {
	uniqueTypes := df.DT().UniqueVals()
	if len(uniqueTypes) == 1 {
		return df.s[0].DataType()
	}
	return "mixed"
}
