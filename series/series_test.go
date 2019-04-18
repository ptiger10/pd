package series

import (
	"testing"
)

func TestUnsupported(t *testing.T) {
	_, err := New([]complex64{10})
	if err == nil {
		t.Errorf("Returned nil error when constructing unsupported series type, want error")
	}
}
