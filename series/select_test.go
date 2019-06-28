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

func TestSubset(t *testing.T) {
	tests := []struct {
		args    []int
		want    *Series
		wantErr bool
	}{
		{[]int{0}, MustNew("foo"), false},
		{[]int{1}, MustNew("bar", Config{Index: 1}), false},
		{[]int{0, 1}, MustNew([]string{"foo", "bar"}), false},
		{[]int{1, 0}, MustNew([]string{"bar", "foo"}, Config{Index: []int{1, 0}}), false},
		{[]int{}, newEmptySeries(), true},
		{[]int{3}, newEmptySeries(), true},
	}
	for _, tt := range tests {
		s := MustNew([]string{"foo", "bar", "baz"})
		got, err := s.Subset(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("s.Subset() error = %v, want %v for args %v", err, tt.wantErr, tt.args)
		}
		if !Equal(got, tt.want) {
			t.Errorf("s.Subset() got %v, want %v for args %v", got, tt.want, tt.args)
		}
	}
}
