package dataframe

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/ptiger10/pd/series"
)

func (df *DataFrame) math(name string, fn func(s *series.Series) float64) *series.Series {
	var vals []interface{}
	var idx []interface{}
	for _, s := range df.s {
		if calc := fn(s); !math.IsNaN(calc) {
			vals = append(vals, calc)
			idx = append(idx, s.Name())
		}
	}
	s, err := newSingleIndexFromSeries(vals, idx, name)
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
