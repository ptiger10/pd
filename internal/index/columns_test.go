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
	config := Config{Cols: []interface{}{"foo", "bar"}, ColsName: "baz"}
	got, err := NewColumnsFromConfig(config, 2)
	if err != nil {
		t.Errorf("NewColumnsFromConfig(): %v", err)
	}
	wantLvl := ColLevel{
		Name:     "baz",
		Labels:   []interface{}{"foo", "bar"},
		LabelMap: map[interface{}][]int{"foo": []int{0}, "bar": []int{1}}}
	want := Columns{NameMap: map[interface{}][]int{"baz": []int{0}},
		Levels: []ColLevel{wantLvl}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewColumnsFromConfig(): got %v, want %v", got.Levels[0], want.Levels[0])
	}
}

func TestColumns_In(t *testing.T) {
	tests := []struct {
		name      string
		positions []int
		want      Columns
		wantErr   bool
	}{
		{"pass 0", []int{0}, NewColumns(NewColLevel([]interface{}{"foo"}, "baz"), NewColLevel([]interface{}{"qux"}, "corge")), false},
		{"pass 1", []int{1}, NewColumns(NewColLevel([]interface{}{"bar"}, "baz"), NewColLevel([]interface{}{"quux"}, "corge")), false},
		{"out of range", []int{2}, Columns{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := NewColumns(
				NewColLevel([]interface{}{"foo", "bar"}, "baz"),
				NewColLevel([]interface{}{"qux", "quux"}, "corge"))
			got, err := col.In(tt.positions)
			if (err != nil) != tt.wantErr {
				t.Errorf("cols.In(): %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cols.In(): got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColLevel_In(t *testing.T) {
	tests := []struct {
		name      string
		positions []int
		want      ColLevel
		wantErr   bool
	}{
		{"pass 0", []int{0}, NewColLevel([]interface{}{"foo"}, "baz"), false},
		{"pass 1", []int{1}, NewColLevel([]interface{}{"bar"}, "baz"), false},
		{"out of range", []int{2}, ColLevel{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := NewColLevel([]interface{}{"foo", "bar"}, "baz")
			got, err := col.In(tt.positions)
			if (err != nil) != tt.wantErr {
				t.Errorf("colsLevel.In(): %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("colsLevel.In(): got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColMaxNameWidth(t *testing.T) {
	cols := NewColumns([]ColLevel{ColLevel{Name: "foo", Labels: []interface{}{"bar"}}, ColLevel{Name: "corge", Labels: []interface{}{"quux"}}}...)
	got := cols.MaxNameWidth()
	want := 5
	if got != want {
		t.Errorf("cols.MaxNameWidth() got %v, want %v", got, want)
	}
}
