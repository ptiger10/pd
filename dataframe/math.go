package dataframe

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/ptiger10/pd/series"
)

func (df *DataFrame) math(name string, fn func(s *series.Series) float64) *series.Series {
	if df == nil {
		s, _ := series.New(nil)
		return s
	}
	var vals []interface{}
	var idx []interface{}
	for _, s := range df.s {
		if calc := fn(s); !math.IsNaN(calc) {
			vals = append(vals, calc)
			idx = append(idx, s.Name())
		}
	}
	s, err := newSingleIndexSeries(vals, idx, name)
	if err != nil {
		log.Printf("%v: %v", fmt.Sprintf("%v()", strings.Title(name)), err)
		return nil
	}
	return s
}

// Sum all numerical or boolean columns.
func (df *DataFrame) Sum() *series.Series {
	return df.math("sum", (*series.Series).Sum)
}

// Mean all numerical or boolean columns.
func (df *DataFrame) Mean() *series.Series {
	return df.math("mean", (*series.Series).Mean)
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
