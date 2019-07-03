package series

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/d4l3k/messagediff"
	"github.com/ptiger10/pd/options"
)

// Modify tests check both inplace and copy functionality in the same test, if both are available
func TestRename(t *testing.T) {
	s := MustNew("foo", Config{Name: "baz"})
	want := "qux"
	s.Rename(want)
	got := s.Name()
	if got != want {
		t.Errorf("Rename() got %v, want %v", got, want)
	}
}

func TestModify_Sort(t *testing.T) {
	testDate1 := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	testDate2 := testDate1.Add(24 * time.Hour)
	testDate3 := testDate2.Add(24 * time.Hour)

	type args struct {
		asc bool
	}
	type want struct {
		series *Series
	}
	var tests = []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{"float",
			MustNew([]float64{3, 5, 1}, Config{Index: []int{0, 1, 2}}), args{true},
			want{MustNew([]float64{1, 3, 5}, Config{Index: []int{2, 0, 1}})}},
		{"float desc",
			MustNew([]float64{3, 5, 1}, Config{Index: []int{0, 1, 2}}), args{false},
			want{MustNew([]float64{5, 3, 1}, Config{Index: []int{1, 0, 2}})}},

		{"int",
			MustNew([]int{3, 5, 1}, Config{Index: []int{0, 1, 2}}), args{true},
			want{MustNew([]int{1, 3, 5}, Config{Index: []int{2, 0, 1}})}},
		{"int desc",
			MustNew([]int{3, 5, 1}, Config{Index: []int{0, 1, 2}}), args{false},
			want{MustNew([]int{5, 3, 1}, Config{Index: []int{1, 0, 2}})}},

		{"string",
			MustNew([]string{"15", "3", "1"}, Config{Index: []int{0, 1, 2}}), args{true},
			want{MustNew([]string{"1", "15", "3"}, Config{Index: []int{2, 0, 1}})}},
		{"string desc",
			MustNew([]string{"15", "3", "1"}, Config{Index: []int{0, 1, 2}}), args{false},
			want{MustNew([]string{"3", "15", "1"}, Config{Index: []int{1, 0, 2}})}},

		{"bool",
			MustNew([]bool{false, true, false}, Config{Index: []int{0, 1, 2}}), args{true},
			want{MustNew([]bool{false, false, true}, Config{Index: []int{0, 2, 1}})}},
		{"bool desc",
			MustNew([]bool{false, true, false}, Config{Index: []int{0, 1, 2}}), args{false},
			want{MustNew([]bool{true, false, false}, Config{Index: []int{1, 0, 2}})}},

		{"datetime",
			MustNew([]time.Time{testDate2, testDate3, testDate1}, Config{Index: []int{0, 1, 2}}), args{true},
			want{MustNew([]time.Time{testDate1, testDate2, testDate3}, Config{Index: []int{2, 0, 1}})}},
		{"datetime desc",
			MustNew([]time.Time{testDate2, testDate3, testDate1}, Config{Index: []int{0, 1, 2}}), args{false},
			want{MustNew([]time.Time{testDate3, testDate2, testDate1}, Config{Index: []int{1, 0, 2}})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			sArchive := tt.input.Copy()
			s.InPlace.Sort(tt.args.asc)
			if !Equal(s, tt.want.series) {
				t.Errorf("InPlace.Sort() got %v, want %v", s, tt.want.series)
			}

			sCopy := sArchive.Sort(tt.args.asc)
			if !Equal(sCopy, tt.want.series) {
				t.Errorf("Series.Sort() got %v, want %v", sCopy, tt.want.series)
			}
			if Equal(sArchive, sCopy) {
				t.Errorf("Series.Sort() retained access to original, want copy")
			}
		})
	}
}

