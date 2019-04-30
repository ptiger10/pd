package constructors

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/new/datatypes"
	"github.com/ptiger10/pd/new/internal/index"
)

func TestIntIndexConstructor(t *testing.T) {
	var tests = []struct {
		data         interface{}
		wantKind     reflect.Kind
		wantLabels   index.IntLabels
		wantLabelMap map[string][]int
	}{
		{
			data:     []int{0, 1, 2},
			wantKind: datatypes.Int, wantLabels: index.IntLabels([]int64{0, 1, 2}),
			wantLabelMap: map[string][]int{"0": []int{0}, "1": []int{1}, "2": []int{2}},
		},
		{
			data:     []int{2, 6, 4, 6},
			wantKind: datatypes.Int, wantLabels: index.IntLabels([]int64{2, 6, 4, 6}),
			wantLabelMap: map[string][]int{"2": []int{0}, "4": []int{2}, "6": []int{1, 3}},
		},
		{
			data:     []int{},
			wantKind: datatypes.Int, wantLabels: nil,
			wantLabelMap: map[string][]int{},
		},
	}
	for _, test := range tests {
		got := SliceInt(test.data)
		if got.Kind != test.wantKind {
			t.Errorf("Data %v returned Kind %v, want %v", test.data, got.Kind, test.wantKind)
		}
		gotLabels := got.Labels.(index.IntLabels)
		if !reflect.DeepEqual(gotLabels, test.wantLabels) {
			t.Errorf("Data %v returned Labels %#v, want %#v", test.data, got.Labels, test.wantLabels)
		}
		if !reflect.DeepEqual(got.LabelMap, test.wantLabelMap) {
			t.Errorf("Data %v returned LabelMap %#v, want %#v", test.data, got.LabelMap, test.wantLabelMap)
		}
	}
}
