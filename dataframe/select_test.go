package dataframe

import (
	"testing"

	"github.com/ptiger10/pd/series"
)

// func TestSubset(t *testing.T) {
// 	tests := []struct {
// 		args    []int
// 		want    *DataFrame
// 		wantErr bool
// 	}{
// 		{[]int{0}, MustNew([]interface{}{"foo"}), false},
// 		{[]int{1}, MustNew([]interface{}{"bar"}, Config{Index: 1}), false},
// 		{[]int{0, 1}, MustNew([]interface{}{[]string{"foo", "bar"}}), false},
// 		{[]int{1, 0}, MustNew([]interface{}{[]string{"bar", "foo"}}, Config{Index: []int{1, 0}}), false},
// 		{[]int{}, newEmptyDataFrame(), true},
// 		{[]int{3}, newEmptyDataFrame(), true},
// 	}
// 	for _, tt := range tests {
// 		df := MustNew([]interface{}{[]string{"foo", "bar", "baz"}})
// 		got, err := df.Subset(tt.args)
// 		if (err != nil) != tt.wantErr {
// 			t.Errorf("s.Subset() error = %v, want %v for args %v", err, tt.wantErr, tt.args)
// 		}
// 		if !Equal(got, tt.want) {
// 			t.Errorf("s.Subset() got %v, want %v for args %v", got, tt.want, tt.args)
// 		}
// 	}
// }

func TestSelectRows(t *testing.T) {
	var err error
	df, err := New(
		[]interface{}{[]string{"foo", "bar", "baz"}},
		Config{Index: []string{"qux", "quux", "corge"}, Col: []string{"foofoo"}})
	got, err := df.selectByRows([]int{0, 1})
	if err != nil {
		t.Errorf("selectByRows(): %v", err)
	}
	want := MustNew([]interface{}{[]string{"foo", "bar"}},
		Config{Index: []string{"qux", "quux"}, Col: []string{"foofoo"}})
	if !Equal(got, want) {
		t.Errorf("selectByRows(): got %v, want %v", got, want)
	}
}

func TestSelectCols(t *testing.T) {
	df := MustNew(
		[]interface{}{[]string{"foo"}, []string{"bar"}},
		Config{Col: []string{"baz", "qux"}})
	got, err := df.selectByCols([]int{1})
	if err != nil {
		t.Errorf("selectByCols(): %v", err)
	}
	want := MustNew([]interface{}{[]string{"bar"}}, Config{Col: []string{"qux"}})
	if !Equal(got, want) {
		t.Errorf("selectByCols(): got %v, want %v", got, want)
	}
}

func TestCol(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo"}, []string{"bar"}},
		Config{Col: []string{"baz", "qux"}})
	got := df.Col("qux")
	want := series.MustNew([]string{"bar"}, series.Config{Name: "qux"})
	if !series.Equal(got, want) {
		t.Errorf("Col(): got %v, want %v", got, want)
	}
}
