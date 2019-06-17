package dataframe

import (
	"fmt"
	"testing"
)

func Test_New(t *testing.T) {
	c := Config{Cols: []interface{}{"fooCol", "barCol"}, Index: []string{"foo", "bar", "baz"}}
	df, err := New([]interface{}{[]int64{1, 2, 3}, []float64{4, 5, 6}}, c)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(df.DT())
	// fmt.Print(df)
	// fmt.Println(df.Sum())
}
