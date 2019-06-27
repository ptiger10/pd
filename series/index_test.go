package series

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ptiger10/pd/internal/index"
)

// func TestIndex_Sort(t *testing.T) {
// 	var tests = []struct {
// 		desc  string
// 		input *Series
// 		asc   bool
// 		want  *Series
// 	}{
// 		{"float", MustNew([]float64{1, 3, 5}, Config{Index: []int{2, 0, 1}}), true,
// 			MustNew([]float64{3, 5, 1}, Config{Index: []int{0, 1, 2}})},
// 		{"float reverse", MustNew([]float64{1, 3, 5}, Config{Index: []int{2, 0, 1}}), false,
// 			MustNew([]float64{1, 5, 3}, Config{Index: []int{2, 1, 0}})},
// 	}
// 	for _, tt := range tests {
// 		s := tt.input
// 		s.Index.Sort(tt.asc)
// 		if !Equal(s, tt.want) {
// 			t.Errorf("Index.Sort() test %v got %v, want %v", tt.desc, s, tt.want)
// 		}
// 	}
// }

func TestIndex_Describe(t *testing.T) {
	singleDefault := MustNew([]string{"foo", "bar", "baz"})
	multiConfig := MustNew([]string{"foo", "bar", "baz"}, Config{MultiIndex: []interface{}{[]int{1, 2, 3}, []string{"qux", "quux", "quuz"}}})
	type args struct {
		atRow     int
		atLevel   int
		valsLevel int
	}
	type want struct {
		len       int
		numLevels int
		at        interface{}
		vals      interface{}
	}
	tests := []struct {
		name  string
		input *Series
		args  args
		want  want
	}{
		{"single default", singleDefault, args{2, 0, 0}, want{3, 1, int64(2), []int64{0, 1, 2}}},
		{"multi from config", multiConfig, args{2, 1, 1}, want{3, 2, "quuz", []string{"qux", "quux", "quuz"}}},
		{"fail: at invalid row", singleDefault, args{10, 0, 0}, want{3, 1, nil, []int64{0, 1, 2}}},
		{"fail: at invalid level", singleDefault, args{2, 10, 0}, want{3, 1, nil, []int64{0, 1, 2}}},
		{"fail: values invalid level", singleDefault, args{2, 0, 10}, want{3, 1, int64(2), nil}},
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
			gotNumLevels := idx.NumLevels()
			wantNumLevels := tt.want.numLevels
			if gotNumLevels != wantNumLevels {
				t.Errorf("Index.NumLevels(): got %v, want %v", gotNumLevels, wantNumLevels)
			}
			gotAt := idx.At(tt.args.atRow, tt.args.atLevel)
			wantAt := tt.want.at
			if gotAt != wantAt {
				t.Errorf("Index.At(): got %v, want %v", gotAt, wantAt)
			}
			gotVals := idx.Values(tt.args.valsLevel)
			wantVals := tt.want.vals
			if !reflect.DeepEqual(gotVals, wantVals) {
				t.Errorf("Index.Values(): got %v, want %v", gotVals, wantVals)
			}
			if strings.Contains(tt.name, "fail:") {
				if buf.String() == "" {
					t.Errorf("Index operation returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestIndex_Nil(t *testing.T) {
	s := newEmptySeries()
	s.index = index.Index{}
	s.Index.Len()
}

func TestConversions(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	type args struct {
		fn    func(Index, int) (*Series, error)
		level int
	}
	type want struct {
		series *Series
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"toFloat", args{Index.LevelToFloat64, 0},
			want{MustNew([]string{"a", "b", "c", "d", "e"}, Config{Index: []float64{1.5, 1.0, 1.0, 0, 1.5566688e+18}}), false}},
		{"fail: toFloat", args{Index.LevelToFloat64, 10}, want{newEmptySeries(), true}},
		{"toInt", args{Index.LevelToInt64, 0},
			want{MustNew([]string{"a", "b", "c", "d", "e"}, Config{Index: []int64{1, 1, 1, 0, 1.5566688e+18}}), false}},
		{"fail: toInt", args{Index.LevelToInt64, 10}, want{newEmptySeries(), true}},
		{"toString", args{Index.LevelToString, 0},
			want{MustNew([]string{"a", "b", "c", "d", "e"}, Config{Index: []string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}}), false}},
		{"fail: toString", args{Index.LevelToString, 10}, want{newEmptySeries(), true}},
		{"toBool", args{Index.LevelToBool, 0},
			want{MustNew([]string{"a", "b", "c", "d", "e"}, Config{Index: []bool{true, true, true, false, true}}), false}},
		{"fail: toBool", args{Index.LevelToBool, 10}, want{newEmptySeries(), true}},
		{"toDateTime", args{Index.LevelToDateTime, 0},
			want{MustNew([]string{"a", "b", "c", "d", "e"}, Config{Index: []time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate}}), false}},
		{"fail: toDateTime", args{Index.LevelToDateTime, 10}, want{newEmptySeries(), true}},
		{"toInterface", args{Index.LevelToInterface, 0},
			want{MustNew([]string{"a", "b", "c", "d", "e"}, Config{Index: []interface{}{1.5, 1, "1", false, testDate}}), false}},
		{"fail: toInterface", args{Index.LevelToInterface, 10}, want{newEmptySeries(), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Index{
				s: MustNew([]string{"a", "b", "c", "d", "e"}, Config{Index: []interface{}{1.5, 1, "1", false, testDate}}),
			}
			got, err := tt.args.fn(idx, tt.args.level)
			if (err != nil) != tt.want.err {
				t.Errorf("Index conversion error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(got, tt.want.series) {
				t.Errorf("Index conversion = %v, \nwant %v", got, tt.want.series)
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
		{"pass", MustNew([]string{"foo", "bar", "baz"}, Config{Name: "corge", IndexName: "foobar"}), args{0},
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
