package dataframe

import (
	"fmt"
	"testing"

	"github.com/ptiger10/pd/series"
)

func Test_New(t *testing.T) {
	df, err := New([]interface{}{[]float64{1, 2, 3}, []float64{4, 5, 6}}, series.Idx([]string{"foo", "bar"}))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(df.Sum())
}
