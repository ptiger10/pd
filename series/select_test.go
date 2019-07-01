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
func TestSubset(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"})
	misaligned := MustNew([]string{"foo", "bar"})
	misaligned.index = misaligned.index.Subset([]int{0})

	tests := []struct {
		name    string
		input   *Series
		args    []int
		want    *Series
		wantErr bool
	}{
		{name: "pass", input: s, args: []int{0}, want: MustNew("foo"), wantErr: false},
		{"pass", s, []int{1}, MustNew("bar", Config{Index: 1}), false},
		{"pass", s, []int{0, 1}, MustNew([]string{"foo", "bar"}), false},
		{"pass", s, []int{1, 0}, MustNew([]string{"bar", "foo"}, Config{Index: []int{1, 0}}), false},
		{"fail: misaligned", misaligned, []int{1}, newEmptySeries(), true},
		{"fail: empty", s, []int{}, newEmptySeries(), true},
		{"fail: invalid", s, []int{3}, newEmptySeries(), true},
		{"fail: nil positions", s, nil, newEmptySeries(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input.Copy()
			sArchive := tt.input.Copy()

			err := s.InPlace.Subset(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Series.InPlace.Subset() error = %v, want %v", err, tt.wantErr)
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(s, tt.want) {
					t.Errorf("Series.InPlace.Subset() got %v, want %v", s, tt.want)
				}
			}

			sCopy, err := sArchive.Subset(tt.args)
			if !Equal(sCopy, tt.want) {
				t.Errorf("Series.Subset() got %v, want %v", sCopy, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Series.Subset() error = %v, want %v", err, tt.wantErr)
			}
			if Equal(sArchive, sCopy) {
				t.Errorf("Series.Subset() retained access to original, want copy")
			}
		})

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
	s := MustNew([]string{"foo", "bar", "baz"})
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
		{name: "ascending", input: s, args: args{start: 0, end: 2}, want: MustNew([]string{"foo", "bar", "baz"})},
		{"single", s, args{1, 1}, MustNew([]string{"bar"}, Config{Index: []int{1}})},
		{"partial", s, args{1, 2}, MustNew([]string{"bar", "baz"}, Config{Index: []int{1, 2}})},
		{"descending", s, args{2, 0}, MustNew([]string{"baz", "bar", "foo"}, Config{Index: []int{2, 1, 0}})},
		{"fail: partial invalid", s, args{10, 0}, MustNew(newEmptySeries)},
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
		name     string
		input    *Series
		args     args
		want     int
		wantFail bool
	}{
		{name: "pass", input: s, args: args{label: "1"}, want: 0, wantFail: false},
		{"fail: empty Series", &Series{}, args{label: "1"}, -1, true},
		{"fail: label not in Series", s, args{label: "100"}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)
			if got := tt.input.SelectLabel(tt.args.label); got != tt.want {
				t.Errorf("Series.SelectLabel() = %v, want %v", got, tt.want)
			}
			if tt.wantFail {
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
		name     string
		input    *Series
		args     args
		want     []int
		wantFail bool
	}{
		{name: "pass", input: s, args: args{labels: []string{"1"}, level: 0}, want: []int{0, 1, 2}, wantFail: false},
		{"pass", s, args{[]string{"qux", "quux"}, 1}, []int{0, 1}, false},
		{"fail: empty Series", &Series{}, args{[]string{"1"}, 0}, []int{}, true},
		{"fail: label not in Series", s, args{[]string{"100"}, 0}, []int{}, true},
		{"fail: invalid level", s, args{[]string{"1"}, 100}, []int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)
			if got := tt.input.SelectLabels(tt.args.labels, tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Series.SelectLabels() = %v, want %v", got, tt.want)
			}
			if tt.wantFail {
				if buf.String() == "" {
					t.Errorf("Series.SelectLabels() returned no log message, want log due to fail")
				}
			}
		})
	}
}
