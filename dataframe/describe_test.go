package dataframe

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

func TestDataFrame_Describe(t *testing.T) {
	type want struct {
		len             int
		numCols         int
		numIdxLevels    int
		numColLevels    int
		dataType        options.DataType
		dataTypePrinter string
		dataTypes       *series.Series
	}
	tests := []struct {
		name  string
		input *DataFrame
		want  want
	}{
		{name: "empty",
			input: newEmptyDataFrame(),
			want: want{
				len: 0, numCols: 0, numIdxLevels: 0, numColLevels: 0,
				dataType: options.None, dataTypePrinter: "empty", dataTypes: series.MustNew(nil),
			}},
		{"default index, col",
			MustNew([]interface{}{"foo"}),
			want{
				len: 1, numCols: 1, numIdxLevels: 1, numColLevels: 1,
				dataType: options.String, dataTypePrinter: "string", dataTypes: series.MustNew("string", series.Config{Name: "datatypes"}),
			}},
		{"multi index, single col",
			MustNew([]interface{}{"foo"}, Config{MultiIndex: []interface{}{"baz", "qux"}}),
			want{
				len: 1, numCols: 1, numIdxLevels: 2, numColLevels: 1,
				dataType: options.String, dataTypePrinter: "string", dataTypes: series.MustNew("string", series.Config{Name: "datatypes"}),
			}},
		{"single index, two cols, mixed types",
			MustNew([]interface{}{"foo", 5}, Config{Col: []string{"baz", "qux"}}),
			want{
				len: 1, numCols: 2, numIdxLevels: 1, numColLevels: 1,
				dataType: options.Unsupported, dataTypePrinter: "mixed", dataTypes: series.MustNew([]string{"string", "int64"}, series.Config{Name: "datatypes"}),
			}},
		{"single index, multi col",
			MustNew([]interface{}{"foo", "bar"}, Config{MultiCol: [][]string{{"baz", "qux"}, {"corge", "fred"}}}),
			want{
				len: 1, numCols: 2, numIdxLevels: 1, numColLevels: 2,
				dataType: options.String, dataTypePrinter: "string", dataTypes: series.MustNew([]string{"string", "string"}, series.Config{Name: "datatypes"}),
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			gotLen := df.Len()
			if gotLen != tt.want.len {
				t.Errorf("df.Len(): got %v, want %v", gotLen, tt.want.len)
			}
			gotNumCols := df.NumCols()
			if gotNumCols != tt.want.numCols {
				t.Errorf("df.NumCols(): got %v, want %v", gotNumCols, tt.want.numCols)
			}
			gotNumIdxLevels := df.IndexLevels()
			if gotNumIdxLevels != tt.want.numIdxLevels {
				t.Errorf("df.IndexLevels(): got %v, want %v", gotNumIdxLevels, tt.want.numIdxLevels)
			}
			gotNumColLevels := df.ColLevels()
			if gotNumColLevels != tt.want.numColLevels {
				t.Errorf("df.ColLevels(): got %v, want %v", gotNumColLevels, tt.want.numColLevels)
			}
			gotDataType := df.dataType()
			if gotDataType != tt.want.dataType {
				t.Errorf("df.gotDataType: got %v, want %v", gotDataType, tt.want.dataType)
			}
			gotDataTypePrinter := df.dataTypePrinter()
			if gotDataTypePrinter != tt.want.dataTypePrinter {
				t.Errorf("df.dataTypePrinter: got %v, want %v", gotDataTypePrinter, tt.want.dataTypePrinter)
			}
			gotDataTypes := df.DataTypes()
			if !series.Equal(gotDataTypes, tt.want.dataTypes) {
				t.Errorf("df.DataTypes(): got %v, want %v", gotDataTypes, tt.want.dataTypes)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo"}, []string{"bar"}}, Config{Index: "corge", Col: []string{"baz", "qux"}})
	type args struct {
		df2 *DataFrame
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  bool
	}{
		{name: "equal", input: df,
			args: args{df2: MustNew([]interface{}{[]string{"foo"}, []string{"bar"}}, Config{Index: "corge", Col: []string{"baz", "qux"}})},
			want: true},
		{"equal empty", newEmptyDataFrame(), args{newEmptyDataFrame()}, true},
		{"equal empty copy", newEmptyDataFrame().Copy(), args{newEmptyDataFrame()}, true},
		{"not equal: values", df,
			args{MustNew([]interface{}{[]string{"foo"}, []string{"bar"}})}, false},
		{"not equal: cols", df,
			args{MustNew([]interface{}{[]string{"foo"}, []string{"bar"}}, Config{Index: "corge", Col: []string{"fred", "qux"}})}, false},
		{"not equal: name", df,
			args{MustNew([]interface{}{[]string{"foo"}, []string{"bar"}}, Config{Index: "corge", Col: []string{"baz", "qux"}, Name: "quux"})}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Equal(tt.input, tt.args.df2)
			if got != tt.want {
				t.Errorf("Equal() got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxColWidth(t *testing.T) {
	type want struct {
		colWidths       []int
		exclusionsTable [][]bool
	}
	tests := []struct {
		name  string
		input *DataFrame
		want  want
	}{
		{name: "empty config", input: MustNew([]interface{}{[]string{"a", "foo"}, []string{"b", "quux"}}),
			want: want{colWidths: []int{3, 4}, exclusionsTable: [][]bool{{false, false}}}},
		{"single level",
			MustNew([]interface{}{[]string{"a", "foo"}, []string{"b", "quux"}},
				Config{Col: []string{"corge", "bar"}, ColName: "grapply"}),
			want{[]int{5, 4}, [][]bool{{false, false}}}},
		{"multi level",
			MustNew([]interface{}{[]string{"a", "foo"}, []string{"b", "quux"}},
				Config{MultiCol: [][]string{{"corge", "bar"}, {"qux", "quuz"}}, MultiColNames: []string{"grapply", "grault"}}),
			want{[]int{5, 4}, [][]bool{{false, false}, {false, false}}}},
		{"nil: empty colWidths", newEmptyDataFrame(), want{nil, [][]bool{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input
			excl := df.makeColumnExclusionsTable()
			got := df.maxColWidths(excl)
			if !reflect.DeepEqual(excl, tt.want.exclusionsTable) {
				t.Errorf("df.makeColumnExclusionsTable() got %v, want %v", excl, tt.want.exclusionsTable)
			}
			if !reflect.DeepEqual(got, tt.want.colWidths) {
				t.Errorf("df.maxColWidths() got %v, want %v", got, tt.want.colWidths)
			}
		})
	}
}

func TestMaxColWidthExcludeRepeat(t *testing.T) {
	df := MustNew(
		[]interface{}{[]string{"a", "b"}, []string{"c", "quux"}},
		Config{MultiCol: [][]string{{"waldo", "waldo"}, {"d", "e"}}})
	excl := [][]bool{{false, true}, {false, false}}
	got := df.maxColWidths(excl)
	want := []int{5, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("df.maxColWidths() got %v, want %v", got, want)
	}
}

func TestHeadTail(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo", "bar", "baz", "qux"}}, Config{Index: []int{0, 1, 2, 3}})
	type args struct {
		n int
	}
	tests := []struct {
		name  string
		input *DataFrame
		fn    func(*DataFrame, int) *DataFrame
		args  args
		want  *DataFrame
	}{
		{name: "head", input: df, fn: (*DataFrame).Head, args: args{n: 2},
			want: MustNew([]interface{}{[]string{"foo", "bar"}}, Config{Index: []int{0, 1}})},
		{name: "head - max", input: df, fn: (*DataFrame).Head, args: args{n: 10},
			want: df},
		{name: "tail", input: df, fn: (*DataFrame).Tail, args: args{n: 2},
			want: MustNew([]interface{}{[]string{"baz", "qux"}}, Config{Index: []int{2, 3}})},
		{name: "tail - max", input: df, fn: (*DataFrame).Tail, args: args{n: 10},
			want: df},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.input, tt.args.n)
			if !Equal(got, tt.want) {
				t.Errorf("df.Head/Tail() got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataFrame_Export(t *testing.T) {
	tests := []struct {
		name  string
		input *DataFrame
		want  [][]interface{}
	}{
		{name: "pass", input: MustNew([]interface{}{"foo"}, Config{Index: "bar", Col: []string{"baz"}}),
			want: [][]interface{}{{nil, "baz"}, {"bar", "foo"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Export()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("df.Export() got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataFrame_ExportToCSV(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  bool
	}{
		{name: "pass", input: MustNew([]interface{}{"foo"}, Config{Index: "bar", Col: []string{"baz"}}),
			want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.ExportToCSV("output_test.csv")
			//TODO: move ReadCSV to dataframe package to rehydrate output and compare to input
		})
	}
}
