package series

import (
	"testing"
	"time"

	"github.com/ptiger10/pd/kinds"
)

func TestTo_Float(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]interface{}{1.5, 1, "1", false, testDate})
	if err != nil {
		t.Error(err)
	}
	newS := s.To.Float()
	wantS, _ := New([]float64{1.5, 1.0, 1.0, 0, 1.5566688e+18})
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.To.Float() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.Float64
	if gotKind := newS.kind; gotKind != wantKind {
		t.Errorf("s.To.Float() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.Kind() == s.Kind() {
		t.Errorf("Conversion to float occurred in place, want copy only")
	}
}

func TestTo_Int(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]interface{}{1.5, 1, "1", false, testDate})
	if err != nil {
		t.Error(err)
	}
	newS := s.To.Int()
	wantS, _ := New([]int64{1, 1.0, 1.0, 0, 1.5566688e+18})
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.To.Int() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.Int64
	if gotKind := newS.kind; gotKind != wantKind {
		t.Errorf("s.To.Int() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.Kind() == s.Kind() {
		t.Errorf("Conversion to int occurred in place, want copy only")
	}
}

func TestTo_String(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]interface{}{1.5, 1, "1", false, testDate})
	if err != nil {
		t.Error(err)
	}
	newS := s.To.String()
	wantS, _ := New([]string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"})
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.To.String() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.String
	if gotKind := newS.kind; gotKind != wantKind {
		t.Errorf("s.To.String() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.Kind() == s.Kind() {
		t.Errorf("Conversion to string occurred in place, want copy only")
	}
}

func TestTo_Bool(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]interface{}{1.5, 1, "1", false, testDate})
	if err != nil {
		t.Error(err)
	}
	newS := s.To.Bool()
	wantS, _ := New([]bool{true, true, true, false, true})
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.To.Bool() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.Bool
	if gotKind := newS.kind; gotKind != wantKind {
		t.Errorf("s.To.Bool() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.Kind() == s.Kind() {
		t.Errorf("Conversion to bool occurred in place, want copy only")
	}
}

func TestTo_DateTime(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]interface{}{1.5, 1, "1", false, testDate})
	if err != nil {
		t.Error(err)
	}
	newS := s.To.DateTime()
	wantS, _ := New([]time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate})
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.To.DateTime() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.DateTime
	if gotKind := newS.kind; gotKind != wantKind {
		t.Errorf("s.To.DateTime() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.Kind() == s.Kind() {
		t.Errorf("Conversion to DateTime occurred in place, want copy only")
	}
}

func TestTo_Interface(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]interface{}{1.5, 1, "1", false, testDate})
	if err != nil {
		t.Error(err)
	}
	newS := s.To.Interface()
	wantS, _ := New([]interface{}{1.5, 1, "1", false, testDate})
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.To.DateTime() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.Interface
	if gotKind := newS.kind; gotKind != wantKind {
		t.Errorf("s.To.DateTime() returned kind %v, want %v", gotKind, wantKind)
	}
}

func TestIndexTo_Float(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
	if err != nil {
		t.Error(err)
	}
	newS := s.Index.To.Float()
	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]float64{1.5, 1.0, 1.0, 0, 1.5566688e+18}))
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.To.Float() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.Float64
	if gotKind := newS.index.Levels[0].Kind; gotKind != wantKind {
		t.Errorf("s.To.Float() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.index.Levels[0].Kind == s.index.Levels[0].Kind {
		t.Errorf("Conversion to float occurred in place, want copy only")
	}
}

func TestIndexTo_Int(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
	if err != nil {
		t.Error(err)
	}
	newS := s.Index.To.Int()
	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]int64{1, 1, 1, 0, 1.5566688e+18}))
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.IndexTo.Int() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.Int64
	if gotKind := newS.index.Levels[0].Kind; gotKind != wantKind {
		t.Errorf("s.IndexTo.Int() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.index.Levels[0].Kind == s.index.Levels[0].Kind {
		t.Errorf("Conversion to int occurred in place, want copy only")
	}
}

func TestIndexTo_String(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
	if err != nil {
		t.Error(err)
	}
	newS := s.Index.To.String()
	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}))
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.IndexTo.String() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.String
	if gotKind := newS.index.Levels[0].Kind; gotKind != wantKind {
		t.Errorf("s.IndexTo.String() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.index.Levels[0].Kind == s.index.Levels[0].Kind {
		t.Errorf("Conversion to string occurred in place, want copy only")
	}
}

func TestIndexTo_Bool(t *testing.T) {
	// testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]int{0, 1, 2, 3}, Idx([]interface{}{1.5, 1, "1", false}))
	if err != nil {
		t.Error(err)
	}

	newS := s.Index.To.Bool()
	wantS, _ := New([]int{0, 1, 2, 3}, Idx([]bool{true, true, true, false}))
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.IndexTo.Bool() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.Bool
	if gotKind := newS.index.Levels[0].Kind; gotKind != wantKind {
		t.Errorf("s.IndexTo.Bool() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.index.Levels[0].Kind == s.index.Levels[0].Kind {
		t.Errorf("Conversion to bool occurred in place, want copy only")
	}
}

func TestIndexTo_DateTime(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
	if err != nil {
		t.Error(err)
	}
	newS := s.Index.To.DateTime()
	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate}))
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.IndexTo.DateTime() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.DateTime
	if gotKind := newS.index.Levels[0].Kind; gotKind != wantKind {
		t.Errorf("s.IndexTo.DateTime() returned kind %v, want %v", gotKind, wantKind)
	}
	if newS.index.Levels[0].Kind == s.index.Levels[0].Kind {
		t.Errorf("Conversion to DateTime occurred in place, want copy only")
	}
}

func TestIndexTo_Interface(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
	if err != nil {
		t.Error(err)
	}
	newS := s.Index.To.Interface()
	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
	if !seriesEquals(newS, wantS) {
		t.Errorf("s.IndexTo.Interface() returned %v, want %v", newS, wantS)
	}
	wantKind := kinds.Interface
	if gotKind := newS.index.Levels[0].Kind; gotKind != wantKind {
		t.Errorf("s.IndexTo.Interface() returned kind %v, want %v", gotKind, wantKind)
	}
}

// func TestConvertIndexMulti(t *testing.T) {
// 	var tests = []struct {
// 		convertTo kinds.Kind
// 		lvl       int
// 	}{
// 		{kinds.Float64, 0},
// 		{kinds.Float64, 1},
// 		{kinds.Int, 0},
// 		{kinds.Int, 1},
// 		{kinds.String, 0},
// 		{kinds.String, 1},
// 		{kinds.Bool, 0},
// 		{kinds.Bool, 1},
// 		{kinds.DateTime, 0},
// 		{kinds.DateTime, 1},
// 	}
// 	for _, test := range tests {
// 		s, err := New([]interface{}{1, 2, 3}, Idx([]int{1, 2, 3}), Idx([]int{10, 20, 30}))
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		newS, err := s.IndexLevelTo(test.lvl, test.convertTo)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if newS.index.Levels[test.lvl].Kind != test.convertTo {
// 			t.Errorf("Conversion of Series with multiIndex level %v to %v returned %v, want %v", test.lvl, test.convertTo, newS.index.Levels[test.lvl].Kind, test.convertTo)
// 		}
// 		// excludes Int because the original test Index is int
// 		if test.convertTo != kinds.Int {
// 			if s.index.Levels[test.lvl].Kind == newS.index.Levels[test.lvl].Kind {
// 				t.Errorf("Conversion to %v occurred in place, want copy only", test.convertTo)
// 			}
// 		}
// 	}
// }
