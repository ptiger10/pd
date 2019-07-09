package index

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

func TestIndex_Nil(t *testing.T) {
	idx := Index{}
	_ = idx.Aligned()
	_ = idx.Copy()
	_ = idx.Len()
	_ = idx.NumLevels()
	idx.Refresh()
}

func TestIndexLevel_Nil(t *testing.T) {
	lvl := Level{}
	_ = lvl.Copy()
	_ = lvl.Len()
	_ = lvl.maxWidth()
	lvl.Refresh()
}

func TestNew(t *testing.T) {
	vals, _ := values.InterfaceFactory([]int{1, 2})
	labels := vals.Values

	type args struct {
		levels []Level
	}
	type want struct {
		index     Index
		len       int
		numLevels int
		maxWidths []int
		unnamed   bool
		datatypes []options.DataType
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"empty", args{nil},
			want{
				index: Index{Levels: []Level{}, NameMap: LabelMap{}},
				len:   0, numLevels: 0, maxWidths: []int{}, unnamed: true, datatypes: []options.DataType{},
			}},
		{"one level",
			args{[]Level{MustNewLevel([]int{1, 2}, "foo")}},
			want{
				index: Index{
					Levels:  []Level{Level{Name: "foo", DataType: options.Int64, LabelMap: LabelMap{"1": []int{0}, "2": []int{1}}, Labels: labels}},
					NameMap: LabelMap{"foo": []int{0}}},
				len: 2, numLevels: 1, maxWidths: []int{3}, unnamed: false, datatypes: []options.DataType{options.Int64},
			}},
		{"two levels",
			args{[]Level{MustNewLevel([]int{1, 2}, "foo"), MustNewLevel([]int{1, 2}, "corge")}},
			want{
				index: Index{
					Levels: []Level{
						Level{Name: "foo", DataType: options.Int64, LabelMap: LabelMap{"1": []int{0}, "2": []int{1}}, Labels: labels},
						Level{Name: "corge", DataType: options.Int64, LabelMap: LabelMap{"1": []int{0}, "2": []int{1}}, Labels: labels}},
					NameMap: LabelMap{"foo": []int{0}, "corge": []int{1}}},
				len: 2, numLevels: 2, maxWidths: []int{3, 5}, unnamed: false, datatypes: []options.DataType{options.Int64, options.Int64},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.levels...)
			if !reflect.DeepEqual(got, tt.want.index) {
				t.Errorf("New(): got %#v, want %#v", got, tt.want.index)
			}
			gotLen := got.Len()
			if gotLen != tt.want.len {
				t.Errorf("Index.Len(): got %v, want %v", gotLen, tt.want.len)
			}
			gotNumLevels := got.NumLevels()
			if !reflect.DeepEqual(gotNumLevels, tt.want.numLevels) {
				t.Errorf("Index.NumLevels(): got %v, want %v", gotNumLevels, tt.want.numLevels)
			}
			gotMaxWidth := got.MaxWidths()
			if !reflect.DeepEqual(gotMaxWidth, tt.want.maxWidths) {
				t.Errorf("Index.MaxWidth(): got %#v, want %#v", gotMaxWidth, tt.want.maxWidths)
			}
			gotUnnamed := got.Unnamed()
			if !reflect.DeepEqual(gotUnnamed, tt.want.unnamed) {
				t.Errorf("Index.GotUnnamed(): got %v, want %v", gotUnnamed, tt.want.unnamed)
			}
			gotDataTypes := got.DataTypes()
			if !reflect.DeepEqual(gotDataTypes, tt.want.datatypes) {
				t.Errorf("Index.GotDataTypes(): got %#v, want %#v", gotDataTypes, tt.want.datatypes)
			}
		})
	}
}

