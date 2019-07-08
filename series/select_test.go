package series

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ptiger10/pd/options"
)

func TestElement(t *testing.T) {
	s, err := New([]string{"", "valid"}, Config{MultiIndex: []interface{}{[]string{"A", "B"}, []int{1, 2}}})
	if err != nil {
		t.Error(err)
	}
	var tests = []struct {
		position int
		wantVal  interface{}
		wantNull bool
		wantIdx  []interface{}
	}{
		{0, "NaN", true, []interface{}{"A", int64(1)}},
		{1, "valid", false, []interface{}{"B", int64(2)}},
	}
	wantIdxTypes := []options.DataType{options.String, options.Int64}
	for _, test := range tests {
		got := s.Element(test.position)
		if got.Value != test.wantVal {
			t.Errorf("Element returned value %v, want %v", got.Value, test.wantVal)
		}
		if got.Null != test.wantNull {
			t.Errorf("Element returned bool %v, want %v", got.Null, test.wantNull)
		}
		if !reflect.DeepEqual(got.Labels, test.wantIdx) {
			t.Errorf("Element returned index %#v, want %#v", got.Labels, test.wantIdx)
		}
		if !reflect.DeepEqual(got.LabelTypes, wantIdxTypes) {
			t.Errorf("Element returned kind %v, want %v", got.LabelTypes, wantIdxTypes)
		}
	}
}

func TestSeries_At(t *testing.T) {
	type args struct {
		position int
	}
	tests := []struct {
		name     string
		input    *Series
		args     args
		want     interface{}
		wantFail bool
	}{
		{name: "pass", input: MustNew([]string{"foo", "bar", "baz"}), args: args{1}, want: "bar", wantFail: false},
		{"soft fail: invalid position", MustNew([]string{"foo", "bar", "baz"}), args{10}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			if got := tt.input.At(tt.args.position); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Series.At() = %v, want %v", got, tt.want)
			}

			if tt.wantFail {
				if buf.String() == "" {
					t.Errorf("Series.At() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestFrom(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"}, Config{Index: []int{0, 1, 2}})
	type args struct {
		start int
		end   int
	}
	tests := []struct {
		name  string
		input *Series
		args  args
		want  *Series
	}{
		{name: "ascending", input: s, args: args{start: 0, end: 2}, want: MustNew([]string{"foo", "bar", "baz"}, Config{Index: []int{0, 1, 2}})},
		{"single", s, args{1, 1}, MustNew([]string{"bar"}, Config{Index: []int{1}})},
		{"partial", s, args{1, 2}, MustNew([]string{"bar", "baz"}, Config{Index: []int{1, 2}})},
		{"descending", s, args{2, 0}, MustNew([]string{"baz", "bar", "foo"}, Config{Index: []int{2, 1, 0}})},
		{"fail: partial invalid", s, args{10, 0}, newEmptySeries()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			got := tt.input.From(tt.args.start, tt.args.end)
			if !Equal(got, tt.want) {
				t.Errorf("Series.From() got %v, want %v", s, tt.want)
			}

			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("Series.From() returned no log message, want log due to fail")
				}
			}
		})

	}
}

func TestSeries_XS(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"}, Config{MultiIndex: []interface{}{[]int{1, 2, 3}, []string{"qux", "quux", "quuz"}}})
	type args struct {
		rowPositions   []int
		levelPositions []int
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
		{name: "pass", input: s, args: args{[]int{0, 1}, []int{0}},
			want: want{series: MustNew([]string{"foo", "bar"}, Config{MultiIndex: []interface{}{[]int{1, 2}}}), err: false}},
		{"pass reverse", s, args{[]int{1, 0}, []int{1}},
			want{MustNew([]string{"bar", "foo"}, Config{MultiIndex: []interface{}{[]string{"quux", "qux"}}}), false}},
		{"pass multi reverse", s, args{[]int{1, 0}, []int{1, 0}},
			want{MustNew([]string{"bar", "foo"}, Config{MultiIndex: []interface{}{[]string{"quux", "qux"}, []int{2, 1}}}), false}},
		{"fail: invalid row position", s, args{[]int{10}, []int{0}}, want{newEmptySeries(), true}},
		{"fail: partial invalid row position", s, args{[]int{0, 10}, []int{0}}, want{newEmptySeries(), true}},
		{"fail: invalid level position", s, args{[]int{0}, []int{10}}, want{newEmptySeries(), true}},
		{"fail: partial invalid level position", s, args{[]int{0}, []int{0, 10}}, want{newEmptySeries(), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.XS(tt.args.rowPositions, tt.args.levelPositions)
			if (err != nil) != tt.want.err {
				t.Errorf("Series.XS() error = %v, want %v", err, tt.want.err)
			}
			if !reflect.DeepEqual(got, tt.want.series) {
				t.Errorf("Series.XS() = %v, want %v", got, tt.want.series)
			}

		})
	}
}

func TestSeries_SelectLabel(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"}, Config{Index: []int{1, 1, 1}})
	type args struct {
		label string
	}
	tests := []struct {
		name  string
		input *Series
		args  args
		want  int
	}{
		{name: "pass", input: s, args: args{label: "1"}, want: 0},
		{"fail: empty Series", newEmptySeries(), args{label: "1"}, -1},
		{"fail: label not in Series", s, args{label: "100"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)
			if got := tt.input.SelectLabel(tt.args.label); got != tt.want {
				t.Errorf("Series.SelectLabel() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("Series.SelectLabel() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestSeries_SelectLabels(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"}, Config{MultiIndex: []interface{}{[]int{1, 1, 1}, []string{"qux", "quux", "quuz"}}})
	type args struct {
		labels []string
		level  int
	}
	tests := []struct {
		name  string
		input *Series
		args  args
		want  []int
	}{
		{name: "pass", input: s, args: args{labels: []string{"1"}, level: 0}, want: []int{0, 1, 2}},
		{"pass", s, args{[]string{"qux", "quux"}, 1}, []int{0, 1}},
		{"fail: empty Series", newEmptySeries(), args{[]string{"1"}, 0}, []int{}},
		{"fail: label not in Series", s, args{[]string{"100"}, 0}, []int{}},
		{"fail: invalid level", s, args{[]string{"1"}, 100}, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)
			if got := tt.input.SelectLabels(tt.args.labels, tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Series.SelectLabels() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("Series.SelectLabels() returned no log message, want log due to fail")
				}
			}
		})
	}
}
