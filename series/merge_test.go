package series

import "testing"

func TestJoin(t *testing.T) {
	s := MustNew([]int{1, 2, 3})
	s2 := MustNew([]float64{4, 5, 6})
	s3 := s.Join(s2)
	want := MustNew([]int{1, 2, 3, 4, 5, 6}, Config{Index: []int{0, 1, 2, 0, 1, 2}})
	if !Equal(s3, want) {
		t.Errorf("s.Join() returned %v, want %v", s3, want)
	}
}

func TestJoinEmpty(t *testing.T) {
	s := MustNew(nil)
	s2 := MustNew([]float64{4, 5, 6})
	s3 := s.Join(s2)
	want := MustNew([]float64{4, 5, 6}, Config{Index: []int{0, 1, 2}})
	if !Equal(s3, want) {
		t.Errorf("s.Join() returned %v, want %v", s3, want)
	}
}

// func Test_InPlace_Join(t *testing.T) {
// 	s, _ := New([]int{1, 2, 3})
// 	sCopy, _ := New([]float64{4, 5, 6})
// 	s.InPlace.Join(sCopy)
// 	want := MustNew([]int{1, 2, 3, 4, 5, 6}, Idx([]int{0, 1, 2, 0, 1, 2}))
// 	if !Equal(s, want) {
// 		t.Errorf("s.InPlace.Join() got %v, want %v", s, want)
// 	}
// }

// // func Test_InPlace_replace(t *testing.T) {
// // 	s, _ := New(1, options.Name("foo"))
// // 	sCopy, _ := New(2, options.Name("bar"))
// // 	s.InPlace.s.replace(sCopy)
// // 	if !Equal(s, *sCopy) {
// // 		t.Errorf("s.InPlace.replace() got %v, want %v", s, sCopy)
// // 	}
// // }

// func Test_InPlace_Join_EmptyBase(t *testing.T) {
// 	s, _ := New(nil)
// 	sCopy, _ := New([]float64{4, 5, 6})
// 	s.InPlace.Join(sCopy)
// 	want := MustNew([]float64{4, 5, 6}, Idx([]int{0, 1, 2}))
// 	if !Equal(s, want) {
// 		t.Errorf("s.InPlace.Join() got %v, want %v", sCopy, want)
// 	}
// }
