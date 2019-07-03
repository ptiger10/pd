package series

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSeries_Join(t *testing.T) {
	single := MustNew("foo", Config{Index: []int{1}, IndexName: "foobar"})
	single2 := MustNew("bar", Config{Index: []int{2}, IndexName: "corge"})
	single3 := MustNew(7.11, Config{Index: []int{2}, IndexName: "corge"})
	multi := MustNew("foo", Config{MultiIndex: []interface{}{[]string{"A"}, []int{1}}, MultiIndexNames: []string{"foobar", "corge"}})
	multi2 := MustNew("bar", Config{MultiIndex: []interface{}{[]string{"B"}, []int{2}}, MultiIndexNames: []string{"waldo", "fred"}})
	type args struct {
		s2 *Series
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
		{name: "singleIndex",
			input: single, args: args{s2: single2},
			want: want{series: MustNew([]string{"foo", "bar"}, Config{Index: []int{1, 2}, IndexName: "foobar"}), err: false}},
		{"replace empty s",
			newEmptySeries(), args{s2: single2},
			want{MustNew([]string{"bar"}, Config{Index: []int{2}, IndexName: "corge"}), false}},
		{"singleIndex convert",
			single, args{single3},
			want{MustNew([]string{"foo", "7.11"}, Config{Index: []int{1, 2}, IndexName: "foobar"}), false}},
		{"multiIndex",
			multi, args{multi2},
			want{MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"A", "B"}, []int{1, 2}}, MultiIndexNames: []string{"foobar", "corge"}}), false}},
		{"fail: empty s2",
			single, args{newEmptySeries()},
			want{single, true}},
		{"fail: nil s2",
			single, args{&Series{}},
			want{single, true}},
		{"fail: invalid num levels",
			single, args{multi},
			want{single, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input.Copy()
			sArchive := tt.input.Copy()
			err := s.InPlace.Join(tt.args.s2)
			if (err != nil) != tt.want.err {
				t.Errorf("InPlace.Join() error = %v, want %v", err, tt.want.err)
				return
			}

			if !Equal(s, tt.want.series) {
				t.Errorf("InPlace.Join() got %v, want %v", s, tt.want.series)
			}

			sCopy, err := sArchive.Join(tt.args.s2)
			if (err != nil) != tt.want.err {
				t.Errorf("Series.Join() error = %v, want %v", err, tt.want.err)
				return
			}
			if !Equal(sCopy, tt.want.series) {
				t.Errorf("Series.Join() got %v, want %v", sCopy, tt.want.series)
			}
			if !strings.Contains(tt.name, "fail") {
				if !strings.Contains(tt.name, "same") {
					if Equal(sArchive, sCopy) {
						t.Errorf("Series.Join() retained access to original, want copy")
					}
				}
			}
		})
	}
}

func TestSeries_LookupSeries(t *testing.T) {
	multi := MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]string{"baz", "qux"}, []int{1, 2}}})
	multi2 := MustNew("corge", Config{MultiIndex: []interface{}{[]string{"baz"}, []int{1}}})
	type args struct {
		s2 *Series
	}
	tests := []struct {
		name     string
		input    *Series
		args     args
		want     *Series
		wantFail bool
	}{
		{name: "single", input: MustNew("foo"), args: args{s2: MustNew("bar")},
			want: MustNew("bar"), wantFail: false},
		{"multi", multi, args{multi2},
			MustNew([]string{"corge", ""}, Config{MultiIndex: []interface{}{[]string{"baz", "qux"}, []int{1, 2}}}), false},
		{"fail", MustNew("foo"), args{multi2},
			newEmptySeries(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			if got := tt.input.LookupSeries(tt.args.s2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Series.LookupSeries() = %v, want %v", got.index, tt.want.index)
			}
			if tt.wantFail {
				if buf.String() == "" {
					t.Errorf("Series.LookupSeries() returned no log message, want log due to fail")
				}
			}
		})
	}
}
