package series

import (
	"log"
	"reflect"
	"testing"

	"github.com/ptiger10/pd/internal/config"
	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/kinds"
	"github.com/ptiger10/pd/opt"
)

// Calls New and panics if error. For use in testing
func mustNew(data interface{}, options ...opt.ConstructorOption) Series {
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

	_, err = New([]float32{-1.5, 0, 1.5}, opt.Name("float32"))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, opt.Name("float64"))
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

	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Float))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Int))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.String))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Bool))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.DateTime))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Interface))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Unsupported))
	if err == nil {
		t.Error("Returned nil error, want error due to unsupported conversion type")
	}

}

func TestMini_single(t *testing.T) {
	mini := config.MiniIndex{
		Data: []int{1, 2, 3},
		Kind: kinds.Int,
		Name: "test",
	}
	got, err := indexFromMiniIndex([]config.MiniIndex{mini}, 3)
	if err != nil {
		t.Error(err)
	}
	lvl, err := index.NewLevelFromSlice([]int64{1, 2, 3}, "test")
	want := index.New(lvl)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MiniIndex returned %v, want %v", got, want)
	}

}
