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

// DT returns the DataTypes of the Series in the DataFrame.
func (df *DataFrame) DT() *series.Series {
	ret, _ := series.New(nil)
	for _, s := range df.s {
		dt := series.MustNew(s.DataType(), series.Idx(s.Name))
		ret.InPlace.Join(dt)
	}
	ret.Name = "datatypes"
	return ret
}

// DataType is the data type of the DataFrame's values. Mimics reflect.Type with the addition of time.Time as DateTime.
func (df *DataFrame) DataType() string {
	uniqueTypes := df.DT().UniqueVals()
	if len(uniqueTypes) == 1 {
		return df.s[0].DataType()
	}
	return "mixed"
}

// Len returns the number of values in each Series of the DataFrame.
func (df *DataFrame) Len() int {
	if df.s != nil {
		return df.s[0].Len()
	}
	return 0
}

// Cols returns the number of columsn in the DataFrame.
func (df *DataFrame) Cols() int {
	if df.s != nil {
		return len(df.s)
	}
	return 0
}

// Levels returns the number of index levels in the DataFrame.
func (df *DataFrame) Levels() int {
	return df.index.Len()
}
