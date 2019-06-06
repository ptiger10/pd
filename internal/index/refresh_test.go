package index

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/internal/values"
)

func Test_RefreshLevel(t *testing.T) {
	var tests = []struct {
		newLabels    values.Values
		wantLabelMap LabelMap
		wantLongest  int
	}{
		{MustCreateNewLevel([]int64{3, 4}, "").Labels, LabelMap{"3": []int{0}, "4": []int{1}}, 1},
		{MustCreateNewLevel([]int64{10, 20}, "").Labels, LabelMap{"10": []int{0}, "20": []int{1}}, 2},
	}
	for _, test := range tests {
		lvl := MustCreateNewLevel([]int64{1, 2}, "")
		origLabelMap := LabelMap{"1": []int{0}, "2": []int{1}}
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
	origLvl, err := NewLevel([]int64{1, 2}, "")
	if err != nil {
		t.Error(err)
	}
	idx := New(origLvl)
	if idx.Levels[0].Name != "" {
		t.Error("Expecting no name")
	}
	newLvl, err := NewLevel([]int64{1, 2}, "ints")
	if err != nil {
		t.Error(err)
	}
	idx.Levels[0] = newLvl
	idx.Refresh()
	wantNameMap := LabelMap{"ints": []int{0}}
	wantName := "ints"
	if !reflect.DeepEqual(idx.NameMap, wantNameMap) {
		t.Errorf("Returned nameMap %v, want %v", idx.NameMap, wantNameMap)
	}
	if idx.Levels[0].Name != wantName {
		t.Errorf("Returned name %v, want %v", idx.Levels[0].Name, wantName)
	}
}
