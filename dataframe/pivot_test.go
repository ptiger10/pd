package dataframe

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/d4l3k/messagediff"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

func TestDataFrame_Pivot_stack(t *testing.T) {
	type args struct {
		level int
	}
	type want struct {
		newIdxPositions []int
		vals            [][]interface{}
		newColLevel     []string
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  want
	}{
		{name: "one index, one group, one column", args: args{level: 1},
			input: MustNew([]interface{}{[]float64{1}},
				Config{Col: []string{"A"},
					MultiIndex: []interface{}{
						[]string{"foo"}, []string{"qux"},
					}}),
			want: want{newIdxPositions: []int{0}, vals: [][]interface{}{{1.0}}, newColLevel: []string{"qux"}}},
		{name: "one index, two groups, one column", args: args{level: 1},
			input: MustNew([]interface{}{[]float64{1, 2}},
				Config{Col: []string{"A"},
					MultiIndex: []interface{}{
						[]string{"foo", "foo"},
						[]string{"qux", "bar"},
					}}),
			want: want{newIdxPositions: []int{0}, vals: [][]interface{}{{1.0, 2.0}}, newColLevel: []string{"qux", "bar"}}},
		{name: "one index, one group, two columns", args: args{level: 1},
			input: MustNew([]interface{}{[]float64{1}, []float64{2}},
				Config{Col: []string{"A", "B"},
					MultiIndex: []interface{}{
						[]string{"foo"},
						[]string{"qux"},
					}}),
			want: want{newIdxPositions: []int{0}, vals: [][]interface{}{{1.0, 2.0}},
				newColLevel: []string{"qux", "qux"}},
		},
		{name: "one index, two groups, two columns", args: args{level: 1},
			input: MustNew([]interface{}{[]float64{1, 2}, []float64{3, 4}},
				Config{Col: []string{"A", "B"},
					MultiIndex: []interface{}{
						[]string{"foo", "foo"},
						[]string{"qux", "bar"},
					}}),
			want: want{newIdxPositions: []int{0}, vals: [][]interface{}{{1.0, 3.0, 2.0, 4.0}},
				newColLevel: []string{"qux", "qux", "bar", "bar"}},
		},
		{name: "two indexes, one group, one column", args: args{level: 1},
			input: MustNew([]interface{}{[]float64{1, 2}},
				Config{Col: []string{"A"},
					MultiIndex: []interface{}{
						[]string{"foo", "bar"},
						[]string{"qux", "qux"},
					}}),
			want: want{newIdxPositions: []int{0, 1}, vals: [][]interface{}{{1.0}, {2.0}},
				newColLevel: []string{"qux"}},
		},
		{name: "two indexes, two groups, one column", args: args{level: 1},
			input: MustNew([]interface{}{[]float64{1, 2}},
				Config{Col: []string{"A"},
					MultiIndex: []interface{}{
						[]string{"foo", "bar"},
						[]string{"qux", "baz"},
					}}),
			want: want{newIdxPositions: []int{0, 1}, vals: [][]interface{}{{1.0, nil}, {nil, 2.0}},
				newColLevel: []string{"qux", "baz"}},
		},
		{name: "two indexes, two groups, two columns", args: args{level: 1},
			input: MustNew([]interface{}{[]float64{1, 2}, []float64{3, 4}},
				Config{Col: []string{"A", "B"},
					MultiIndex: []interface{}{
						[]string{"foo", "bar"},
						[]string{"qux", "baz"},
					}}),
			want: want{newIdxPositions: []int{0, 1}, vals: [][]interface{}{{1.0, 3.0, nil, nil}, {nil, nil, 2.0, 4.0}},
				newColLevel: []string{"qux", "qux", "baz", "baz"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newIdxPositions, vals, newColLevel := tt.input.stack(tt.args.level)
			if !reflect.DeepEqual(newIdxPositions, tt.want.newIdxPositions) {
				t.Errorf("DataFrame.stackIndex() newIdxPositions = %v, want %v", newIdxPositions, tt.want.newIdxPositions)
			}
			if !reflect.DeepEqual(vals, tt.want.vals) {
				t.Errorf("DataFrame.stackIndex() valsMatrix = %v, want %v", vals, tt.want.vals)
				diff, _ := messagediff.PrettyDiff(vals, tt.want.vals)
				fmt.Println(diff)
			}
			if !reflect.DeepEqual(newColLevel, tt.want.newColLevel) {
				t.Errorf("DataFrame.stackIndex() newColLevel = %v, want %v", newColLevel, tt.want.newColLevel)
			}

		})
	}
}

// TODO: Support different types
// func TestDataFrame_Pivot_stackColumn(t *testing.T) {
// 	type args struct {
// 		col int
// 	}
// 	tests := []struct {
// 		name  string
// 		input *DataFrame
// 		args  args
// 		want  *DataFrame
// 	}{
// 		{name: "one group",
// 			input: MustNew([]interface{}{[]int64{1, 1, 1}, []string{"foo", "bar", "baz"}},
// 				Config{Col: []string{"A", "B"}}),
// 			args: args{col: 0},
// 			want: MustNew([]interface{}{[]string{"foo", "bar", "baz"}},
// 				Config{Col: []string{"1"}, ColName: "A"})},
// 		{name: "two groups",
// 			input: MustNew([]interface{}{[]int64{1, 1, 2}, []string{"foo", "bar", "baz"}},
// 				Config{Col: []string{"A", "B"}}),
// 			args: args{col: 0},
// 			want: MustNew([]interface{}{[]string{"foo", "bar", ""}, []string{"", "", "baz"}},
// 				Config{Col: []string{"1", "2"}, ColName: "A"})},
// 		// {name: "different types",
// 		// 	input: MustNew([]interface{}{[]int64{1, 1, 2}, []string{"foo", "bar", "baz"}},
// 		// 		Config{Col: []string{"A", "B"}}),
// 		// 	args: args{col: 0},
// 		// 	want: MustNew([]interface{}{[]string{"foo", "bar", ""}, []string{"", "", "baz"}},
// 		// 		Config{Col: []string{"1", "2"}, ColName: "A"})},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := tt.input.stackCol(tt.args.col)
// 			if !Equal(got, tt.want) {
// 				t.Errorf("DataFrame.stackColumn() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
		{name: "one group, one index level, one column",
			input: MustNew([]interface{}{[]string{"foo"}},
				Config{MultiIndex: []interface{}{[]string{"bar"}, []string{"baz"}}, Col: []string{"A"}}),
			args: args{col: 1},
			want: MustNew([]interface{}{[]string{"foo"}},
				Config{Index: []string{"bar"}, MultiCol: [][]string{{"baz"}, {"A"}}})},
		{name: "two groups, one index level, one column",
			input: MustNew([]interface{}{[]string{"foo", "bar"}},
				Config{MultiIndex: []interface{}{[]string{"baz", "baz"}, []string{"qux", "quuz"}}, Col: []string{"A"}}),
			args: args{col: 1},
			want: MustNew([]interface{}{[]string{"foo"}, []string{"bar"}},
				Config{Index: []string{"baz"}, MultiCol: [][]string{{"qux", "quuz"}, {"A", "A"}}})},
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
	multi := MustNew([]interface{}{
		[]string{"foo", "foo", "bar", "bar"}, []string{"baz", "baz", "baz", "qux"},
		[]string{"qux", "quux", "quuz", "quuz"}, []int{1, 2, 3, 4}, []int{5, 6, 7, 8}},
		Config{Col: []string{"A", "B", "C", "D", "E"}})
	df := MustNew([]interface{}{
		[]string{"foo", "foo", "foo"}, []string{"bar", "bar", "baz"}, []int{1, 2, 3}},
		Config{Col: []string{"A", "B", "C"}})

	type args struct {
		index   int
		data    int
		col     int
		aggFunc string
	}
	tests := []struct {
		name    string
		input   *DataFrame
		args    args
		want    *DataFrame
		wantErr bool
	}{
		{name: "NaN: sum",
			input: multi,
			args:  args{index: 0, data: 3, col: 1, aggFunc: "sum"},
			want: MustNew([]interface{}{[]string{"3", "3"}, []string{"NaN", "4"}},
				Config{Col: []string{"baz", "qux"}, Index: []string{"foo", "bar"}}),
			wantErr: false},
		{name: "mean",
			input: df,
			args:  args{index: 0, data: 2, col: 1, aggFunc: "mean"},
			want: MustNew([]interface{}{1.5, 3.0},
				Config{Col: []string{"bar", "baz"}, Index: "foo"}),
			wantErr: false},
		{name: "median",
			input: df,
			args:  args{index: 0, data: 2, col: 1, aggFunc: "median"},
			want: MustNew([]interface{}{1.5, 3.0},
				Config{Col: []string{"bar", "baz"}, Index: "foo"}),
			wantErr: false},
		{name: "min",
			input: df,
			args:  args{index: 0, data: 2, col: 1, aggFunc: "min"},
			want: MustNew([]interface{}{1.0, 3.0},
				Config{Col: []string{"bar", "baz"}, Index: "foo"}),
			wantErr: false},
		{name: "max",
			input: df,
			args:  args{index: 0, data: 2, col: 1, aggFunc: "max"},
			want: MustNew([]interface{}{2.0, 3.0},
				Config{Col: []string{"bar", "baz"}, Index: "foo"}),
			wantErr: false},
		{name: "std",
			input: df,
			args:  args{index: 0, data: 2, col: 1, aggFunc: "std"},
			want: MustNew([]interface{}{0.5, 0.0},
				Config{Col: []string{"bar", "baz"}, Index: "foo"}),
			wantErr: false},
		{name: "fail", input: df, args: args{index: 0, data: 2, col: 1, aggFunc: "unsupported"},
			want: newEmptyDataFrame(), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Pivot(tt.args.index, tt.args.data, tt.args.col, tt.args.aggFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("df.Pivot() error = %v, want %v", err, tt.wantErr)
				return
			}
			// convert to compare NaN
			if strings.Contains(tt.name, "NaN") {
				gotConverted, _ := got.Convert(options.String.String())
				if !Equal(gotConverted, tt.want) {
					t.Errorf("df.Pivot() = %v, want %v", gotConverted, tt.want)
				}
			} else {
				if !Equal(got, tt.want) {
					t.Errorf("df.Pivot() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
