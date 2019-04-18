package series

import (
	"reflect"
	"testing"
)

func Test_ConstructorUnsupported(t *testing.T) {
	_, err := New([]complex64{10})
	if err == nil {
		t.Errorf("Returned nil error when constructing unsupported series type, want error")
	}

	_, err = New([]interface{}{10})
	if err == nil {
		t.Errorf("Returned nil error when constructing interface without supplying series type, want error")
	}

	_, err = New([]interface{}{10}, SeriesType(reflect.Ptr))
	if err == nil {
		t.Errorf("Returned nil error when constructing interface with unsupported series type, want error")
	}
}