func TestNewFromConfig(t *testing.T) {
	type args struct {
		config Config
		n      int
	}
	type want struct {
		index Index
		err   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"both nil and unnamed",
			args{Config{}, 2},
			want{NewDefault(2), false}},
		{"both nil but named",
			args{Config{IndexName: "foo"}, 2},
			want{New(NewDefaultLevel(2, "foo")), false}},
		{"singleIndex",
			args{Config{Index: []string{"foo", "bar"}, IndexName: "baz"}, 2},
			want{New(MustNewLevel([]string{"foo", "bar"}, "baz")), false}},
		{"multiIndex",
			args{Config{MultiIndex: []interface{}{[]string{"foo", "bar"}, []string{"baz", "qux"}}, MultiIndexNames: []string{"quux", "quuz"}}, 2},
			want{New(MustNewLevel([]string{"foo", "bar"}, "quux"), MustNewLevel([]string{"baz", "qux"}, "quuz")), false}},
		{"singleIndex interpolated",
			args{Config{Index: []interface{}{"foo", "bar"}, IndexName: "baz"}, 2},
			want{New(MustNewLevel([]string{"foo", "bar"}, "baz")), false}},
		{"singleIndex not interpolated",
			args{Config{Manual: true, Index: []interface{}{"foo", "bar"}, IndexName: "baz"}, 2},
			want{New(MustNewLevel([]interface{}{"foo", "bar"}, "baz")), false}},
		{"multiIndex interpolated",
			args{Config{MultiIndex: []interface{}{[]interface{}{"foo", "bar"}}, MultiIndexNames: []string{"baz"}}, 2},
			want{New(MustNewLevel([]string{"foo", "bar"}, "baz")), false}},
		{"multiIndex not interpolated",
			args{Config{Manual: true, MultiIndex: []interface{}{[]interface{}{"foo", "bar"}}, MultiIndexNames: []string{"baz"}}, 2},
			want{New(MustNewLevel([]interface{}{"foo", "bar"}, "baz")), false}},
		{"fail: singleIndex unsupported type",
			args{Config{Index: complex64(1)}, 2},
			want{Index{}, true}},
		{"fail: multiIndex unsupported type",
			args{Config{MultiIndex: []interface{}{complex64(1)}}, 2},
			want{Index{}, true}},
		{"fail: both not nil",
			args{Config{
				Index:      "foo",
				MultiIndex: []interface{}{"foo"}}, 2},
			want{Index{}, true}},
		{"fail: wrong multiindex names length",
			args{Config{
				MultiIndex:      []interface{}{"foo"},
				MultiIndexNames: []string{"bar", "baz"}}, 2},
			want{Index{}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFromConfig(tt.args.config, tt.args.n)
			if (err != nil) != tt.want.err {
				t.Errorf("NewFromConfig() error = %v, want %v", err, tt.want.err)
			}
			if !reflect.DeepEqual(got, tt.want.index) {
				t.Errorf("NewFromConfig(): got %v, want %v", got, tt.want.index)
			}
		})
	}
}

func Test_NewDefault(t *testing.T) {
	got := NewDefault(2)
	want := New(Level{Labels: values.MustCreateValuesFromInterface([]int64{0, 1}).Values,
		LabelMap: LabelMap{"0": []int{0}, "1": []int{1}}, Name: "", DataType: options.Int64, IsDefault: true})
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Default constructor returned %v, want %v", got, want)
	}
	gotLen := len(got.Levels)
	wantLen := 1
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}
}

func TestIndex_NewLevel(t *testing.T) {
	type args struct {
		data interface{}
		name string
	}
	type want struct {
		level Level
		err   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "empty", args: args{data: nil, name: ""},
			want: want{level: Level{}, err: false}},
		{"pass", args{"foo", "bar"},
			want{
				Level{
					Name: "bar", DataType: options.String,
					Labels: values.MustCreateValuesFromInterface("foo").Values, LabelMap: LabelMap{"foo": []int{0}}},
				false}},
		{"[]interface no interpolation", args{[]interface{}{"foo"}, ""},
			want{
				Level{
					DataType: options.Interface,
					Labels:   values.MustCreateValuesFromInterface([]interface{}{"foo"}).Values, LabelMap: LabelMap{"foo": []int{0}}},
				false}},
		{"fail unsupported", args{complex64(1), "bar"},
			want{Level{}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLevel(tt.args.data, tt.args.name)
			if (err != nil) != tt.want.err {
				t.Errorf("NewLevel() error = %v, want %v", err, tt.want.err)
			}
			if !reflect.DeepEqual(got, tt.want.level) {
				t.Errorf("NewLevel() = %v, want %v", got, tt.want.level)
			}
		})
	}
}

