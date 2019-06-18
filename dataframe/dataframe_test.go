package dataframe

import (
	"fmt"
	"testing"

	"github.com/ptiger10/pd/series"
)

func TestDataTypes(t *testing.T) {
	df, err := New([]interface{}{"foo"})
	if err != nil {
		t.Errorf("df.DT(): %v", err)
	}
	got := df.DataTypes()
	want, err := series.New("string", series.Config{Index: "0", Name: "datatypes"})
	if err != nil {
		fmt.Println(err)
	}
	if !series.Equal(got, want) {
		t.Errorf("df.DT() returned %v, want %v", got, want)
	}
}

func TestRowsIn(t *testing.T) {
	var err error
	df, _ := New(
		[]interface{}{[]string{"foo", "bar", "baz"}, []string{"qux", "quux", "corge"}},
		Config{Cols: []interface{}{"foofoo", "barbar"}})
	df, err = df.rowsIn([]int{0, 1})
	if err != nil {
		t.Errorf("rowsIn(): %v", err)
	}
	fmt.Println(df)
}

func TestColsIn(t *testing.T) {
	var err error
	df, _ := New(
		[]interface{}{[]string{"foo", "bar", "baz"}, []string{"qux", "quux", "corge"}},
		Config{Cols: []interface{}{"foofoo", "barbar"}})
	df, err = df.colsIn([]int{1})
	if err != nil {
		t.Errorf("colsIn(): %v", err)
	}
	fmt.Println(df)
}

func TestDataType(t *testing.T) {
	df, _ := New([]interface{}{"foo"})
	got := df.dataType()
	want := "string"
	if got != want {
		t.Errorf("df.dataType() returned %v, want %v", got, want)
	}
}
