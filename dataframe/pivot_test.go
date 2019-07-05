package dataframe

import (
	"fmt"
	"testing"

	"github.com/d4l3k/messagediff"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// TODO: Support different types
func TestDataFrame_Pivot_stackColumn(t *testing.T) {
	type args struct {
		col int
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  *DataFrame
	}{
		{name: "one group",
			input: MustNew([]interface{}{[]int64{1, 1, 1}, []string{"foo", "bar", "baz"}},
				Config{Col: []string{"A", "B"}}),
			args: args{col: 0},
			want: MustNew([]interface{}{[]string{"foo", "bar", "baz"}},
				Config{Col: []string{"1"}, ColName: "A"})},
		{name: "two groups",
			input: MustNew([]interface{}{[]int64{1, 1, 2}, []string{"foo", "bar", "baz"}},
				Config{Col: []string{"A", "B"}}),
			args: args{col: 0},
			want: MustNew([]interface{}{[]string{"foo", "bar", ""}, []string{"", "", "baz"}},
				Config{Col: []string{"1", "2"}, ColName: "A"})},
		// {name: "different types",
		// 	input: MustNew([]interface{}{[]int64{1, 1, 2}, []string{"foo", "bar", "baz"}},
		// 		Config{Col: []string{"A", "B"}}),
		// 	args: args{col: 0},
		// 	want: MustNew([]interface{}{[]string{"foo", "bar", ""}, []string{"", "", "baz"}},
		// 		Config{Col: []string{"1", "2"}, ColName: "A"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.stackCol(tt.args.col)
			if !Equal(got, tt.want) {
				t.Errorf("DataFrame.stackColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: Multiple index levels
func TestDataFrame_Pivot_stackIndex(t *testing.T) {
	type args struct {
		col int
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  *DataFrame
	}{
		{name: "one group, one index level",
			input: MustNew([]interface{}{[]string{"foo", "bar", "baz"}},
				Config{Index: []int64{1, 1, 1}, IndexName: "A"}),
			args: args{col: 0},
			want: MustNew([]interface{}{[]string{"foo", "bar", "baz"}},
				Config{Col: []string{"1"}, ColName: "A"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.stackIndex(tt.args.col)
			if !Equal(got, tt.want) {
				t.Errorf("DataFrame.stackIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranspose(t *testing.T) {
	tests := []struct {
		name  string
		input *DataFrame
		want  *DataFrame
	}{
		{name: "pass",
			input: MustNew([]interface{}{[]string{"qux", "quux", "quuz"}, []string{"foo", "bar", "baz"}},
				Config{Col: []string{"A", "B"}, Index: []string{"1", "2", "3"}}),
			want: MustNew([]interface{}{[]string{"qux", "foo"}, []string{"quux", "bar"}, []string{"quuz", "baz"}},
				Config{Col: []string{"1", "2", "3"}, Index: []string{"A", "B"}})},
		// TODO: support transposing index name onto column (impossible to add to transposeSeries)
		{name: "pass with column name",
			input: MustNew([]interface{}{[]string{"qux"}, []string{"foo"}},
				Config{Col: []string{"A", "B"}, Index: []string{"1"}, IndexName: "corge"}),
			want: MustNew([]interface{}{[]string{"qux", "foo"}},
				Config{Col: []string{"1"}, ColName: "corge", Index: []string{"A", "B"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Transpose()
			if !Equal(got, tt.want) {
				t.Errorf("DataFrame.Transpose() = %v, want %v", got, tt.want)
				diff, _ := messagediff.PrettyDiff(got, tt.want)
				fmt.Println(diff)
			}
		})
	}
}

func TestTransposeSeries(t *testing.T) {
	multi := MustNew([]interface{}{0, 1, 2}, Config{MultiIndex: []interface{}{"qux"}, MultiCol: [][]string{{"foo", "bar", "baz"}, {"4", "5", "6"}}})
	multi.cols.Levels[1].DataType = options.Int64

	type args struct {
		s *series.Series
	}
	tests := []struct {
		name string
		args args
		want *DataFrame
	}{
		{name: "single index", args: args{s: series.MustNew([]int{0, 1, 4},
			series.Config{Index: []string{"foo", "bar", "baz"}, Name: "qux"})},
			want: MustNew([]interface{}{0, 1, 4}, Config{MultiIndex: []interface{}{"qux"}, Col: []string{"foo", "bar", "baz"}})},
		{name: "multi index", args: args{s: series.MustNew([]int{0, 1, 2},
			series.Config{MultiIndex: []interface{}{[]string{"foo", "bar", "baz"}, []int{4, 5, 6}}, Name: "qux"})},
			want: multi},
		{name: "int in name to int index", args: args{s: series.MustNew([]int{0, 1, 2},
			series.Config{Index: []string{"foo", "bar", "baz"}, Name: "0"})},
			want: MustNew([]interface{}{0, 1, 2}, Config{Index: 0, Col: []string{"foo", "bar", "baz"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transposeSeries(tt.args.s)
			if !Equal(got, tt.want) {
				t.Errorf("transposeSeries() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestDataFrame_Pivot(t *testing.T) {
	type args struct {
		index   int
		data    int
		col     int
		aggFunc string
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  *DataFrame
	}{
		{name: "pass",
			input: MustNew([]interface{}{[]string{"foo", "foo", "foo"}, []string{"bar", "bar", "baz"}, []int{1, 2, 3}},
				Config{Col: []string{"A", "B", "C"}}),
			args: args{index: 0, data: 2, col: 1, aggFunc: "sum"},
			want: MustNew([]interface{}{3.0, 3.0},
				Config{Col: []string{"bar", "baz"}, Index: "foo"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.pivot(tt.args.index, tt.args.data, tt.args.col, tt.args.aggFunc)
			if !Equal(got, tt.want) {
				t.Errorf("transposeSeries() = %v, want %v", got, tt.want)
			}

		})
	}
}