func TestIndex_InterpolatedNewLevel(t *testing.T) {
	type args struct {
		data interface{}
		name string
	}
	type want struct {
		level Level
		err   bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "empty", args: args{data: nil, name: ""},
			want: want{level: Level{}, err: false}},
		{name: "no interpolation", args: args{data: "foo", name: ""},
			want: want{level: Level{DataType: options.String,
				Labels: values.MustCreateValuesFromInterface("foo").Values, LabelMap: LabelMap{"foo": []int{0}}},
				err: false}},
		{name: "interpolated string", args: args{data: []interface{}{"foo"}, name: ""},
			want: want{level: Level{DataType: options.String,
				Labels: values.MustCreateValuesFromInterface("foo").Values, LabelMap: LabelMap{"foo": []int{0}}},
				err: false}},
		{"fail unsupported", args{complex64(1), "bar"},
			want{Level{}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InterpolatedNewLevel(tt.args.data, tt.args.name)
			if (err != nil) != tt.want.err {
				t.Errorf("InterpolatedNewLevel() error = %v, want %v", err, tt.want.err)
			}
			if !reflect.DeepEqual(got, tt.want.level) {
				t.Errorf("InterpolatedNewLevel() = %v, want %v", got, tt.want.level)
			}
		})
	}
}

func TestNewIndexLevel_Copy(t *testing.T) {
	tests := []struct {
		name  string
		input Level
		want  Level
	}{
		{name: "empty nil", input: MustNewLevel(nil, ""), want: Level{}},
		{"pass", MustNewLevel("foo", "bar"), MustNewLevel("foo", "bar")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Copy()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Level.Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIndex_Copy(t *testing.T) {
	tests := []struct {
		name  string
		input Index
		want  Index
	}{
		{name: "empty", input: New(), want: Index{Levels: []Level{}, NameMap: LabelMap{}}},
		{"pass", New(MustNewLevel("foo", "bar")), New(MustNewLevel("foo", "bar"))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Copy()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index.Copy() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestMustNew_fail(t *testing.T) {
	options.RestoreDefaults()
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	MustNewLevel(complex64(1), "")
	if buf.String() == "" {
		t.Errorf("MustNew() returned no log message, want log due to fail")
	}
}

func TestElements(t *testing.T) {
	idx := New(MustNewLevel([]string{"foo", "bar", "baz"}, "a"), MustNewLevel([]int64{1, 2, 3}, "b"))
	got := idx.Elements(0)
	want := Elements{Labels: []interface{}{"foo", int64(1)}, DataTypes: []options.DataType{options.String, options.Int64}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Elements(): got %#v, want %#v", got, want)
	}
}

func TestAligned(t *testing.T) {
	vals, _ := values.InterfaceFactory([]int{1})
	labels1 := vals.Values
	vals2, _ := values.InterfaceFactory([]int{1, 2})
	labels2 := vals2.Values

	tests := []struct {
		name  string
		input Index
		err   bool
	}{
		{"aligned", Index{Levels: []Level{Level{Labels: labels1}, Level{Labels: labels1}}}, false},
		{"fail: misaligned", Index{Levels: []Level{Level{Labels: labels1}, Level{Labels: labels2}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Aligned()
			if (err != nil) != tt.err {
				t.Errorf("Aligned() got %v, want %v", err, tt.err)
			}
		})
	}
}

func TestIndex_SwapLevels(t *testing.T) {
	idx := New(MustNewLevel("foo", ""), MustNewLevel("bar", ""))
	type args struct {
		i int
		j int
	}
	type want struct {
		index Index
		err   bool
	}
	tests := []struct {
		name  string
		input Index
		args  args
		want  want
	}{
		{name: "pass", input: idx.Copy(), args: args{0, 1},
			want: want{index: New(MustNewLevel("bar", ""), MustNewLevel("foo", "")), err: false}},
		{"reverse order", idx.Copy(), args{1, 0},
			want{New(MustNewLevel("bar", ""), MustNewLevel("foo", "")), false}},
		{"fail: i", idx.Copy(), args{10, 1},
			want{idx, true}},
		{"fail: j", idx.Copy(), args{0, 10},
			want{idx, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.SwapLevels(tt.args.i, tt.args.j)
			if (err != nil) != tt.want.err {
				t.Errorf("Index.SwapLevels() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(tt.input, tt.want.index) {
				t.Errorf("Index.SwapLevels() = \n%#v, want \n%#v", tt.input, tt.want.index)
			}
		})
	}
}

func TestIndex_InsertLevel(t *testing.T) {
	idx := New(MustNewLevel("foo", ""), MustNewLevel("bar", ""))
	type args struct {
		pos    int
		values interface{}
		name   string
	}
	type want struct {
		index Index
		err   bool
	}
	tests := []struct {
		name  string
		input Index
		args  args
		want  want
	}{
		{name: "0", input: idx.Copy(), args: args{0, "baz", ""},
			want: want{index: New(MustNewLevel("baz", ""), MustNewLevel("foo", ""), MustNewLevel("bar", "")), err: false}},
		{"1", idx.Copy(), args{1, "baz", ""},
			want{index: New(MustNewLevel("foo", ""), MustNewLevel("baz", ""), MustNewLevel("bar", "")), err: false}},
		{"1", idx.Copy(), args{2, "baz", ""},
			want{index: New(MustNewLevel("foo", ""), MustNewLevel("bar", ""), MustNewLevel("baz", "")), err: false}},
		{"fail: invalid position", idx.Copy(), args{10, "baz", ""},
			want{idx.Copy(), true}},
		{"fail: unsupported value", idx.Copy(), args{2, complex64(1), ""},
			want{idx.Copy(), true}},
		{"fail: incorrect length of values", idx.Copy(), args{2, []string{"corge", "waldo"}, ""},
			want{idx.Copy(), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.InsertLevel(tt.args.pos, tt.args.values, tt.args.name)
			if (err != nil) != tt.want.err {
				t.Errorf("Index.InsertLevel() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(tt.input, tt.want.index) {
				t.Errorf("Index.InsertLevel() = %v, want %v", tt.input, tt.want.index)
			}
		})
	}
}

func TestSubset(t *testing.T) {
	lvl := MustNewLevel([]string{"foo", "bar", "baz"}, "")
	type args struct {
		pos []int
	}
	type want struct {
		index Index
		err   bool
	}
	tests := []struct {
		name  string
		input Index
		args  args
		want  want
	}{
		{name: "subsetRows multiIndex",
			input: New(lvl, lvl),
			args:  args{pos: []int{0, 1}},
			want: want{index: New(MustNewLevel([]string{"foo", "bar"}, ""), MustNewLevel([]string{"foo", "bar"}, "")),
				err: false}},
		{"subsetRows singleIndex",
			New(lvl),
			args{[]int{0, 1}},
			want{New(MustNewLevel([]string{"foo", "bar"}, "")), false}},
		{"fail: invalid row",
			New(lvl),
			args{[]int{10}},
			want{New(lvl), true}},
		{"fail: no rows",
			New(lvl),
			args{[]int{}},
			want{New(lvl), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Subset(tt.args.pos)
			if (err != nil) != tt.want.err {
				t.Errorf("Index.Subset() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(tt.input, tt.want.index) {
				t.Errorf("Subset(): got %v, want %v", tt.input, tt.want.index)
			}
		})
	}
}

func TestSubsetLevels(t *testing.T) {
	idx := New(MustNewLevel([]string{"foo", "bar", "baz"}, ""), MustNewLevel([]string{"qux", "quux", "quuz"}, ""))
	type args struct {
		pos []int
	}
	type want struct {
		index Index
		err   bool
	}
	tests := []struct {
		name  string
		input Index
		args  args
		want  want
	}{

		{name: "subsetLevels multiIndex",
			input: idx,
			args:  args{pos: []int{1}},
			want:  want{New(MustNewLevel([]string{"qux", "quux", "quuz"}, "")), false}},
		{"fail: invalid level levels",
			idx,
			args{[]int{10}},
			want{idx, true}},
		{"fail: no levels",
			idx,
			args{[]int{}},
			want{idx, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.SubsetLevels(tt.args.pos)
			if !reflect.DeepEqual(tt.input, tt.want.index) {
				t.Errorf("Index.SubsetLevels(): got %v, want %v", tt.input, tt.want.index)
			}
		})
	}
}

func TestSet(t *testing.T) {
	idx := New(MustNewLevel([]string{"foo", "bar", "baz"}, ""), MustNewLevel([]string{"qux", "quux", "quuz"}, ""))
	type args struct {
		row   int
		level int
		val   interface{}
	}
	type want struct {
		index Index
		err   bool
	}
	tests := []struct {
		name  string
		input Index
		args  args
		want  want
	}{
		{"string", idx, args{0, 0, "corge"},
			want{New(MustNewLevel([]string{"corge", "bar", "baz"}, ""), MustNewLevel([]string{"qux", "quux", "quuz"}, "")), false}},
		{"null string", idx, args{0, 0, ""},
			want{New(MustNewLevel([]string{"NaN", "bar", "baz"}, ""), MustNewLevel([]string{"qux", "quux", "quuz"}, "")), false}},
		{"float", idx, args{2, 1, 1.5},
			want{New(MustNewLevel([]string{"foo", "bar", "baz"}, ""), MustNewLevel([]string{"qux", "quux", "1.5"}, "")), false}},
		{"fail: invalid row", idx, args{10, 0, "corge"},
			want{idx, true}},
		{"fail: invalid level", idx, args{0, 10, "corge"},
			want{idx, true}},
		{"fail: unsupported", idx, args{0, 0, complex64(1)},
			want{idx, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := tt.input.Copy()
			err := idx.Set(tt.args.row, tt.args.level, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("Set() error = %v, want %v", err, tt.want.err)
			}
			if !reflect.DeepEqual(idx, tt.want.index) {
				t.Errorf("Set(): got %v, want %v", idx, tt.want.index)
			}
		})
	}
}

func TestDropLevel(t *testing.T) {
	lvl := MustNewLevel([]string{"foo", "bar", "baz"}, "")
	single := New(lvl)
	multi := New(lvl, lvl)

	type args struct {
		pos int
	}
	type want struct {
		index Index
		err   bool
	}
	tests := []struct {
		name  string
		input Index
		args  args
		want  want
	}{
		{"one level: default index", single, args{0}, want{NewDefault(3), false}},
		{"two levels", multi, args{0}, want{single, false}},
		{"fail: invalid", single, args{10}, want{single, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.DropLevel(tt.args.pos)
			if (err != nil) != tt.want.err {
				t.Errorf("DropLevel() error = %v, want %v", err, tt.want.err)
			}
			if !reflect.DeepEqual(tt.input, tt.want.index) {
				t.Errorf("DropLevel(): got %v, want %v", tt.input, tt.want.index)
			}
		})
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

func Test_LevelCopy(t *testing.T) {
	idxLvl := MustNewLevel([]int{1, 2, 3}, "")
	idxLvl.Name = "foo"
	copyLvl := idxLvl.Copy()
	if !reflect.DeepEqual(idxLvl.LabelMap, copyLvl.LabelMap) {
		t.Error("Level.Copy() did not copy LabelMap")
	}
	if copyLvl.Name != "foo" {
		t.Error("Level.Copy() did not copy Name")
	}
	if copyLvl.DataType != options.Int64 {
		t.Error("Level.Copy() did not copy DataType")
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
		{MustNewLevel([]int64{3, 4}, "").Labels, LabelMap{"3": []int{0}, "4": []int{1}}, 1},
		{MustNewLevel([]int64{10, 20}, "").Labels, LabelMap{"10": []int{0}, "20": []int{1}}, 2},
	}
	for _, test := range tests {
		lvl := MustNewLevel([]int64{1, 2}, "")
		origLabelMap := LabelMap{"1": []int{0}, "2": []int{1}}
		if !reflect.DeepEqual(lvl.LabelMap, origLabelMap) {
			t.Errorf("Returned labelMap %v, want %v", lvl.LabelMap, origLabelMap)
		}

		lvl.Labels = test.newLabels
		lvl.Refresh()
		if !reflect.DeepEqual(lvl.LabelMap, test.wantLabelMap) {
			t.Errorf("Returned labelMap %v, want %v", lvl.LabelMap, test.wantLabelMap)
		}
		if lvl.maxWidth() != test.wantLongest {
			t.Errorf("Returned longest length %v, want %v", lvl.maxWidth(), test.wantLongest)
		}
	}
}

func TestConvertIndex_int(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		lvl       Level
		convertTo options.DataType
	}{
		// Float
		{MustNewLevel([]float64{1, 2, 3}, ""), options.Float64},
		{MustNewLevel([]float64{1, 2, 3}, ""), options.Int64},
		{MustNewLevel([]float64{1, 2, 3}, ""), options.String},
		{MustNewLevel([]float64{1, 2, 3}, ""), options.Bool},
		{MustNewLevel([]float64{1, 2, 3}, ""), options.DateTime},
		{MustNewLevel([]float64{1, 2, 3}, ""), options.Interface},

		// Int
		{MustNewLevel([]int64{1, 2, 3}, ""), options.Float64},
		{MustNewLevel([]int64{1, 2, 3}, ""), options.Int64},
		{MustNewLevel([]int64{1, 2, 3}, ""), options.String},
		{MustNewLevel([]int64{1, 2, 3}, ""), options.Bool},
		{MustNewLevel([]int64{1, 2, 3}, ""), options.DateTime},
		{MustNewLevel([]int64{1, 2, 3}, ""), options.Interface},

		// String
		{MustNewLevel([]string{"1", "2", "3"}, ""), options.Float64},
		{MustNewLevel([]string{"1", "2", "3"}, ""), options.Int64},
		{MustNewLevel([]string{"1", "2", "3"}, ""), options.String},
		{MustNewLevel([]string{"1", "2", "3"}, ""), options.Bool},
		{MustNewLevel([]string{"1", "2", "3"}, ""), options.DateTime},
		{MustNewLevel([]string{"1", "2", "3"}, ""), options.Interface},

		// Bool
		{MustNewLevel([]bool{true, false, false}, ""), options.Float64},
		{MustNewLevel([]bool{true, false, false}, ""), options.Int64},
		{MustNewLevel([]bool{true, false, false}, ""), options.String},
		{MustNewLevel([]bool{true, false, false}, ""), options.Bool},
		{MustNewLevel([]bool{true, false, false}, ""), options.DateTime},
		{MustNewLevel([]bool{true, false, false}, ""), options.Interface},

		// DateTime
		{MustNewLevel([]time.Time{testDate}, ""), options.Float64},
		{MustNewLevel([]time.Time{testDate}, ""), options.Int64},
		{MustNewLevel([]time.Time{testDate}, ""), options.String},
		{MustNewLevel([]time.Time{testDate}, ""), options.Bool},
		{MustNewLevel([]time.Time{testDate}, ""), options.DateTime},
		{MustNewLevel([]time.Time{testDate}, ""), options.Interface},

		// Interface
		{MustNewLevel([]interface{}{1, "2", true}, ""), options.Float64},
		{MustNewLevel([]interface{}{1, "2", true}, ""), options.Int64},
		{MustNewLevel([]interface{}{1, "2", true}, ""), options.String},
		{MustNewLevel([]interface{}{1, "2", true}, ""), options.Bool},
		{MustNewLevel([]interface{}{1, "2", true}, ""), options.DateTime},
		{MustNewLevel([]interface{}{1, "2", true}, ""), options.Interface},
	}
	for _, tt := range tests {
		err := tt.lvl.Convert(tt.convertTo)
		if err != nil {
			t.Error(err)
		}
		if tt.lvl.DataType != tt.convertTo {
			t.Errorf("index.Convert() = %v, want %v", tt.lvl, tt.convertTo)
		}
	}
}

func TestConvert_Numeric_Datetime(t *testing.T) {
	n := int64(1556668800000000000)
	wantVal := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		lvl Level
	}{
		{MustNewLevel([]int64{n}, "")},
		{MustNewLevel([]float64{float64(n)}, "")},
	}
	for _, test := range tests {
		_ = test.lvl.Convert(options.DateTime)
		elem := test.lvl.Labels.Element(0)
		gotVal := elem.Value.(time.Time)
		if gotVal != wantVal {
			t.Errorf("Error converting %v to datetime: returned %v, want %v", test.lvl, gotVal, wantVal)
		}
	}
}

func TestConvert_Unsupported(t *testing.T) {
	var tests = []struct {
		datatype options.DataType
	}{
		{options.None},
		{options.Unsupported},
	}
	for _, test := range tests {
		lvl := MustNewLevel([]float64{1, 2, 3}, "")
		err := lvl.Convert(test.datatype)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.datatype)
		}
	}
}
