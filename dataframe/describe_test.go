package dataframe

import (
	"fmt"
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
		{name: "default index, col",
			input: MustNew([]interface{}{"foo"}),
			want: want{
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
		{"empty",
			newEmptyDataFrame(),
			want{
				len: 0, numCols: 0, numIdxLevels: 0, numColLevels: 0,
				dataType: "None", dataTypes: series.MustNew(nil),
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
				t.Errorf("df.dataType(): got %v, want %v", gotDataType, tt.want.dataType)
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
	tests := []struct {
		config Config
		want   []int
	}{
		{Config{}, []int{3, 4}},
		{Config{Col: []string{"corge", "bar"}, ColName: "grapply"}, []int{5, 4}},
		{Config{MultiCol: [][]string{{"corge", "bar"}, {"qux", "quuz"}}, MultiColNames: []string{"grapply", "grault"}}, []int{5, 4}},
	}

	for _, tt := range tests {
		df := MustNew([]interface{}{[]string{"a", "foo"}, []string{"b", "quux"}}, tt.config)
		excl := df.makeExclusionsTable()
		got := df.maxColWidths(excl)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("df.maxColWidths() got %v, want %v for df \n%v", got, tt.want, df)
		}
	}
}

func TestMaxColWidthExclusions(t *testing.T) {
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

func TestMakeExclusionTable(t *testing.T) {
	df := MustNew([]interface{}{"foo", "bar"}, Config{MultiCol: [][]string{{"baz", "qux"}, {"quux", "quuz"}}})
	got := df.makeExclusionsTable()
	want := [][]bool{{false, false}, {false, false}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("df.MakeExclusionsTable() got %v, want %v", got, want)
	}
}

func TestNames(t *testing.T) {
	df := newEmptyDataFrame()
	fmt.Println(df.NumCols())
	got := df.cols.Names()
	want := []string{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("df.cols.Names() for nil got %v, want %v", got, want)
	}
}
