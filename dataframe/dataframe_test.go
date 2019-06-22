package dataframe

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

func TestDataTypes(t *testing.T) {
	df, err := New([]interface{}{"foo"})
	if err != nil {
		t.Errorf("df.DataTypes(): %v", err)
	}
	got := df.DataTypes()
	want, err := series.New("string", series.Config{Index: "0", Name: "datatypes"})
	if err != nil {
		fmt.Println(err)
	}
	if !series.Equal(got, want) {
		t.Errorf("df.DataTypes() returned %v, want %v", got, want)
	}
}

func TestRowsIn(t *testing.T) {
	var err error
	df, err := New(
		[]interface{}{[]string{"foo", "bar", "baz"}},
		Config{Index: []string{"qux", "quux", "corge"}, Cols: []interface{}{"foofoo"}})
	got, err := df.rowsIn([]int{0, 1})
	if err != nil {
		t.Errorf("rowsIn(): %v", err)
	}
	want := MustNew([]interface{}{[]string{"foo", "bar"}},
		Config{Index: []string{"qux", "quux"}, Cols: []interface{}{"foofoo"}})
	if !Equal(got, want) {
		t.Errorf("rowsIn(): got %v, want %v", got, want)
	}
}

func TestColsIn(t *testing.T) {
	df := MustNew(
		[]interface{}{[]string{"foo"}, []string{"bar"}},
		Config{Cols: []interface{}{"baz", "qux"}})
	got, err := df.colsIn([]int{1})
	if err != nil {
		t.Errorf("colsIn(): %v", err)
	}
	want := MustNew([]interface{}{[]string{"bar"}}, Config{Cols: []interface{}{"qux"}})
	if !Equal(got, want) {
		t.Errorf("colsIn(): got %v, want %v", got, want)
	}
}

func TestCol(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo"}, []string{"bar"}},
		Config{Cols: []interface{}{"baz", "qux"}})
	got := df.Col("qux")
	want := series.MustNew([]string{"bar"}, series.Config{Name: "qux"})
	if !series.Equal(got, want) {
		t.Errorf("Col(): got %v, want %v", got, want)
	}
}

func TestEqual(t *testing.T) {
	df := MustNew(
		[]interface{}{[]string{"foo"}, []string{"bar"}},
		Config{Cols: []interface{}{"baz", "qux"}})
	df2 := MustNew(
		[]interface{}{[]string{"foo"}, []string{"bar"}},
		Config{Cols: []interface{}{"baz", "qux"}})
	if !Equal(df, df2) {
		t.Errorf("Equal() did not return true for equivalent df")
	}
}

func TestCols(t *testing.T) {

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
		{Config{Cols: []interface{}{"corge", "bar"}, ColsName: "grapply"}, []int{5, 4}},
		{Config{MultiCols: [][]interface{}{{"corge", "bar"}, {"qux", "quuz"}}, MultiColsNames: []string{"grapply", "grault"}}, []int{5, 4}},
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
		Config{MultiCols: [][]interface{}{{"waldo", "waldo"}, {"d", "e"}}})
	excl := [][]bool{{false, true}, {false, false}}
	got := df.maxColWidths(excl)
	want := []int{5, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("df.maxColWidths() got %v, want %v for df \n%v", got, want, df)
	}
}

func TestMakeExclusionTable(t *testing.T) {
	df := MustNew([]interface{}{"foo", "bar"}, Config{MultiCols: [][]interface{}{{"baz", "qux"}, {"quux", "quuz"}}})
	got := df.makeExclusionsTable()
	want := [][]bool{{false, false}, {false, false}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("df.MakeExclusionsTable() got %v, want %v", got, want)
	}
}

func TestSubset(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo", "bar"}, []string{"baz", "qux"}},
		Config{Index: []string{"quux", "quuz"}})
	got := df.Subset([]int{1})
	want := MustNew([]interface{}{[]string{"bar"}, []string{"qux"}},
		Config{Index: []string{"quuz"}})
	if !Equal(got, want) {
		t.Errorf("df.Subset() got %v, want %v", got, want)
	}
}

func TestSubset_Empty(t *testing.T) {
	options.SetLogWarnings(false)
	df := MustNew([]interface{}{[]string{"foo", "bar"}})
	got := df.Subset([]int{})
	want := newEmptyDataFrame()
	if !Equal(got, want) {
		t.Errorf("df.Subset() got %v, want %v", got, want)
	}
	options.SetLogWarnings(true)
}
