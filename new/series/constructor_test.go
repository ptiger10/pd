package series

import (
	"log"
	"reflect"
	"testing"

	constructIdx "github.com/ptiger10/pd/new/internal/index/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

// Calls New and panics if error. For use in testing
func mustNew(data interface{}, options ...Option) Series {
	s, err := New(data, options...)
	if err != nil {
		log.Panicf("mustNew returned an error: %v", err)
	}
	return s
}

// Float Tests
// ------------------------------------------------
func TestConstructor_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5})
	if err != nil {
		t.Error(err)
	}

	_, err = New([]float64{-1.5, 0, 1.5})
	if err != nil {
		t.Error(err)
	}
}

func TestConstructor_Name_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5}, Name("float32"))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, Name("float64"))
	if err != nil {
		t.Error(err)
	}
}

func TestConstructor_Index_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5}, Index([]float32{1, 2, 3}))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, Index([]float64{1, 2, 3}))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, Index([]float64{1}))
	if err == nil {
		t.Errorf("Returned nil error, want error due to mismatched value/index lengths")
	}
}

func TestConstructor_Kind_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.Float))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.Int))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.String))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.Bool))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.DateTime))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, Kind(kinds.Interface))
	if err == nil {
		t.Error("Returned nil error, want error due to unsupported conversion type")
	}
}

func TestMini_single(t *testing.T) {
	mini := miniIndex{
		data: []int{1, 2, 3},
		kind: kinds.Int,
		name: "test",
	}
	got, err := indexFromMiniIndex([]miniIndex{mini}, 3)
	if err != nil {
		t.Error(err)
	}
	want := constructIdx.New(
		constructIdx.SliceInt([]int64{1, 2, 3}, "test"),
	)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MiniIndex returned %v, want %v", got, want)
	}

}