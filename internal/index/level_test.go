package index

import (
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

func TestLevel(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		data       interface{}
		wantLabels values.Values
		wantKind   options.DataType
		wantName   string
	}{
		{
			data:       []float64{0, 1, 2},
			wantLabels: MustCreateNewLevel([]float64{0, 1, 2}, "").Labels,
			wantKind:   options.Float64,
			wantName:   "test",
		},
		{
			data:       []int{0, 1, 2},
			wantLabels: MustCreateNewLevel([]int64{0, 1, 2}, "").Labels,
			wantKind:   options.Int64,
			wantName:   "test",
		},
		{
			data:       []uint{0, 1, 2},
			wantLabels: MustCreateNewLevel([]int64{0, 1, 2}, "").Labels,
			wantKind:   options.Int64,
			wantName:   "test",
		},
		{
			data:       []string{"0", "1", "2"},
			wantLabels: MustCreateNewLevel([]string{"0", "1", "2"}, "").Labels,
			wantKind:   options.String,
			wantName:   "test",
		},
		{
			data:       []bool{true, true, false},
			wantLabels: MustCreateNewLevel([]bool{true, true, false}, "").Labels,
			wantKind:   options.Bool,
			wantName:   "test",
		},
		{
			data:       []time.Time{testDate},
			wantLabels: MustCreateNewLevel([]time.Time{testDate}, "").Labels,
			wantKind:   options.DateTime,
			wantName:   "test",
		},
		{
			data:       []interface{}{1.5, 1, "", false, testDate},
			wantLabels: MustCreateNewLevel([]interface{}{1.5, 1, "", false, testDate}, "").Labels,
			wantKind:   options.Interface,
			wantName:   "test",
		},
	}
	for _, test := range tests {
		lvl, err := NewLevel(test.data, "test")
		if err != nil {
			t.Errorf("Unable to construct level from %v: %v", test.data, err)
		}
		if !reflect.DeepEqual(lvl.Labels, test.wantLabels) {
			t.Errorf("%T test returned labels %#v, want %#v", test.data, lvl.Labels, test.wantLabels)
		}
		if lvl.DataType != test.wantKind {
			t.Errorf("%T test returned kind %v, want %v", test.data, lvl.DataType, test.wantKind)
		}
		if lvl.Name != test.wantName {
			t.Errorf("%T test returned name %v, want %v", test.data, lvl.Name, test.wantName)
		}
	}
}

func TestLevel_Unsupported(t *testing.T) {
	data := []complex64{1, 2, 3}
	_, err := NewLevel(data, "")
	if err == nil {
		t.Errorf("Returned nil error, expected error due to unsupported type %T", data)
	}
}

func Test_LevelCopy(t *testing.T) {
	idxLvl := MustCreateNewLevel([]int{1, 2, 3}, "")
	idxLvl.Name = "foo"
	copyLvl := idxLvl.Copy()
	if !reflect.DeepEqual(idxLvl.LabelMap, copyLvl.LabelMap) {
		t.Error("Level.Copy() did not copy LabelMap")
	}
	if copyLvl.Name != "foo" {
		t.Error("Level.Copy() did not copy Name")
	}
	if copyLvl.DataType != options.Int64 {
		t.Error("Level.Copy() did not copy Kind")
	}
	if copyLvl.MaxWidth() != 1 {
		t.Error("Level.Copy() did not copy Kind")
	}
	if reflect.ValueOf(idxLvl.Labels).Pointer() == reflect.ValueOf(copyLvl.Labels).Pointer() {
		t.Error("Level.Copy() returned original labels, want fresh copy")
	}
	if reflect.ValueOf(idxLvl.LabelMap).Pointer() == reflect.ValueOf(copyLvl.LabelMap).Pointer() {
		t.Error("Level.Copy() returned original map, want fresh copy")
	}
}

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
		if lvl.MaxWidth() != test.wantLongest {
			t.Errorf("Returned longest length %v, want %v", lvl.MaxWidth(), test.wantLongest)
		}
	}
}
