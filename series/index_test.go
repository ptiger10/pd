package series

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ptiger10/pd/internal/index"
)

func TestIndex_Less(t *testing.T) {
	s := MustNew([]int{1, 2, 3, 4}, Config{Index: []int{2, 0, 1, 1}})
	tests := []struct {
		i    int
		j    int
		want bool
	}{
		{0, 1, false},
		{0, 2, false},
		{1, 2, true},
		{2, 3, false},
	}
	for _, tt := range tests {
		idx := Index{s: s}
		got := idx.Less(tt.i, tt.j)
		if got != tt.want {
			t.Errorf("Index.Less() got %v, want %v", got, tt.want)
		}
	}
}

func TestIndex_Swap(t *testing.T) {
	s := MustNew([]int{1, 2}, Config{Index: []int{2, 0}})
	tests := []struct {
		i    int
		j    int
		want *Series
	}{
		{0, 1, MustNew([]int{2, 1}, Config{Index: []int{0, 2}})},
		{1, 0, MustNew([]int{2, 1}, Config{Index: []int{0, 2}})},
	}
	for _, tt := range tests {
		idx := Index{s: s.Copy()}
		idx.Swap(tt.i, tt.j)
		if !Equal(idx.s, tt.want) {
			t.Errorf("Index.Swap() got %v, want %v", idx.s, tt.want)
		}
	}
}

