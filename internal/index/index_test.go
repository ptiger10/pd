package index

import (
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

func TestNew(t *testing.T) {
	empty, _ := values.InterfaceFactory(nil)
	labelsEmpty := empty.Values
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
		unnamed bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"empty", args{nil},
			want{
				index: Index{Levels: []Level{Level{Labels: labelsEmpty, LabelMap: LabelMap{}}}, NameMap: LabelMap{"": []int{0}}},
				len:   0, numLevels: 1, maxWidths: []int{0}, unnamed: true,
			}},
		{"one level",
			args{[]Level{MustNewLevel([]int{1, 2}, "foo")}},
			want{
				index: Index{
					Levels:  []Level{Level{Name: "foo", DataType: options.Int64, LabelMap: LabelMap{"1": []int{0}, "2": []int{1}}, Labels: labels}},
					NameMap: LabelMap{"foo": []int{0}}},
				len: 2, numLevels: 1, maxWidths: []int{3}, unnamed: false,
			}},
		{"two cols",
			args{[]Level{MustNewLevel([]int{1, 2}, "foo"), MustNewLevel([]int{1, 2}, "corge")}},
			want{
				index: Index{
					Levels: []Level{
						Level{Name: "foo", DataType: options.Int64, LabelMap: LabelMap{"1": []int{0}, "2": []int{1}}, Labels: labels},
						Level{Name: "corge", DataType: options.Int64, LabelMap: LabelMap{"1": []int{0}, "2": []int{1}}, Labels: labels}},
					NameMap: LabelMap{"foo": []int{0}, "corge": []int{1}}},
				len: 2, numLevels: 2, maxWidths: []int{3, 5}, unnamed: false,
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.levels...)
			if !reflect.DeepEqual(got, tt.want.index) {
				t.Errorf("New(): got %#v, want %#v", got.Levels[0], tt.want.index.Levels[0])
			}
			gotLen := got.Len()
			if gotLen != tt.want.len {
				t.Errorf("Index.Len(): got %v, want %v", gotLen, tt.want.len)
			}
			gotNumLevels := got.NumLevels()
			if !reflect.DeepEqual(gotNumLevels, tt.want.numLevels) {
				t.Errorf("Index.MaxWidth(): got %v, want %v", gotNumLevels, tt.want.numLevels)
			}
			gotMaxWidth := got.MaxWidths()
			if !reflect.DeepEqual(gotMaxWidth, tt.want.maxWidths) {
				t.Errorf("Index.MaxWidth(): got %v, want %v", gotMaxWidth, tt.want.maxWidths)
			}
			gotUnnamed := got.Unnamed()
			if !reflect.DeepEqual(gotUnnamed, tt.want.unnamed) {
				t.Errorf("Index.GotUnnamed(): got %v, want %v", gotUnnamed, tt.want.unnamed)
			}
		})
	}
}

func Test_NewDefault(t *testing.T) {
	got := NewDefault(3)
	lvl, err := NewLevel([]int64{0, 1, 2}, "")
	if err != nil {
		t.Error(err)
	}
	want := New(lvl)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Default constructor returned %v, want %v", got, want)
	}
	gotLen := len(got.Levels)
	wantLen := 1
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}
}

func Test_NewMulti(t *testing.T) {
	lvl1, err := NewLevel([]int64{0, 1, 2}, "")
	if err != nil {
		t.Error(err)
	}
	lvl2, err := NewLevel([]int64{100, 101, 102}, "")
	if err != nil {
		t.Error(err)
	}
	index := New(lvl1, lvl2)
	gotLen := len(index.Levels)
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}

}

func Test_Copy(t *testing.T) {
	idx := New(MustNewLevel([]int{1, 2, 3}, ""))
	copyIdx := idx.Copy()
	for i := 0; i < len(idx.Levels); i++ {
		if reflect.ValueOf(idx.Levels[i].Labels).Pointer() == reflect.ValueOf(copyIdx.Levels[i].Labels).Pointer() {
			t.Errorf("index.Copy() returned original labels at level %v, want fresh copy", i)
		}
		if reflect.ValueOf(idx.Levels[i].LabelMap).Pointer() == reflect.ValueOf(copyIdx.Levels[i].LabelMap).Pointer() {
			t.Errorf("index.Copy() returned original map at level %v, want fresh copy", i)
		}
	}
}

func Test_Drop_oneLevel(t *testing.T) {
	idx := New(MustNewLevel([]int{1, 2, 3}, ""))
	err := idx.Drop(0)
	if err != nil {
		t.Errorf("idx.Drop(): %v", err)
	}
	want := New(MustNewLevel([]int{1, 2, 3}, ""))
	if !reflect.DeepEqual(idx, want) {
		t.Errorf("idx.Drop() for one level returned %v, want %v", idx, want)
	}
}

func Test_Drop_multilevel(t *testing.T) {
	idx := New(MustNewLevel([]int{1, 2, 3}, ""), MustNewLevel([]int{4, 5, 6}, ""))
	idx.Drop(1)
	want := New(MustNewLevel([]int{1, 2, 3}, ""))
	if !reflect.DeepEqual(idx, want) {
		t.Errorf("idx.Drop() for multilevel returned %v, want %v", idx, want)
	}
}

func Test_Droplevels(t *testing.T) {
	idx := New(MustNewLevel([]int{1, 2, 3}, ""), MustNewLevel([]int{4, 5, 6}, ""), MustNewLevel([]int{7, 8, 9}, ""))
	err := idx.dropLevels([]int{2, 0})
	if err != nil {
		t.Errorf("idx.Droplevels(): %v", err)
	}
	want := New(MustNewLevel([]int{4, 5, 6}, ""))
	if !reflect.DeepEqual(idx, want) {
		t.Errorf("idx.Drop() for multilevel returned %v, want %v", idx, want)
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
			wantLabels: MustNewLevel([]float64{0, 1, 2}, "").Labels,
			wantKind:   options.Float64,
			wantName:   "test",
		},
		{
			data:       []int{0, 1, 2},
			wantLabels: MustNewLevel([]int64{0, 1, 2}, "").Labels,
			wantKind:   options.Int64,
			wantName:   "test",
		},
		{
			data:       []uint{0, 1, 2},
			wantLabels: MustNewLevel([]int64{0, 1, 2}, "").Labels,
			wantKind:   options.Int64,
			wantName:   "test",
		},
		{
			data:       []string{"0", "1", "2"},
			wantLabels: MustNewLevel([]string{"0", "1", "2"}, "").Labels,
			wantKind:   options.String,
			wantName:   "test",
		},
		{
			data:       []bool{true, true, false},
			wantLabels: MustNewLevel([]bool{true, true, false}, "").Labels,
			wantKind:   options.Bool,
			wantName:   "test",
		},
		{
			data:       []time.Time{testDate},
			wantLabels: MustNewLevel([]time.Time{testDate}, "").Labels,
			wantKind:   options.DateTime,
			wantName:   "test",
		},
		{
			data:       []interface{}{1.5, 1, "", false, testDate},
			wantLabels: MustNewLevel([]interface{}{1.5, 1, "", false, testDate}, "").Labels,
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
	for _, test := range tests {
		lvl, err := test.lvl.Convert(test.convertTo)
		if err != nil {
			t.Error(err)
		}
		if lvl.DataType != test.convertTo {
			t.Errorf("Attempted conversion to %v returned %v", test.convertTo, lvl.DataType)
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
		lvl, _ := test.lvl.Convert(options.DateTime)
		elem := lvl.Labels.Element(0)
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
		_, err := lvl.Convert(test.datatype)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.datatype)
		}
	}
}
