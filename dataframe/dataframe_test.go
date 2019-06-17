package dataframe

import (
	"fmt"
	"testing"

	"github.com/ptiger10/pd/series"
)

func TestDT(t *testing.T) {
	df, _ := New([]interface{}{"foo"})
	got := df.DT()
	want := series.MustNew("string", series.Idx("0"))
	if !series.Equal(got, want) {
		t.Errorf("df.DT() returned %v, want %v", got, want)
	}

}

func TestDataType(t *testing.T) {
	df, _ := New([]interface{}{"foo"})
	// fmt.Println(df.s[0])
	fmt.Println(df.DataType())
}
