package dataframe

import (
	"fmt"
	"reflect"
	"testing"

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

func TestEqual(t *testing.T) {
	df := MustNew(
		[]interface{}{[]string{"foo"}, []string{"bar"}},
		Config{Col: []string{"baz", "qux"}})
	df2 := MustNew(
		[]interface{}{[]string{"foo"}, []string{"bar"}},
		Config{Col: []string{"baz", "qux"}})
	if !Equal(df, df2) {
		t.Errorf("Equal() did not return true for equivalent df")
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
