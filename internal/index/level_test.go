package index

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

func mustCreateNewLevel(data interface{}) Level {
	lvl, err := NewLevel(data, "")
	if err != nil {
		log.Fatalf("mustCreateNewLevel returned an error: %v", err)
	}
	return lvl
}

func TestLevel(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		data       interface{}
		wantLabels values.Values
		wantKind   kinds.Kind
		wantName   string
	}{
		{
			data:       []float64{0, 1, 2},
			wantLabels: mustCreateNewLevel([]float64{0, 1, 2}).Labels,
			wantKind:   kinds.Float,
			wantName:   "test",
		},
		{
			data:       []int{0, 1, 2},
			wantLabels: mustCreateNewLevel([]int64{0, 1, 2}).Labels,
			wantKind:   kinds.Int,
			wantName:   "test",
		},
		{
			data:       []uint{0, 1, 2},
			wantLabels: mustCreateNewLevel([]int64{0, 1, 2}).Labels,
			wantKind:   kinds.Int,
			wantName:   "test",
		},
		{
			data:       []string{"0", "1", "2"},
			wantLabels: mustCreateNewLevel([]string{"0", "1", "2"}).Labels,
			wantKind:   kinds.String,
			wantName:   "test",
		},
		{
			data:       []bool{true, true, false},
			wantLabels: mustCreateNewLevel([]bool{true, true, false}).Labels,
			wantKind:   kinds.Bool,
			wantName:   "test",
		},
		{
			data:       []time.Time{testDate},
			wantLabels: mustCreateNewLevel([]time.Time{testDate}).Labels,
			wantKind:   kinds.DateTime,
			wantName:   "test",
		},
		{
			data:       []interface{}{1.5, 1, "", false, testDate},
			wantLabels: mustCreateNewLevel([]interface{}{1.5, 1, "", false, testDate}).Labels,
			wantKind:   kinds.Interface,
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
		if lvl.Kind != test.wantKind {
			t.Errorf("%T test returned kind %v, want %v", test.data, lvl.Kind, test.wantKind)
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
