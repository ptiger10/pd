package dataframe

import (
	"math"
	"testing"

	"github.com/ptiger10/pd/series"
)

func Test_Math(t *testing.T) {
	type want struct {
		data   interface{}
		config series.Config
	}
	tests := []struct {
		name string
		fn   func(*DataFrame) *series.Series
		want want
	}{
		{"Sum", (*DataFrame).Sum,
			want{[]float64{9, -4, 0}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "sum"}},
		},
		{"Mean", (*DataFrame).Mean,
			want{[]float64{3, -2, 0}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "mean"}},
		},
		{"Min", (*DataFrame).Min,
			want{[]float64{1, -3, -5}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "min"}},
		},
		{"Max", (*DataFrame).Max,
			want{[]float64{5, -1, 5}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "max"}},
		},
		{"Std", (*DataFrame).Std,
			want{[]float64{1.632993161855452, 1, 4.08248290463863}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "std"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df, err := New(
				[]interface{}{[]float64{1, 3, 5}, []float64{-3, math.NaN(), -1}, []float64{-5, 0, 5}},
				Config{Cols: []interface{}{"foo", "bar", "baz"}})
			if err != nil {
				t.Errorf("%v() error: %v", tt.name, err)
			}
			got := tt.fn(df)
			want, _ := series.New(tt.want.data, tt.want.config)
			if !series.Equal(got, want) {
				t.Errorf("%v() got %v, want %v", tt.name, got, want)
			}
		})
	}
}
