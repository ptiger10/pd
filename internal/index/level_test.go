package index

import (
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/datatypes"
	"github.com/ptiger10/pd/internal/values"
)

func TestLevel(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		data       interface{}
		wantLabels values.Values
		wantKind   datatypes.DataType
		wantName   string
	}{
		{
			data:       []float64{0, 1, 2},
			wantLabels: MustCreateNewLevel([]float64{0, 1, 2}, "").Labels,
			wantKind:   datatypes.Float64,
			wantName:   "test",
		},
		{
			data:       []int{0, 1, 2},
			wantLabels: MustCreateNewLevel([]int64{0, 1, 2}, "").Labels,
			wantKind:   datatypes.Int64,
			wantName:   "test",
		},
		{
			data:       []uint{0, 1, 2},
			wantLabels: MustCreateNewLevel([]int64{0, 1, 2}, "").Labels,
			wantKind:   datatypes.Int64,
			wantName:   "test",
		},
		{
			data:       []string{"0", "1", "2"},
			wantLabels: MustCreateNewLevel([]string{"0", "1", "2"}, "").Labels,
			wantKind:   datatypes.String,
			wantName:   "test",
		},
		{
			data:       []bool{true, true, false},
			wantLabels: MustCreateNewLevel([]bool{true, true, false}, "").Labels,
			wantKind:   datatypes.Bool,
			wantName:   "test",
		},
		{
			data:       []time.Time{testDate},
			wantLabels: MustCreateNewLevel([]time.Time{testDate}, "").Labels,
			wantKind:   datatypes.DateTime,
			wantName:   "test",
		},
		{
			data:       []interface{}{1.5, 1, "", false, testDate},
			wantLabels: MustCreateNewLevel([]interface{}{1.5, 1, "", false, testDate}, "").Labels,
			wantKind:   datatypes.Interface,
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
	if copyLvl.Name != "foo" {
		t.Error("Level.Copy() did not copy Name")
	}
	if copyLvl.DataType != datatypes.Int64 {
		t.Error("Level.Copy() did not copy Kind")
	}
	if copyLvl.Longest != 1 {
		t.Error("Level.Copy() did not copy Kind")
	}
	if reflect.ValueOf(idxLvl.Labels).Pointer() == reflect.ValueOf(copyLvl.Labels).Pointer() {
		t.Error("Level.Copy() returned original labels, want fresh copy")
	}
	if reflect.ValueOf(idxLvl.LabelMap).Pointer() == reflect.ValueOf(copyLvl.LabelMap).Pointer() {
		t.Error("Level.Copy() returned original map, want fresh copy")
	}
}
