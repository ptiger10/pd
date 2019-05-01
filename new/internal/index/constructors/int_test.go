package constructors

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/new/internal/index"
	"github.com/ptiger10/pd/new/internal/values"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

func TestLevelConstructor_int_slice(t *testing.T) {
	var tests = []struct {
		data         interface{}
		name         string
		wantKind     reflect.Kind
		wantLabels   values.Values
		wantLabelMap index.LabelMap
	}{
		{
			data:         []int8{0, 1, 2},
			name:         "",
			wantKind:     kinds.Int,
			wantLabels:   constructVal.SliceInt([]int{0, 1, 2}),
			wantLabelMap: index.LabelMap{"0": []int{0}, "1": []int{1}, "2": []int{2}},
		},
		{
			data:         []int32{2, 4, 6, 4},
			name:         "test",
			wantKind:     kinds.Int,
			wantLabels:   constructVal.SliceInt([]int{2, 4, 6, 4}),
			wantLabelMap: index.LabelMap{"2": []int{0}, "4": []int{1, 3}, "6": []int{2}},
		},
		{
			data:         []int{},
			wantKind:     kinds.Int,
			wantLabels:   values.IntValues(nil),
			wantLabelMap: index.LabelMap{},
		},
	}
	for _, test := range tests {
		got := SliceInt(test.data, test.name)
		if got.Kind != test.wantKind {
			t.Errorf("Returned Kind %v, want %v\n", got.Kind, test.wantKind)
		}
		if got.Name != test.name {
			t.Errorf("Returned Name %v, want %v\n", got.Name, test.name)
		}
		gotLabels := got.Labels.(values.Values)
		if !reflect.DeepEqual(gotLabels, test.wantLabels) {
			t.Errorf("Data %v returned Labels %#v, want %#v\n", test.data, got.Labels, test.wantLabels)
		}
		if !reflect.DeepEqual(got.LabelMap, test.wantLabelMap) {
			t.Errorf("Data %v returned LabelMap %#v, want %#v\n", test.data, got.LabelMap, test.wantLabelMap)
		}
	}
}
