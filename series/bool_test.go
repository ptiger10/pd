package series

import (
	"testing"
)

func TestAdd_BoolUnsupported(t *testing.T) {
	s, _ := New([]bool{true, true, false, true})
	_, err := s.AddConst(1)
	if err == nil {
		t.Error("Returned nil error when adding constant to Bool, want error")
	}
}
