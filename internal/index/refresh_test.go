package index_test

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/internal/index"
	constructIdx "github.com/ptiger10/pd/internal/index/constructors"

	"github.com/ptiger10/pd/internal/values"
	constructVal "github.com/ptiger10/pd/internal/values/constructors"
)

func Test_RefreshLevel(t *testing.T) {
	var tests = []struct {
		newLabels    values.Values
		wantLabelMap index.LabelMap
		wantLongest  int
	}{
		{constructVal.SliceInt([]int64{3, 4}), index.LabelMap{"3": []int{0}, "4": []int{1}}, 1},
		{constructVal.SliceInt([]int64{10, 20}), index.LabelMap{"10": []int{0}, "20": []int{1}}, 2},
	}
	for _, test := range tests {
		lvl := constructIdx.SliceInt([]int64{1, 2}, "")
		origLabelMap := index.LabelMap{"1": []int{0}, "2": []int{1}}
		if !reflect.DeepEqual(lvl.LabelMap, origLabelMap) {
			t.Errorf("Returned labelMap %v, want %v", lvl.LabelMap, origLabelMap)
		}

		lvl.Labels = test.newLabels
		lvl.Refresh()
		if !reflect.DeepEqual(lvl.LabelMap, test.wantLabelMap) {
			t.Errorf("Returned labelMap %v, want %v", lvl.LabelMap, test.wantLabelMap)
		}
		if lvl.Longest != test.wantLongest {
			t.Errorf("Returned longest length %v, want %v", lvl.Longest, test.wantLongest)
		}
	}
}

func Test_RefreshIndex(t *testing.T) {
	var tests = []struct {
		newLevel    index.Level
		wantNameMap index.LabelMap
		wantName    string
	}{
		{constructIdx.SliceInt([]int64{1, 2}, "ints"), index.LabelMap{"ints": []int{0}}, "ints"},
	}
	for _, test := range tests {
		orig := constructIdx.SliceInt([]int64{1, 2}, "")
		idx := constructIdx.New(orig)
		idx.Levels[0] = test.newLevel
		idx.Refresh()
		if !reflect.DeepEqual(idx.NameMap, test.wantNameMap) {
			t.Errorf("Returned nameMap %v, want %v", idx.NameMap, test.wantNameMap)
		}
		if idx.Levels[0].Name != test.wantName {
			t.Errorf("Returned name %v, want %v", idx.Levels[0].Name, test.wantName)
		}
	}
}
