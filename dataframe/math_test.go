package dataframe

import (
	"math"
	"testing"

	"github.com/ptiger10/pd/series"
)

func Test_Math(t *testing.T) {
	df := MustNew([]interface{}{[]float64{1, 3, 5}, []float64{-3, math.NaN(), -1}, []float64{-5, 0, 5}},
		Config{Col: []string{"foo", "bar", "baz"}})
	tests := []struct {
		name  string
		input *DataFrame
		fn    func(*DataFrame) *series.Series
		want  *series.Series
	}{
		{name: "Empty", input: newEmptyDataFrame(), fn: (*DataFrame).Sum, want: series.MustNew(nil)},
		{"Sum", df, (*DataFrame).Sum,
			series.MustNew([]float64{9, -4, 0}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "sum"}),
		},
		{"Mean", df, (*DataFrame).Mean,
			series.MustNew([]float64{3, -2, 0}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "mean"}),
		},
		{"Min", df, (*DataFrame).Min,
			series.MustNew([]float64{1, -3, -5}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "min"}),
		},
		{"Max", df, (*DataFrame).Max,
			series.MustNew([]float64{5, -1, 5}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "max"}),
		},
		{"Std", df, (*DataFrame).Std,
			series.MustNew([]float64{1.632993161855452, 1, 4.08248290463863}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "std"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.input)
			if !series.Equal(got, tt.want) {
				t.Errorf("%v() got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
