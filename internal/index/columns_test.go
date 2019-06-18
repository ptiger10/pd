package index

import (
	"reflect"
	"testing"
)

func TestNewColLen(t *testing.T) {
	got := ColLevel{Labels: []interface{}{"foo", "bar"}}.Len()
	want := 2
	if got != want {
		t.Errorf("ColLevel.Len(): got %v, want %v", got, want)
	}
}

func TestNewColLevel(t *testing.T) {
	got := NewColLevel([]interface{}{"foo", "bar"}, "foobar")
	want := ColLevel{
		Name:     "foobar",
		Labels:   []interface{}{"foo", "bar"},
		LabelMap: LabelMap{"foo": []int{0}, "bar": []int{1}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewColLevel(): got %v, want %v", got, want)
	}
}

func TestNewColumnFromConfig(t *testing.T) {
	config := Config{Cols: []interface{}{"foo", "bar"}}
	got, err := NewColumnsFromConfig(2, config)
	if err != nil {
		t.Errorf("NewColumnsFromConfig(): %v", err)
	}
	want := Columns{NameMap: map[string][]int{"": []int{0}},
		Levels: []ColLevel{ColLevel{
			Labels:   []interface{}{"foo", "bar"},
			LabelMap: map[string][]int{"foo": []int{0}, "bar": []int{1}}}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewColumnsFromConfig(): got %v, want %v", got.Levels[0], want.Levels[0])
	}
}
