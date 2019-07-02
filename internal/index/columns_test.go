package index

import (
	"reflect"
	"testing"
)

func TestNewColumns(t *testing.T) {
	type args struct {
		levels []ColLevel
	}
	type want struct {
		columns      Columns
		len          int
		numLevels    int
		maxNameWidth int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"empty", args{nil},
			want{
				columns: Columns{Levels: []ColLevel{ColLevel{Labels: nil}}, NameMap: LabelMap{"": []int{0}}},
				len:     0, numLevels: 1, maxNameWidth: 0,
			}},
		{"one col",
			args{[]ColLevel{NewColLevel([]string{"1", "2"}, "foo")}},
			want{Columns{
				Levels:  []ColLevel{ColLevel{Name: "foo", LabelMap: LabelMap{"1": []int{0}, "2": []int{1}}, Labels: []string{"1", "2"}}},
				NameMap: LabelMap{"foo": []int{0}}},
				2, 1, 3,
			}},
		{"two cols",
			args{[]ColLevel{NewColLevel([]string{"1", "2"}, "foo"), NewColLevel([]string{"3", "4"}, "corge")}},
			want{Columns{
				Levels: []ColLevel{
					ColLevel{Name: "foo", LabelMap: LabelMap{"1": []int{0}, "2": []int{1}}, Labels: []string{"1", "2"}},
					ColLevel{Name: "corge", LabelMap: LabelMap{"3": []int{0}, "4": []int{1}}, Labels: []string{"3", "4"}}},
				NameMap: LabelMap{"foo": []int{0}, "corge": []int{1}}},
				2, 2, 5,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewColumns(tt.args.levels...)
			if !reflect.DeepEqual(got, tt.want.columns) {
				t.Errorf("NewColumns(): got %v, want %v", got, tt.want.columns)
			}
			gotLen := got.Len()
			if gotLen != tt.want.len {
				t.Errorf("Columns.Len(): got %v, want %v", gotLen, tt.want.len)
			}
			gotLevels := got.NumLevels()
			if gotLevels != tt.want.numLevels {
				t.Errorf("Columns.NumLevels(): got %v, want %v", gotLevels, tt.want.numLevels)
			}
			gotMaxWidth := got.MaxNameWidth()
			if gotMaxWidth != tt.want.maxNameWidth {
				t.Errorf("Columns.MaxWidth(): got %v, want %v", gotMaxWidth, tt.want.maxNameWidth)
			}
		})
	}
}

func TestNewDefaultColumns(t *testing.T) {
	got := NewDefaultColumns(3)
	want := NewColumns(NewColLevel([]string{"0", "1", "2"}, ""))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewDefaultColumns: got %v, want %v", got, want)
	}
}

func TestNewColLevel(t *testing.T) {
	got := NewColLevel([]string{"foo", "bar"}, "foobar")
	want := ColLevel{
		Name:     "foobar",
		Labels:   []string{"foo", "bar"},
		LabelMap: LabelMap{"foo": []int{0}, "bar": []int{1}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewColLevel(): got %v, want %v", got, want)
	}
}

func TestNewColumnFromConfig(t *testing.T) {
	type args struct {
		config Config
		n      int
	}
	type want struct {
		columns Columns
		err     bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"both nil and unnamed",
			args{Config{}, 2},
			want{NewDefaultColumns(2), false}},
		{"both nil but named",
			args{Config{ColName: "foo"}, 2},
			want{NewColumns(NewDefaultColLevel(2, "foo")), false}},
		{"singleLevel",
			args{Config{Col: []string{"foo", "bar"}, ColName: "baz"}, 2},
			want{NewColumns(NewColLevel([]string{"foo", "bar"}, "baz")), false}},
		{"multiLevel",
			args{Config{MultiCol: [][]string{{"foo", "bar"}, {"baz", "qux"}}, MultiColNames: []string{"quux", "quuz"}}, 2},
			want{NewColumns(NewColLevel([]string{"foo", "bar"}, "quux"), NewColLevel([]string{"baz", "qux"}, "quuz")), false}},
		{"fail: both not nil",
			args{Config{
				Col:      []string{"foo"},
				MultiCol: [][]string{{"foo"}}}, 2},
			want{Columns{}, true}},
		{"fail: wrong multiindex names length",
			args{Config{
				MultiCol:      [][]string{{"foo"}},
				MultiColNames: []string{"foo", "bar"}}, 2},
			want{Columns{}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewColumnsFromConfig(tt.args.config, tt.args.n)
			if (err != nil) != tt.want.err {
				t.Errorf("NewColumnsFromConfig() error = %v, want %v", err, tt.want.err)
			}
			if !reflect.DeepEqual(got, tt.want.columns) {
				t.Errorf("NewColumnsFromConfig(): got %v, want %v", got, tt.want.columns)
			}
		})
	}
}

func TestColumns_Nil(t *testing.T) {
	cols := Columns{}
	_ = cols.Copy()
	_ = cols.Len()
	_ = cols.NumLevels()
	cols.Refresh()
}

func TestCol_Refresh(t *testing.T) {
	columns := NewColumns(NewColLevel([]string{"foo"}, "bar"))
	columns.Levels[0] = NewDefaultColLevel(5, "baz")
	columns.Refresh()
	want := NewColumns(NewDefaultColLevel(5, "baz"))
	if !reflect.DeepEqual(columns, want) {
		t.Errorf("Col.Refresh() got %v, want %v", columns, want)
	}
	// Empty or nil columns do not trigger an error
	columns.Levels = make([]ColLevel, 0)
	columns.Refresh()
	columns = NewColumns()
	columns.Refresh()

}

func TestCol_Subset(t *testing.T) {
	tests := []struct {
		name      string
		positions []int
		want      Columns
		wantErr   bool
	}{
		{"pass 0", []int{0}, NewColumns(NewColLevel([]string{"foo"}, "baz"), NewColLevel([]string{"qux"}, "corge")), false},
		{"pass 1", []int{1}, NewColumns(NewColLevel([]string{"bar"}, "baz"), NewColLevel([]string{"quux"}, "corge")), false},
		{"fail: out of range", []int{2}, Columns{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := NewColumns(
				NewColLevel([]string{"foo", "bar"}, "baz"),
				NewColLevel([]string{"qux", "quux"}, "corge"))
			got, err := col.Subset(tt.positions)
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

func TestColLevel_Subset(t *testing.T) {
	tests := []struct {
		name      string
		positions []int
		want      ColLevel
		wantErr   bool
	}{
		{"pass 0", []int{0}, NewColLevel([]string{"foo"}, "baz"), false},
		{"pass 1", []int{1}, NewColLevel([]string{"bar"}, "baz"), false},
		{"out of range", []int{2}, ColLevel{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := NewColLevel([]string{"foo", "bar"}, "baz")
			got, err := col.Subset(tt.positions)
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
