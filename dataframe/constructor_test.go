package dataframe

import (
	"fmt"
	"testing"

	"github.com/ptiger10/pd/series"
)

func Test_New(t *testing.T) {
	c := Config{Columns: []string{"fooCol", "barCol"}}
	df, err := NewWithConfig(c, []interface{}{[]float64{1, 2, 3}, []float64{4, 5, 6}}, series.Idx([]string{"foo", "bar", "baz"}))
	if err != nil {
		t.Error(err)
	}
	fmt.Print(df)
	// fmt.Println(df.Sum())
}
