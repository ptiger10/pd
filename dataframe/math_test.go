package dataframe

import (
	"testing"

	"github.com/ptiger10/pd/series"
)

func Test_Math(t *testing.T) {
	type args struct {
		data   []interface{}
		config Config
	}
	type want struct {
		data   interface{}
		config series.Config
	}
	tests := []struct {
		name string
		args args
		fn   func(*DataFrame) *series.Series
		want want
	}{
		{
			"Sum",
			args{[]interface{}{[]float64{1, 2}, []float64{3, 4}}, Config{Cols: []interface{}{"foo", "bar"}}},
			(*DataFrame).Sum,
			want{[]float64{3, 7}, series.Config{Index: []string{"foo", "bar"}, Name: "sum"}},
		},
		{
			"Mean",
			args{[]interface{}{[]float64{1, 3, 5}, []float64{-3, -2, -1}, []float64{-5, 0, 5}}, Config{Cols: []interface{}{"foo", "bar", "baz"}}},
			(*DataFrame).Mean,
			want{[]float64{3, -2, 0}, series.Config{Index: []string{"foo", "bar", "baz"}, Name: "mean"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df, err := New(tt.args.data, tt.args.config)
			if err != nil {
				t.Errorf("%v(): %v", tt.name, err)
			}
			got := tt.fn(df)
			want, _ := series.New(tt.want.data, tt.want.config)
			if !series.Equal(got, want) {
				t.Errorf("%v() = %v, want %v", tt.name, got, want)
			}
		})
	}
}

func Test_Sum(t *testing.T) {
	c := Config{Cols: []interface{}{"fooCol", "barCol"}, Index: []string{"foo", "bar", "baz"}}
	df, err := New([]interface{}{[]float64{1, 2, 3}, []float64{4, 5, 6}}, c)
	if err != nil {
		t.Error(err)
	}
	got := df.Sum()
	want, _ := series.New([]float64{6, 15}, series.Config{Index: []string{"fooCol", "barCol"}, Name: "sum"})
	if !series.Equal(got, want) {
		t.Errorf("df.Sum() returned %v, want %v", got, want)
	}
}

func Test_Mean(t *testing.T) {
	c := Config{Cols: []interface{}{"fooCol", "barCol"}, Index: []string{"foo", "bar", "baz"}}
	df, err := New([]interface{}{[]float64{1, 2, 3}, []float64{4, 5, 6}}, c)
	if err != nil {
		t.Error(err)
	}
	got := df.Mean()
	want, _ := series.New([]float64{2, 5}, series.Config{Index: []string{"fooCol", "barCol"}, Name: "sum"})
	if !series.Equal(got, want) {
		t.Errorf("df.Mean() returned %v, want %v", got, want)
	}
}
