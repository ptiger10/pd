package dataframe

import (
	"testing"

	"github.com/ptiger10/pd/series"
)

func TestDT(t *testing.T) {
	df, err := New([]interface{}{"foo"})
	if err != nil {
		t.Errorf("df.DT(): %v", err)
	}
	got := df.DT()
	want := series.MustNew("string", series.Config{Index: "0", Name: "datatypes"})
	if !series.Equal(got, want) {
		t.Errorf("df.DT() returned %v, want %v", got, want)
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
