package series

import "testing"

func Test_XS(t *testing.T) {

	s := MustNew([]float64{1, 2}, Config{MultiIndex: []interface{}{[]string{"foo", "bar"}, []string{"baz", "qux"}}})
	got, err := s.XS([]int{1}, []int{1})
	if err != nil {
		t.Error(err)
	}
	want := MustNew(2.0, Config{Index: "qux"})
	if !Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
