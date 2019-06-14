package dataframe

import (
	"math"

	"github.com/ptiger10/pd/series"
)

// Sum all numerical or boolean columns.
func (df *DataFrame) Sum() *series.Series {
	ret, _ := series.New(nil)
	for _, s := range df.s {
		if calc := s.Sum(); !math.IsNaN(calc) {
			newS := series.MustNew(calc, series.Idx(s.Name))
			ret.InPlace.Join(newS)
		}
	}
	return ret
}

// Mean all numerical or boolean columns.
func (df *DataFrame) Mean() *series.Series {
	ret, _ := series.New(nil)
	for _, s := range df.s {
		if calc := s.Mean(); !math.IsNaN(calc) {
			newS := series.MustNew(calc, series.Idx(s.Name))
			ret.InPlace.Join(newS)
		}
	}
	return ret
}
