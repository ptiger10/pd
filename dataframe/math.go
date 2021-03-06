package dataframe

import (
	"math"

	"github.com/ptiger10/pd/series"
)

func (df *DataFrame) math(name string, fn func(s *series.Series) float64) *series.Series {
	if Equal(df, newEmptyDataFrame()) {
		return series.MustNew(nil)
	}
	var vals []interface{}
	var idx []interface{}
	for m := 0; m < df.NumCols(); m++ {
		s := df.hydrateSeries(m)
		if calc := fn(s); !math.IsNaN(calc) {
			vals = append(vals, calc)
			idx = append(idx, s.Name())
		}
	}
	ret := series.MustNew(nil)
	for i := 0; i < len(vals); i++ {
		// ducks safe method because values are assumed to be supported
		s := series.MustNew(vals[i], series.Config{Index: idx[i], Name: name})
		ret.InPlace.Join(s)
	}

	return ret
}

// Sum all numerical or boolean columns.
func (df *DataFrame) Sum() *series.Series {
	return df.math("sum", (*series.Series).Sum)
}

// Mean of all numerical or boolean columns.
func (df *DataFrame) Mean() *series.Series {
	return df.math("mean", (*series.Series).Mean)
}

// Median of all numerical or boolean columns.
func (df *DataFrame) Median() *series.Series {
	return df.math("median", (*series.Series).Median)
}

// Min all numerical columns.
func (df *DataFrame) Min() *series.Series {
	return df.math("min", (*series.Series).Min)
}

// Max all numerical columns.
func (df *DataFrame) Max() *series.Series {
	return df.math("max", (*series.Series).Max)
}

// Std returns the standard deviation of all numerical columns.
func (df *DataFrame) Std() *series.Series {
	return df.math("std", (*series.Series).Std)
}
