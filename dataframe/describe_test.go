package dataframe

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ptiger10/pd/series"
)

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

func TestDataTypes(t *testing.T) {
	df, err := New([]interface{}{"foo"}, Config{Index: "bar"})
	if err != nil {
		t.Errorf("df.DataTypes(): %v", err)
	}
	got := df.DataTypes()
	want, err := series.New("string", series.Config{Name: "datatypes"})
	if err != nil {
		fmt.Println(err)
	}
	if !series.Equal(got, want) {
		t.Errorf("df.DataTypes() returned %#v, want %#v", got, want)
	}
}

func TestDataType(t *testing.T) {
	df, _ := New([]interface{}{"foo"})
	got := df.dataType()
	want := "string"
	if got != want {
		t.Errorf("df.dataType() returned %v, want %v", got, want)
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