func TestModify_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	var tests = []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"0,1", args{0, 1}, false},
		{"1,0", args{1, 0}, false},
		{"fail i", args{2, 0}, true},
		{"fail j", args{0, 2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MustNew([]string{"foo", "bar"}, Config{Index: []int{0, 1}})
			sArchive := s.Copy()
			want := MustNew([]string{"bar", "foo"}, Config{Index: []int{1, 0}})

			sCopy, err := sArchive.Swap(tt.args.i, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("Series.Swap() error = %v, want %v", err, tt.wantErr)
				return
			}

			// intentionally skip fail case
			if !strings.Contains(tt.name, "fail") {
				s.InPlace.Swap(tt.args.i, tt.args.j)
				if !Equal(s, want) {
					t.Errorf("InPlace.Swap() got %v, want %v", s, want)
				}
				if !Equal(sCopy, want) {
					t.Errorf("Series.Swap() got %v, want %v", sCopy, want)
				}
				if Equal(sArchive, sCopy) {
					t.Errorf("Series.Swap() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModify_Insert(t *testing.T) {
	multi := MustNew([]string{"foo"}, Config{MultiIndex: []interface{}{"baz", 1}})

	type args struct {
		pos int
		val interface{}
		idx []interface{}
	}
	type want struct {
		series *Series
		err    bool
	}
	var tests = []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{name: "emptySeries",
			input: newEmptySeries(),
			args:  args{pos: 0, val: "foo", idx: []interface{}{"bar"}},
			want:  want{series: MustNew("foo", Config{Index: "bar"}), err: false}},
		{"singleIndex",
			MustNew([]string{"foo"}, Config{Index: "baz"}),
			args{0, "bar", []interface{}{"qux"}},
			want{series: MustNew([]string{"bar", "foo"}, Config{Index: []string{"qux", "baz"}}), err: false}},
		{"default index, empty arg",
			MustNew([]string{"foo", "bar"}),
			args{1, "baz", nil},
			want{series: MustNew([]string{"foo", "baz", "bar"}), err: false}},
		{"multiIndex",
			multi,
			args{1, "bar", []interface{}{"qux", 2}},
			want{series: MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"baz", "qux"}, []int{1, 2}}}), err: false}},
		{"fail: index too long",
			multi,
			args{1, "bar", []interface{}{"foo", "baz", "qux"}},
			want{nil, true}},
		{"fail: invalid position",
			multi,
			args{3, "bar", []interface{}{"C", 3}},
			want{nil, true}},
		{"fail: unsupported index value",
			MustNew([]string{"foo"}, Config{Index: "A"}),
			args{0, "bar", []interface{}{complex64(1)}},
			want{nil, true}},
		{"fail: unsupported series value",
			MustNew([]string{"foo"}, Config{Index: "A"}),
			args{0, complex64(1), []interface{}{"A"}},
			want{nil, true}},
		{"fail: unsupported value inserting into empty series",
			newEmptySeries(),
			args{0, complex64(1), []interface{}{"A"}},
			want{nil, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			sArchive := tt.input.Copy()
			err := s.InPlace.Insert(tt.args.pos, tt.args.val, tt.args.idx...)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.Insert() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(s, tt.want.series) {
					t.Errorf("InPlace.Insert() got %v, want %v", s, tt.want.series)
				}
			}

			sCopy, err := sArchive.Insert(tt.args.pos, tt.args.val, tt.args.idx...)
			if (err != nil) != tt.want.err {
				t.Errorf("Series.Insert() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(sCopy, tt.want.series) {
					t.Errorf("Series.Insert() got %v, want %v", sCopy, tt.want.series)
					diff, _ := messagediff.PrettyDiff(s, tt.want.series)
					fmt.Println(diff)
				}
				if Equal(sArchive, sCopy) {
					t.Errorf("Series.Insert() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModify_Append(t *testing.T) {
	type args struct {
		val interface{}
		idx []interface{}
	}
	type want struct {
		series *Series
		err    bool
	}
	var tests = []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{"singleIndex",
			MustNew([]string{"foo"}, Config{Index: []int{1}}),
			args{val: "bar", idx: []interface{}{2}},
			want{series: MustNew([]string{"foo", "bar"}, Config{Index: []int{1, 2}}), err: false}},
		{"multiIndex",
			MustNew([]string{"foo"}, Config{MultiIndex: []interface{}{[]string{"A"}, []int{1}}}),
			args{"bar", []interface{}{"B", 2}},
			want{MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"A", "B"}, []int{1, 2}}}), false}},
		{"singleIndex: default index and nil index values",
			MustNew("foo"),
			args{"bar", nil},
			want{MustNew([]string{"foo", "bar"}), false}},
		{"fail multiIndex: too many index values",
			MustNew([]string{"foo"}),
			args{"bar", []interface{}{"qux", "baz"}},
			want{nil, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			sArchive := tt.input.Copy()
			err := s.InPlace.Append(tt.args.val, tt.args.idx...)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.Append() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(s, tt.want.series) {
					t.Errorf("InPlace.Append() got %v, want %v", s, tt.want.series)
				}
			}
			sCopy, err := sArchive.Append(tt.args.val, tt.args.idx...)
			if (err != nil) != tt.want.err {
				t.Errorf("Series.Append() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(sCopy, tt.want.series) {
					t.Errorf("Series.Append() got %v, want %v", sCopy, tt.want.series)
				}
				if Equal(sArchive, sCopy) {
					t.Errorf("Series.Append() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModify_Set(t *testing.T) {
	type args struct {
		rowPositions int
		val          interface{}
	}
	type want struct {
		series *Series
		err    bool
	}
	var tests = []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{"singleRow",
			MustNew("foo"), args{rowPositions: 0, val: "bar"},
			want{series: MustNew("bar"), err: false}},
		{"fail: invalid index singleRow",
			MustNew("foo"), args{1, "bar"},
			want{MustNew("foo"), true}},
		{"fail: unsupported value",
			MustNew("foo"), args{0, complex64(1)},
			want{MustNew("foo"), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			sArchive := tt.input.Copy()
			err := s.InPlace.Set(tt.args.rowPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(s, tt.want.series) {
				t.Errorf("InPlace.Set() got %v, want %v", s, tt.want.series)
			}

			sCopy, err := sArchive.Set(tt.args.rowPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("Series.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(sCopy, tt.want.series) {
					t.Errorf("Series.Set() got %v, want %v", sCopy, tt.want.series)
				}
				if Equal(sArchive, sCopy) {
					t.Errorf("Series.Set() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModify_SetRows(t *testing.T) {
	type args struct {
		rowPositions []int
		val          interface{}
	}
	type want struct {
		series *Series
		err    bool
	}
	var tests = []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{"singleRow",
			MustNew("foo"), args{rowPositions: []int{0}, val: "bar"},
			want{series: MustNew("bar"), err: false}},
		{"multiRow",
			MustNew([]string{"foo", "bar"}), args{[]int{0, 1}, "baz"},
			want{MustNew([]string{"baz", "baz"}), false}},
		{"fail: singleRow",
			MustNew("foo"), args{rowPositions: []int{0}, val: complex64(1)},
			want{series: MustNew("foo"), err: true}},
		{"fail: invalid index singleRow",
			MustNew("foo"), args{[]int{1}, "bar"},
			want{MustNew("foo"), true}},
		{"fail: partial success on multiRow",
			MustNew("foo"), args{[]int{0, 2}, "bar"},
			want{MustNew("foo"), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			sArchive := tt.input.Copy()
			err := s.InPlace.SetRows(tt.args.rowPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(s, tt.want.series) {
				t.Errorf("InPlace.Set() got %v, want %v", s, tt.want.series)
			}

			sCopy, err := sArchive.SetRows(tt.args.rowPositions, tt.args.val)
			if (err != nil) != tt.want.err {
				t.Errorf("Series.Set() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(sCopy, tt.want.series) {
					t.Errorf("Series.Set() got %v, want %v", sCopy, tt.want.series)
				}
				if Equal(sArchive, sCopy) {
					t.Errorf("Series.Set() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModify_Drop(t *testing.T) {
	type args struct {
		rowPositions int
	}
	type want struct {
		series *Series
		err    bool
	}
	var tests = []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{"drop to 0",
			MustNew("foo", Config{Index: []int{0}}), args{rowPositions: 0},
			want{series: newEmptySeries(), err: false}},
		{"singleRow",
			MustNew([]string{"foo", "bar", "baz"}, Config{Index: []int{0, 1, 2}}), args{1},
			want{MustNew([]string{"foo", "baz"}, Config{Index: []int{0, 2}}), false}},
		{"fail: invalid index singleRow",
			MustNew("foo"), args{1},
			want{MustNew("foo"), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			sArchive := tt.input.Copy()
			err := s.InPlace.Drop(tt.args.rowPositions)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.Drop() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(s, tt.want.series) {
				t.Errorf("InPlace.Drop() got %v, want %v", s, tt.want.series)
			}

			sCopy, err := sArchive.Drop(tt.args.rowPositions)
			if (err != nil) != tt.want.err {
				t.Errorf("Series.Drop() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(sCopy, tt.want.series) {
					t.Errorf("Series.Drop() got %v, want %v", sCopy, tt.want.series)
				}
				if Equal(sArchive, sCopy) {
					t.Errorf("Series.Drop() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModify_DropRows(t *testing.T) {
	type args struct {
		rowPositions []int
	}
	type want struct {
		series *Series
		err    bool
	}
	var tests = []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{"drop to 0",
			MustNew("foo"), args{rowPositions: []int{0}},
			want{series: newEmptySeries(), err: false}},
		{"singleRow",
			MustNew([]string{"foo", "bar", "baz"}, Config{Index: []int{0, 1, 2}}), args{[]int{1}},
			want{MustNew([]string{"foo", "baz"}, Config{Index: []int{0, 2}}), false}},
		{"multiRow",
			MustNew([]string{"foo", "bar", "baz"}, Config{Index: []int{0, 1, 2}}), args{[]int{1, 2}},
			want{MustNew([]string{"foo"}, Config{Index: []int{0}}), false}},
		{"multiRow reverse",
			MustNew([]string{"foo", "bar", "baz"}, Config{Index: []int{0, 1, 2}}), args{[]int{2, 1}},
			want{MustNew([]string{"foo"}, Config{Index: []int{0}}), false}},
		{"fail: invalid index singleRow",
			MustNew("foo"), args{[]int{1}},
			want{MustNew("foo"), true}},
		{"fail: partial success on multiRow",
			MustNew([]string{"foo", "bar"}), args{[]int{0, 2}},
			want{MustNew([]string{"foo", "bar"}), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			sArchive := tt.input.Copy()
			err := s.InPlace.DropRows(tt.args.rowPositions)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.DropRows() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(s, tt.want.series) {
				t.Errorf("InPlace.DropRows() got %#v, want %#v", s, tt.want.series)
			}

			sCopy, err := sArchive.DropRows(tt.args.rowPositions)
			if (err != nil) != tt.want.err {
				t.Errorf("Series.DropRows() error = %v, want %v", err, tt.want.err)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(sCopy, tt.want.series) {
					t.Errorf("Series.DropRows() got %v, want %v", sCopy, tt.want.series)
				}
				if Equal(sArchive, sCopy) {
					t.Errorf("Series.DropRows() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModify_DropDuplicates(t *testing.T) {
	var tests = []struct {
		name  string
		input *Series
		want  *Series
	}{
		{"single",
			MustNew([]string{"foo", "foo", "bar"}, Config{
				Index: []int{0, 0, 1}}),
			MustNew([]string{"foo", "bar"}, Config{Index: []int{0, 1}})},
		{"multi",
			MustNew([]string{"foo", "foo", "bar"}, Config{
				MultiIndex: []interface{}{[]int{0, 0, 1}, []int{2, 2, 3}}}),
			MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]int{0, 1}, []int{2, 3}}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			sArchive := tt.input.Copy()
			s.InPlace.DropDuplicates()
			if !Equal(s, tt.want) {
				t.Errorf("InPlace.DropDuplicates() got %v, want %v", s, tt.want)
			}

			sCopy := sArchive.DropDuplicates()
			if !strings.Contains(tt.name, "fail") {
				if !Equal(sCopy, tt.want) {
					t.Errorf("Series.DropDuplicates() got %v, want %v", sCopy, tt.want)
				}
				if Equal(sArchive, sCopy) {
					t.Errorf("Series.DropDuplicates() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModify_DropNull(t *testing.T) {
	type want struct {
		series *Series
	}
	var tests = []struct {
		name  string
		input *Series
		want  want
	}{
		{"control: no null rows",
			MustNew("foo"),
			want{series: MustNew("foo")}},
		{"null row",
			MustNew([]string{"foo", "", "baz"}, Config{Index: []int{0, 1, 2}}),
			want{MustNew([]string{"foo", "baz"}, Config{Index: []int{0, 2}})}},
		{"null row reverse",
			MustNew([]string{"baz", "", "foo"}, Config{Index: []int{0, 1, 2}}),
			want{MustNew([]string{"baz", "foo"}, Config{Index: []int{0, 2}})}},
		{"all null rows",
			MustNew([]string{"", ""}),
			want{newEmptySeries()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			sArchive := tt.input.Copy()
			s.InPlace.DropNull()
			if !Equal(s, tt.want.series) {
				t.Errorf("InPlace.DropNull() got %v, want %v", s, tt.want.series)
			}

			sCopy := sArchive.DropNull()
			if !Equal(sCopy, tt.want.series) {
				t.Errorf("Series.DropNull() got %v, want %v", sCopy, tt.want.series)
			}
			if !strings.Contains(tt.name, "control") {
				if Equal(sArchive, sCopy) {
					t.Errorf("Series.DropNull() retained access to original, want copy")
				}
			}
		})
	}
}

func TestModifyInPlace_DatatypeConversion(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	type args struct {
		To (func(InPlace))
	}
	type want struct {
		series   *Series
		datatype options.DataType
	}
	var tests = []struct {
		name string
		args args
		want want
	}{
		{"float", args{(InPlace).ToFloat64}, want{MustNew([]float64{1.5, 1.0, 1.0, 0, 1.5566688e+18}), options.Float64}},
		{"int", args{(InPlace).ToInt64}, want{MustNew([]int64{1, 1, 1, 0, 1.5566688e+18}), options.Int64}},
		{"string", args{(InPlace).ToString}, want{MustNew([]string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}), options.String}},
		{"bool", args{(InPlace).ToBool}, want{MustNew([]bool{true, true, true, false, true}), options.Bool}},
		{"datetime", args{(InPlace).ToDateTime}, want{MustNew([]time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate}), options.DateTime}},
		{"control", args{(InPlace).ToInterface}, want{MustNew([]interface{}{1.5, 1, "1", false, testDate}), options.Interface}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MustNew([]interface{}{1.5, 1, "1", false, testDate})
			tt.args.To(s.InPlace)
			if !Equal(s, tt.want.series) {
				t.Errorf("InPlace.To... got %v, want %v", s, tt.want.series)
			}
			if s.datatype != tt.want.datatype {
				t.Errorf("InPlace.To... got datatype %v, want %v", s.datatype, tt.want.datatype)
			}
		})
	}
}

func TestModify_DatatypeConversion(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	type args struct {
		To (func(*Series) *Series)
	}
	type want struct {
		series   *Series
		datatype options.DataType
	}
	var tests = []struct {
		name string
		args args
		want want
	}{
		{"float", args{(*Series).ToFloat64}, want{MustNew([]float64{1.5, 1.0, 1.0, 0, 1.5566688e+18}), options.Float64}},
		{"int", args{(*Series).ToInt64}, want{MustNew([]int64{1, 1, 1, 0, 1.5566688e+18}), options.Int64}},
		{"string", args{(*Series).ToString}, want{MustNew([]string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}), options.String}},
		{"bool", args{(*Series).ToBool}, want{MustNew([]bool{true, true, true, false, true}), options.Bool}},
		{"datetime", args{(*Series).ToDateTime}, want{MustNew([]time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate}), options.DateTime}},
		{"control: interface", args{(*Series).ToInterface}, want{MustNew([]interface{}{1.5, 1, "1", false, testDate}), options.Interface}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MustNew([]interface{}{1.5, 1, "1", false, testDate})
			got := tt.args.To(s)
			if !Equal(got, tt.want.series) {
				t.Errorf("Series.To... got %v, want %v", got, tt.want.series)
			}
			if got.datatype != tt.want.datatype {
				t.Errorf("Series.To... got datatype %v, want %v", got.datatype, tt.want.datatype)
			}
			if !strings.Contains(tt.name, "control") {
				if s.DataType() == got.DataType() {
					t.Errorf("Series.To... retained access to original, want copy")
				}
			}
		})
	}
}
