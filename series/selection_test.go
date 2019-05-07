package series

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ptiger10/pd/options"
)

func TestMain(m *testing.M) {
	options.SetLogWarnings(false)
	exitCode := m.Run()
	options.RestoreDefaults()
	os.Exit(exitCode)
}

func TestAt_singleindex(t *testing.T) {
	var tests = []struct {
		input string
		want  Series
	}{
		{"10", mustNew([]string{"hot"}, Index([]int{10}))},
		{"100", mustNew([]string{"dog"}, Index([]int{100}))},
		{"1000", mustNew([]string{"log"}, Index([]int{1000}))},
	}
	for _, test := range tests {
		s, _ := New([]string{"hot", "dog", "log"}, Index([]int{10, 100, 1000}))
		got := s.At(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Returned Series %v, want %v", got, test.want)
		}
	}
}

func TestAt_fail(t *testing.T) {
	s, _ := New([]string{"hot", "dog", "log"}, Index([]int{10, 100, 1000}))
	got := s.At("NotPresent")
	if !reflect.DeepEqual(got, s) {
		t.Errorf("Returned %v, want original series", got)
	}
	fmt.Println(got)
}

func TestAt_multiindex(t *testing.T) {
	var tests = []struct {
		input string
		want  Series
	}{
		{"10", mustNew([]string{"hot"}, Index([]int{10}), Index([]string{"A"}))},
		{"100", mustNew([]string{"dog"}, Index([]int{100}), Index([]string{"B"}))},
		{"1000", mustNew([]string{"log"}, Index([]int{1000}), Index([]string{"C"}))},
	}
	for _, test := range tests {
		s, _ := New(
			[]string{"hot", "dog", "log"},
			Index([]int{10, 100, 1000}),
			Index([]string{"A", "B", "C"}),
		)
		got := s.At(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Returned Series %v, want %v", got, test.want)
		}
	}
}

func TestSelect_ByLabel_fail(t *testing.T) {
	s, _ := New([]int{1, 5, 10, 15}, Index([]int{0, 1, 0, 1}))
	_, err := s.Select(ByLabels([]string{"0, 2"}))
	if err == nil {
		t.Errorf("Returned nil error, expected error due to label indexing out range")
	}
}

func TestSelect_ByRows_fail(t *testing.T) {
	s, _ := New([]int{1, 5, 10, 15})
	_, err := s.Select(ByRows([]int{0, 4}))
	if err == nil {
		t.Errorf("Returned nil error, expected error due to row indexing out range")
	}
}

func TestSelect_ByLevels_fail(t *testing.T) {
	s, _ := New([]int{1, 5, 10, 15})
	_, err := s.Select(ByIndexLevels([]int{1}))
	if err == nil {
		t.Errorf("Returned nil error, expected error due to level indexing out range")
	}
}

func TestSelect_ByNames_fail(t *testing.T) {
	s, _ := New([]int{1, 5, 10, 15}, Index([]int{0, 1, 2, 3}, Name("foo")))
	_, err := s.Select(ByIndexNames([]string{"bar"}))
	if err == nil {
		t.Errorf("Returned nil error, expected error due to level indexing out range")
	}
}

// func TestSelect_ByNames_fail(t *testing.T) {
// 	s, _ := New([]int{1, 5, 10, 15}, Index([]int{0, 1, 2, 3}, Name("foo")))
// 	_, err := s.Select(ByIndexNames([]string{"bar"}))
// 	if err == nil {
// 		t.Errorf("Returned nil error, expected error due to level indexing out range")
// 	}
// }

func TestSelect_Failures(t *testing.T) {
	var tests = []struct {
		options  []SelectionOption
		errorMsg string
	}{
		{[]SelectionOption{ByRows([]int{0, 3})}, "row indexing out range"},
		{[]SelectionOption{ByLabels([]string{"0", "3"})}, "label indexing out range"},
		{[]SelectionOption{ByIndexLevels([]int{0, 2})}, "index level out range"},
		{[]SelectionOption{ByIndexNames([]string{"foo", "baz"})}, "index name not in index"},
		{[]SelectionOption{ByLabels([]string{"0"}), ByRows([]int{0})}, "multiple row selectors supplied"},
		{[]SelectionOption{ByIndexNames([]string{"foo"}), ByIndexLevels([]int{0})}, "multiple index level selectors supplied"},
		{[]SelectionOption{ByIndexLevels([]int{0, 1}), ByLabels([]string{"0"})}, "multiple levels plus labels supplied."},
	}
	for _, test := range tests {
		s, _ := New([]int{1, 5, 10}, Index([]int{0, 1, 2}, Name("foo")), Index([]string{"A", "B", "C"}, Name("bar")))
		_, err := s.Select(test.options...)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to %v", test.errorMsg)
		}
	}
}

func TestSelect_Success(t *testing.T) {
	var tests = []struct {
		options []SelectionOption
		desc    string
	}{
		{[]SelectionOption{ByRows([]int{0})}, "row only"},
		{[]SelectionOption{ByRows([]int{0, 2})}, "rows only"},
		{[]SelectionOption{ByLabels([]string{"0"})}, "label only"},
		{[]SelectionOption{ByLabels([]string{"0", "2"})}, "labels only"},
		{[]SelectionOption{ByIndexLevels([]int{0})}, "index level only"},
		{[]SelectionOption{ByIndexLevels([]int{0, 1})}, "index levels only"},
		{[]SelectionOption{ByIndexNames([]string{"foo"})}, "index name only"},
		{[]SelectionOption{ByIndexNames([]string{"foo", "bar"})}, "index names only"},
		{[]SelectionOption{ByRows([]int{0}), ByIndexLevels([]int{0})}, "rows and one index level"},
		{[]SelectionOption{ByRows([]int{0}), ByIndexLevels([]int{0, 1})}, "rows and multiple index levels"},
		{[]SelectionOption{ByRows([]int{0}), ByIndexNames([]string{"foo"})}, "rows and one index name"},
		{[]SelectionOption{ByRows([]int{0}), ByIndexNames([]string{"foo", "bar"})}, "rows and multiple index names"},
	}
	for _, test := range tests {
		s, _ := New([]int{1, 5, 10}, Index([]int{0, 1, 2}, Name("foo")), Index([]string{"A", "B", "C"}, Name("bar")))
		_, err := s.Select(test.options...)
		if err != nil {
			t.Errorf("Returned error when supplying %v, want no error", test.desc)
		}
	}
}