func TestIndex_Sort(t *testing.T) {
	var tests = []struct {
		name  string
		input *Series
		asc   bool
		want  *Series
	}{
		{"float", MustNew([]float64{1, 3, 5}, Config{Index: []int{2, 0, 1}}), true,
			MustNew([]float64{3, 5, 1}, Config{Index: []int{0, 1, 2}})},
		{"float reverse", MustNew([]float64{1, 3, 5}, Config{Index: []int{2, 0, 1}}), false,
			MustNew([]float64{1, 5, 3}, Config{Index: []int{2, 1, 0}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{s: tt.input}
			idx.Sort(tt.asc)
			if !Equal(idx.s, tt.want) {
				t.Errorf("Index.Sort() got %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_Describe(t *testing.T) {
	singleDefault := MustNew([]string{"foo", "bar", "baz"})
	multiConfig := MustNew([]string{"foo", "bar", "baz"}, Config{MultiIndex: []interface{}{[]int{1, 2, 3}, []string{"qux", "quux", "quuz"}}})
	type args struct {
		atRow   int
		atLevel int
	}
	type want struct {
		len int
		at  interface{}
	}
	tests := []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{name: "single default", input: singleDefault, args: args{atRow: 2, atLevel: 0},
			want: want{len: 3, at: int64(2)}},
		{"multi from config", multiConfig, args{2, 1}, want{3, "quuz"}},
		{"soft fail: at invalid row", singleDefault, args{10, 0}, want{3, nil}},
		{"soft fail: at invalid level", singleDefault, args{2, 10}, want{3, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			idx := Index{
				s: tt.input,
			}

			gotLen := idx.Len()
			wantLen := tt.want.len
			if gotLen != wantLen {
				t.Errorf("Index.Len(): got %v, want %v", gotLen, wantLen)
			}
			gotAt := idx.At(tt.args.atRow, tt.args.atLevel)
			wantAt := tt.want.at
			if gotAt != wantAt {
				t.Errorf("Index.At(): got %v, want %v", gotAt, wantAt)
			}
			if strings.Contains(tt.name, "fail:") {
				if buf.String() == "" {
					t.Errorf("Index operation returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestIndex_Values(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"}, Config{MultiIndex: []interface{}{[]float64{1, 2, 3}, []string{"qux", "quux", "quuz"}}})
	got := s.Index.Values()
	want := [][]interface{}{{1.0, 2.0, 3.0}, {"qux", "quux", "quuz"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Index.Values(): got %v, want %v", got, want)
	}
}

func TestIndex_Nil(t *testing.T) {
	s := newEmptySeries()
	s.index = index.Index{}
	s.Index.Len()
}

func TestIndex_Convert(t *testing.T) {
	type args struct {
		dataType string
		level    int
	}
	type want struct {
		s   *Series
		err bool
	}
	tests := []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{name: "pass", input: MustNew("foo", Config{Index: "bar"}),
			args: args{dataType: "bool", level: 0},
			want: want{s: MustNew("foo", Config{Index: true}), err: false}},
		{"fail: invalid column", MustNew("foo", Config{Index: "bar"}),
			args{dataType: "bool", level: 10},
			want{MustNew("foo", Config{Index: "bar"}), true}},
		{"fail: unsupported type", MustNew("foo", Config{Index: "bar"}),
			args{dataType: "corge", level: 0},
			want{MustNew("foo", Config{Index: "bar"}), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Index.Convert(tt.args.dataType, tt.args.level)
			if (err != nil) != tt.want.err {
				t.Errorf("s.Index.Convert() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(tt.input, tt.want.s) {
				t.Errorf("s.Index.Convert() got %v, want %v", tt.input, tt.want.s)
			}
		})
	}
}

func TestIndex_Flip(t *testing.T) {
	type args struct {
		level int
	}
	type want struct {
		series *Series
		err    bool
	}
	tests := []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{"pass", MustNew([]string{"foo", "bar", "baz"}, Config{Name: "corge", IndexName: "foobar", Index: []int64{0, 1, 2}}), args{0},
			want{MustNew([]int64{0, 1, 2}, Config{Name: "foobar", Index: []string{"foo", "bar", "baz"}, IndexName: "corge"}), false}},
		{"fail: invalid", MustNew("foo"), args{10},
			want{newEmptySeries(), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.input,
			}
			got, err := idx.Flip(tt.args.level)
			if (err != nil) != tt.want.err {
				t.Errorf("Index.Flip() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(got, tt.want.series) {
				t.Errorf("Index.Flip() = %v, \nwant %v", got, tt.want.series)
			}
		})
	}
}

func TestIndex_RenameLevel(t *testing.T) {
	type args struct {
		level int
		name  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Series
		wantErr bool
	}{
		{"pass 0", args{0, "corge"},
			MustNew("foo", Config{MultiIndex: []interface{}{"bar", "baz"}, MultiIndexNames: []string{"corge", "quuz"}}), false},
		{"pass 1", args{1, "corge"},
			MustNew("foo", Config{MultiIndex: []interface{}{"bar", "baz"}, MultiIndexNames: []string{"qux", "corge"}}), false},
		{"fail", args{10, "corge"},
			MustNew("foo", Config{MultiIndex: []interface{}{"bar", "baz"}, MultiIndexNames: []string{"qux", "quuz"}}), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: MustNew("foo", Config{MultiIndex: []interface{}{"bar", "baz"}, MultiIndexNames: []string{"qux", "quuz"}}),
			}
			if err := idx.RenameLevel(tt.args.level, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Index.RenameLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !Equal(idx.s, tt.want) {
				t.Errorf("Index.RenameLevel(): got %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_Reindex(t *testing.T) {
	s := MustNew("foo", Config{Index: "bar"})
	type args struct {
		level int
	}
	tests := []struct {
		name    string
		args    args
		want    *Series
		wantErr bool
	}{
		{name: "pass 0", args: args{0},
			want:    MustNew("foo"),
			wantErr: false},
		{"fail invalid level", args{10}, s, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: s.Copy(),
			}
			if err := idx.Reindex(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("Index.Reindex() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !Equal(idx.s, tt.want) {
				t.Errorf("Index.Reindex(): got %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_DropNull(t *testing.T) {

	type args struct {
		level int
	}
	tests := []struct {
		name    string
		input   *Series
		args    args
		want    *Series
		wantErr bool
	}{
		{name: "single", input: MustNew([]string{"foo", "bar"}, Config{Index: []string{"baz", ""}}), args: args{0},
			want: MustNew("foo", Config{Index: "baz"}), wantErr: false},
		{name: "multi", input: MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"qux", "quux"}, []string{"baz", ""}}}), args: args{1},
			want: MustNew("foo", Config{MultiIndex: []interface{}{"qux", "baz"}}), wantErr: false},
		{name: "fail: invalid", input: MustNew("foo"), args: args{10},
			want: MustNew("foo"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.input.Copy(),
			}
			err := idx.DropNull(tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index.DropNull() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !Equal(idx.s, tt.want) {
				t.Errorf("Index.DropNull(): got %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_SwapLevels(t *testing.T) {
	s := MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"baz", "qux"}, []string{"quux", "quuz"}}})
	type fields struct {
		s *Series
	}
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Series
		wantErr bool
	}{
		{name: "pass", fields: fields{s}, args: args{0, 1},
			want:    MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"quux", "quuz"}, []string{"baz", "qux"}}}),
			wantErr: false},
		{"reverse order", fields{s}, args{1, 0},
			MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"quux", "quuz"}, []string{"baz", "qux"}}}),
			false},
		{"fail: i", fields{s}, args{10, 1},
			s, true},
		{"fail: j", fields{s}, args{0, 10},
			s, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.fields.s.Copy(),
			}
			err := idx.SwapLevels(tt.args.i, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index.SwapLevels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(idx.s, tt.want) {
				t.Errorf("Index.SwapLevels() = %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_InsertLevel(t *testing.T) {
	s := MustNew([]string{"foo", "bar"}, Config{
		MultiIndex: []interface{}{[]string{"baz", "qux"}, []string{"quux", "quuz"}}})
	type fields struct {
		s *Series
	}
	type args struct {
		pos    int
		values interface{}
		name   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Series
		wantErr bool
	}{
		{name: "0", fields: fields{s}, args: args{0, []string{"corge", "fred"}, ""},
			want:    MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"corge", "fred"}, []string{"baz", "qux"}, []string{"quux", "quuz"}}}),
			wantErr: false},
		{"1", fields{s}, args{1, []string{"corge", "fred"}, ""},
			MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"baz", "qux"}, []string{"corge", "fred"}, []string{"quux", "quuz"}}}),
			false},
		{"2", fields{s}, args{2, []string{"corge", "fred"}, ""},
			MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"baz", "qux"}, []string{"quux", "quuz"}, []string{"corge", "fred"}}}),
			false},
		{"fail: invalid position", fields{s}, args{10, []string{"corge", "fred"}, ""},
			s, true},
		{"fail: unsupported value", fields{s}, args{2, []complex64{1, 2}, ""},
			s, true},
		{"fail: misaligned length", fields{s}, args{2, "corge", ""},
			s, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.fields.s.Copy(),
			}
			err := idx.InsertLevel(tt.args.pos, tt.args.values, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index.InsertLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(idx.s, tt.want) {
				t.Errorf("Index.InsertLevel() = %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_AppendLevel(t *testing.T) {
	s := MustNew([]string{"foo", "bar"}, Config{
		MultiIndex: []interface{}{[]string{"baz", "qux"}, []string{"quux", "quuz"}}})
	type fields struct {
		s *Series
	}
	type args struct {
		values interface{}
		name   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Series
		wantErr bool
	}{
		{name: "pass", fields: fields{s}, args: args{[]string{"corge", "fred"}, ""},
			want: MustNew([]string{"foo", "bar"}, Config{
				MultiIndex: []interface{}{[]string{"baz", "qux"}, []string{"quux", "quuz"}, []string{"corge", "fred"}}}),
			wantErr: false},
		{"fail: unsupported value", fields{s}, args{[]complex64{1, 2}, ""},
			s, true},
		{"fail: misaligned length", fields{s}, args{"corge", ""},
			s, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.fields.s.Copy(),
			}
			err := idx.AppendLevel(tt.args.values, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index.AppendLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(idx.s, tt.want) {
				t.Errorf("Index.AppendLevel() = %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_Set(t *testing.T) {
	type fields struct {
		s *Series
	}
	type args struct {
		row   int
		level int
		val   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Series
		wantErr bool
	}{
		{name: "0, 0", fields: fields{MustNew("foo", Config{Index: 0})}, args: args{0, 0, 100},
			want:    MustNew("foo", Config{Index: 100}),
			wantErr: false},
		{"fail: unsupported", fields{MustNew("foo")}, args{1, 0, complex64(1)},
			MustNew("foo"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.fields.s.Copy(),
			}
			err := idx.Set(tt.args.row, tt.args.level, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(idx.s, tt.want) {
				t.Errorf("Index.Set() = %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_DropLevel(t *testing.T) {
	s := MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"baz", "qux"}, []string{"quux", "quuz"}}})
	type fields struct {
		s *Series
	}
	type args struct {
		level int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Series
		wantErr bool
	}{
		{name: "pass", fields: fields{s}, args: args{0},
			want:    MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"quux", "quuz"}}}),
			wantErr: false},
		{"fail: invalid level", fields{s}, args{10},
			s, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.fields.s.Copy(),
			}
			err := idx.DropLevel(tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index.DropLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(idx.s, tt.want) {
				t.Errorf("Index.DropLevel() = %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_SelectName(t *testing.T) {
	s := MustNew("foo", Config{MultiIndex: []interface{}{"bar", "baz", "qux"}, MultiIndexNames: []string{"quux", "quuz", "quux"}})
	type fields struct {
		s *Series
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{name: "pass", fields: fields{s}, args: args{"quux"}, want: 0},
		{name: "soft fail: invalid name", fields: fields{s}, args: args{"fred"}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.fields.s,
			}
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)
			if got := idx.SelectName(tt.args.name); got != tt.want {
				t.Errorf("Index.SelectName() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail:") {
				if buf.String() == "" {
					t.Errorf("Index.SelectName() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestIndex_SelectNames(t *testing.T) {
	s := MustNew("foo", Config{MultiIndex: []interface{}{"bar", "baz", "qux"}, MultiIndexNames: []string{"quux", "quuz", "quux"}})
	type fields struct {
		s *Series
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{name: "pass", fields: fields{s}, args: args{[]string{"quux"}}, want: []int{0, 2}},
		{name: "soft fail: invalid name", fields: fields{s}, args: args{[]string{"fred"}}, want: []int{}},
		{name: "soft fail: partial invalid name", fields: fields{s}, args: args{[]string{"quux", "fred"}}, want: []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.fields.s,
			}
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)
			if got := idx.SelectNames(tt.args.names); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Index.SelectNames(): got %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail:") {
				if buf.String() == "" {
					t.Errorf("Index.SelectNames() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestIndex_SubsetLevels(t *testing.T) {
	s := MustNew("foo", Config{MultiIndex: []interface{}{"bar", "baz", "qux"}})
	type fields struct {
		s *Series
	}
	type args struct {
		levelPositions []int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Series
		wantErr bool
	}{
		{name: "one level", fields: fields{s}, args: args{[]int{0}},
			want: MustNew("foo", Config{MultiIndex: []interface{}{"bar"}}), wantErr: false},
		{name: "multiple levels", fields: fields{s}, args: args{[]int{0, 1}},
			want: MustNew("foo", Config{MultiIndex: []interface{}{"bar", "baz"}}), wantErr: false},
		{"fail: invalid level", fields{s}, args{[]int{10}},
			s, true},
		{"fail: no levels", fields{s}, args{[]int{}},
			s, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: tt.fields.s.Copy(),
			}
			err := idx.SubsetLevels(tt.args.levelPositions)
			if (err != nil) != tt.wantErr {
				t.Errorf("Index.Subset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(idx.s, tt.want) {
				t.Errorf("Index.Subset() = %v, want %v", idx.s, tt.want)
			}
		})
	}
}

func TestIndex_Filter(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"}, Config{Index: []string{"bamboo", "leaves", "taboo"}})
	fn := func(val interface{}) bool {
		v, ok := val.(string)
		if !ok {
			return false
		}
		if strings.HasSuffix(v, "boo") {
			return true
		}
		return false
	}
	type args struct {
		level int
		fn    func(interface{}) bool
	}
	tests := []struct {
		name  string
		input *Series
		args  args
		want  []int
	}{
		{name: "pass", input: s, args: args{level: 0, fn: fn}, want: []int{0, 2}},
		{"fail", s, args{10, fn}, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			got := s.Index.Filter(tt.args.level, tt.args.fn)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("s.Filter() got %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("s.Filter() returned no log message, want log due to fail")
				}
			}
		})
	}
}
