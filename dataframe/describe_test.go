package dataframe

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/series"
)

func TestDataFrame_Describe(t *testing.T) {
	type want struct {
		len          int
		numCols      int
		numIdxLevels int
		numColLevels int
		dataType     string
		dataTypes    *series.Series
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
				dataType: "empty", dataTypes: series.MustNew(nil),
			}},
		{"default index, col",
			MustNew([]interface{}{"foo"}),
			want{
				len: 1, numCols: 1, numIdxLevels: 1, numColLevels: 1,
				dataType: "string", dataTypes: series.MustNew("string", series.Config{Name: "datatypes"}),
			}},
		{"multi index, single col",
			MustNew([]interface{}{"foo"}, Config{MultiIndex: []interface{}{"baz", "qux"}}),
			want{
				len: 1, numCols: 1, numIdxLevels: 2, numColLevels: 1,
				dataType: "string", dataTypes: series.MustNew("string", series.Config{Name: "datatypes"}),
			}},
		{"single index, two cols",
			MustNew([]interface{}{"foo", "bar"}, Config{Col: []string{"baz", "qux"}}),
			want{
				len: 1, numCols: 2, numIdxLevels: 1, numColLevels: 1,
				dataType: "string", dataTypes: series.MustNew([]string{"string", "string"}, series.Config{Name: "datatypes"}),
			}},
		{"single index, multi col",
			MustNew([]interface{}{"foo", "bar"}, Config{MultiCol: [][]string{{"baz", "qux"}, {"corge", "fred"}}}),
			want{
				len: 1, numCols: 2, numIdxLevels: 1, numColLevels: 2,
				dataType: "string", dataTypes: series.MustNew([]string{"string", "string"}, series.Config{Name: "datatypes"}),
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
			gotDataType := df.dataTypePrinter()
			if gotDataType != tt.want.dataType {
				t.Errorf("df.dataTypePrinter: got %v, want %v", gotDataType, tt.want.dataType)
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
		{name: "equal", input: df, args: args{df2: df}, want: true},
		{"equal empty", newEmptyDataFrame(), args{newEmptyDataFrame()}, true},
		{"equal empty copy", newEmptyDataFrame().Copy(), args{newEmptyDataFrame()}, true},
		{"not equal", df, args{MustNew([]interface{}{[]string{"foo"}, []string{"bar"}})}, false},
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
		name   string
		config Config
		want   want
	}{
		{name: "empty config", config: Config{},
			want: want{colWidths: []int{3, 4}, exclusionsTable: [][]bool{{false, false}}}},
		{"single level",
			Config{Col: []string{"corge", "bar"}, ColName: "grapply"},
			want{[]int{5, 4}, [][]bool{{false, false}}}},
		{"multi level",
			Config{MultiCol: [][]string{{"corge", "bar"}, {"qux", "quuz"}}, MultiColNames: []string{"grapply", "grault"}},
			want{[]int{5, 4}, [][]bool{{false, false}, {false, false}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
		df := MustNew([]interface{}{[]string{"a", "foo"}, []string{"b", "quux"}}, tt.config)
		excl := df.makeColumnExclusionsTable()
		got := df.maxColWidths(excl)
		if !reflect.DeepEqual(got, tt.want.colWidths) {
			t.Errorf("df.maxColWidths() got %v, want %v", got, tt.want.colWidths)
		}
		if !reflect.DeepEqual(excl, tt.want.exclusionsTable) {
			t.Errorf("df.makeColumnExclusionsTable() got %v, want %v", excl, tt.want.exclusionsTable)
		}
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
		t.Errorf("df.maxColWidths() got %v, want %v for df \n%v", got, want, df)
	}
}
