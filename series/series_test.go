package series

import (
	"testing"
)

func TestUnsupported(t *testing.T) {
	_, err := New([]complex64{10})
	if err == nil {
		t.Errorf("Returned nil error when constructing unsupported series type, want error")
	}

	_, err = New([]interface{}{10})
	if err == nil {
		t.Errorf("Returned nil error when constructing interface without supplying series type, want error")
	}
}
