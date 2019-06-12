package series

import (
	"testing"
	"time"

	"github.com/ptiger10/pd/datatypes"
	"github.com/ptiger10/pd/opt"
)

func TestInsert(t *testing.T) {
	var tests = []struct {
		pos  int
		val  interface{}
		idx  []interface{}
		want Series
	}{
		{0, "baz", []interface{}{"C", 3},
			mustNew([]string{"baz", "foo", "bar"}, Idx([]string{"C", "A", "B"}), Idx([]int{3, 1, 2}))},
		{1, "baz", []interface{}{"C", 3},
			mustNew([]string{"foo", "baz", "bar"}, Idx([]string{"A", "C", "B"}), Idx([]int{1, 3, 2}))},
		{2, "baz", []interface{}{"C", 3},
			mustNew([]string{"foo", "bar", "baz"}, Idx([]string{"A", "B", "C"}), Idx([]int{1, 2, 3}))},
	}
	for _, test := range tests {
		s, _ := New([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		newS, err := s.Insert(test.pos, test.val, test.idx)
		if err != nil {
			t.Errorf("s.Insert(): %v", err)
		}
		if !seriesEquals(newS, test.want) {
			t.Errorf("s.insert() returned %v, want %v", s, test.want)
		}
		if seriesEquals(newS, s) {
			t.Error("s.insert() maintained reference to original Series, want fresh copy")
		}
	}
}

func TestAppend(t *testing.T) {
	var tests = []struct {
		val  interface{}
		idx  []interface{}
		want Series
	}{
		{"baz", []interface{}{"C", 3},
			mustNew([]string{"foo", "bar", "baz"}, Idx([]string{"A", "B", "C"}), Idx([]int{1, 2, 3}))},
	}
	for _, test := range tests {
		s, _ := New([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		newS := s.Append(test.val, test.idx)
		if !seriesEquals(newS, test.want) {
			t.Errorf("s.Append() returned %v, want %v", s, test.want)
		}
		if seriesEquals(newS, s) {
			t.Error("s.Append() maintained reference to original Series, want fresh copy")
		}
	}
}

func TestDrop(t *testing.T) {
	var tests = []struct {
		pos  int
		want Series
	}{
		{0, mustNew([]string{"bar"}, Idx([]string{"B"}), Idx([]int{2}))},
		{1, mustNew([]string{"foo"}, Idx([]string{"A"}), Idx([]int{1}))},
	}
	for _, test := range tests {
		s, _ := NewPointer([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		newS, err := s.Drop(test.pos)
		if err != nil {
			t.Errorf("s.Drop(): %v", err)
		}
		if !seriesEquals(newS, test.want) {
			t.Errorf("s.Drop() returned %v, want %v", s, test.want)
		}
		if seriesEquals(newS, *s) {
			t.Error("s.Drop() maintained reference to original Series, want fresh copy")
		}
	}
}

func TestJoin(t *testing.T) {
	s, _ := New([]int{1, 2, 3})
	s2, _ := New([]float64{4, 5, 6})
	s3 := s.Join(s2)
	want := mustNew([]int{1, 2, 3, 4, 5, 6}, Idx([]int{0, 1, 2, 0, 1, 2}))
	if !seriesEquals(s3, want) {
		t.Errorf("s.Join() returned %v, want %v", s3, want)
	}
}

func TestInsertInPlace(t *testing.T) {
	var tests = []struct {
		pos  int
		val  interface{}
		idx  []interface{}
		want Series
	}{
		{0, "baz", []interface{}{"C", 3},
			mustNew([]string{"baz", "foo", "bar"}, Idx([]string{"C", "A", "B"}), Idx([]int{3, 1, 2}))},
		{1, "baz", []interface{}{"C", 3},
			mustNew([]string{"foo", "baz", "bar"}, Idx([]string{"A", "C", "B"}), Idx([]int{1, 3, 2}))},
		{2, "baz", []interface{}{"C", 3},
			mustNew([]string{"foo", "bar", "baz"}, Idx([]string{"A", "B", "C"}), Idx([]int{1, 2, 3}))},
	}
	for _, test := range tests {
		s, err := New([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		if err != nil {
			t.Error(err)
		}
		s.InPlace.Insert(test.pos, test.val, test.idx)
		if !seriesEquals(s, test.want) {
			t.Errorf("s.insert() returned %v, want %v", s, test.want)
		}
	}
}

func TestAppendInPlace(t *testing.T) {
	var tests = []struct {
		val  interface{}
		idx  []interface{}
		want Series
	}{
		{"baz", []interface{}{"C", 3},
			mustNew([]string{"foo", "bar", "baz"}, Idx([]string{"A", "B", "C"}), Idx([]int{1, 2, 3}))},
	}
	for _, test := range tests {
		s, _ := New([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		s.InPlace.Append(test.val, test.idx)
		if !seriesEquals(s, test.want) {
			t.Errorf("s.Append() returned %v, want %v", s, test.want)
		}
	}
}

func TestDropInPlace(t *testing.T) {
	var tests = []struct {
		pos  int
		want Series
	}{
		{0, mustNew([]string{"bar"}, Idx([]string{"B"}), Idx([]int{2}))},
		{1, mustNew([]string{"foo"}, Idx([]string{"A"}), Idx([]int{1}))},
	}
	for _, test := range tests {
		s, _ := NewPointer([]string{"foo", "bar"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
		s.InPlace.Drop(test.pos)
		if !seriesEquals(*s, test.want) {
			t.Errorf("s.insert() returned %v, want %v", s, test.want)
		}
	}
}

func Test_InPlace_Join(t *testing.T) {
	s, _ := New([]int{1, 2, 3})
	s2, _ := New([]float64{4, 5, 6})
	s.InPlace.Join(s2)
	want := mustNew([]int{1, 2, 3, 4, 5, 6}, Idx([]int{0, 1, 2, 0, 1, 2}))
	if !seriesEquals(s, want) {
		t.Errorf("s.Join() returned %v, want %v", s, want)
	}
}

func Test_InPlace_replace(t *testing.T) {
	s, _ := NewPointer(1, opt.Name("foo"))
	s2, _ := NewPointer(2, opt.Name("bar"))
	s.InPlace.s.replace(s2)
	if !seriesEquals(*s, *s2) {
		t.Errorf("s.InPlace.replace() returned %v, want %v", s, s2)
	}
}

func Test_InPlace_Join_EmptyBase(t *testing.T) {
	s, _ := NewPointer(nil)
	s2, _ := New([]float64{4, 5, 6})
	s.InPlace.Join(s2)
	want := mustNew([]float64{4, 5, 6}, Idx([]int{0, 1, 2}))
	if !seriesEquals(*s, want) {
		t.Errorf("s.InPlace.Join() returned %v, want %v", s2, want)
	}
}

func Test_InPlace_Sort(t *testing.T) {
	var tests = []struct {
		desc string
		orig Series
		asc  bool
		want Series
	}{
		{"float", mustNew([]float64{3, 5, 1}), true, mustNew([]float64{1, 3, 5}, Idx([]int{2, 0, 1}))},
		{"float reverse", mustNew([]float64{3, 5, 1}), false, mustNew([]float64{5, 3, 1}, Idx([]int{1, 0, 2}))},

		{"int", mustNew([]int{3, 5, 1}), true, mustNew([]int{1, 3, 5}, Idx([]int{2, 0, 1}))},
		{"int reverse", mustNew([]int{3, 5, 1}), false, mustNew([]int{5, 3, 1}, Idx([]int{1, 0, 2}))},

		{"string", mustNew([]string{"3", "5", "1"}), true, mustNew([]string{"1", "3", "5"}, Idx([]int{2, 0, 1}))},
		{"string reverse", mustNew([]string{"3", "5", "1"}), false, mustNew([]string{"5", "3", "1"}, Idx([]int{1, 0, 2}))},

		{"bool", mustNew([]bool{false, true, false}), true, mustNew([]bool{false, false, true}, Idx([]int{0, 2, 1}))},
		{"bool reverse", mustNew([]bool{false, true, false}), false, mustNew([]bool{true, false, false}, Idx([]int{1, 0, 2}))},

		{
			"datetime",
			mustNew([]time.Time{time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}),
			true,
			mustNew([]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC)}, Idx([]int{2, 0, 1})),
		},
		{
			"datetime reverse",
			mustNew([]time.Time{time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}),
			false,
			mustNew([]time.Time{time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 2, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}, Idx([]int{1, 0, 2})),
		},
	}
	for _, test := range tests {
		s := test.orig
		s.InPlace.Sort(test.asc)
		if !seriesEquals(s, test.want) {
			t.Errorf("series.Sort() test %v returned %v, want %v", test.desc, s, test.want)
		}
	}
}

func Test_Index_Sort(t *testing.T) {
	var tests = []struct {
		desc string
		orig Series
		asc  bool
		want Series
	}{
		{"float", mustNew([]float64{1, 3, 5}, Idx([]int{2, 0, 1})), true, mustNew([]float64{3, 5, 1}, Idx([]int{0, 1, 2}))},
		{"float reverse", mustNew([]float64{1, 3, 5}, Idx([]int{2, 0, 1})), false, mustNew([]float64{1, 5, 3}, Idx([]int{2, 1, 0}))},
	}
	for _, test := range tests {
		s := test.orig
		s.Index.Sort(test.asc)
		if !seriesEquals(s, test.want) {
			t.Errorf("series.Index.Sort() test %v returned %v, want %v", test.desc, s, test.want)
		}
	}
}

// [START Convert tests]

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
	wantDataType := datatypes.Float64
	if gotDataType := newS.datatype; gotDataType != wantDataType {
		t.Errorf("s.To.Float() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.DataType() == s.DataType() {
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
	wantDataType := datatypes.Int64
	if gotDataType := newS.datatype; gotDataType != wantDataType {
		t.Errorf("s.To.Int() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.DataType() == s.DataType() {
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
	wantDataType := datatypes.String
	if gotDataType := newS.datatype; gotDataType != wantDataType {
		t.Errorf("s.To.String() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.DataType() == s.DataType() {
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
	wantDataType := datatypes.Bool
	if gotDataType := newS.datatype; gotDataType != wantDataType {
		t.Errorf("s.To.Bool() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.DataType() == s.DataType() {
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
	wantDataType := datatypes.DateTime
	if gotDataType := newS.datatype; gotDataType != wantDataType {
		t.Errorf("s.To.DateTime() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.DataType() == s.DataType() {
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
	wantDataType := datatypes.Interface
	if gotDataType := newS.datatype; gotDataType != wantDataType {
		t.Errorf("s.To.DateTime() returned kind %v, want %v", gotDataType, wantDataType)
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
	wantDataType := datatypes.Float64
	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
		t.Errorf("s.To.Float() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
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
	wantDataType := datatypes.Int64
	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
		t.Errorf("s.IndexTo.Int() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
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
	wantDataType := datatypes.String
	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
		t.Errorf("s.IndexTo.String() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
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
	wantDataType := datatypes.Bool
	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
		t.Errorf("s.IndexTo.Bool() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
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
	wantDataType := datatypes.DateTime
	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
		t.Errorf("s.IndexTo.DateTime() returned kind %v, want %v", gotDataType, wantDataType)
	}
	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
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
	wantDataType := datatypes.Interface
	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
		t.Errorf("s.IndexTo.Interface() returned kind %v, want %v", gotDataType, wantDataType)
	}
}

// func TestConvertIndexMulti(t *testing.T) {
// 	var tests = []struct {
// 		convertTo datatypes.DataType
// 		lvl       int
// 	}{
// 		{datatypes.Float64, 0},
// 		{datatypes.Float64, 1},
// 		{datatypes.Int, 0},
// 		{datatypes.Int, 1},
// 		{datatypes.String, 0},
// 		{datatypes.String, 1},
// 		{datatypes.Bool, 0},
// 		{datatypes.Bool, 1},
// 		{datatypes.DateTime, 0},
// 		{datatypes.DateTime, 1},
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
// 		if newS.index.Levels[test.lvl].DataType != test.convertTo {
// 			t.Errorf("Conversion of Series with multiIndex level %v to %v returned %v, want %v", test.lvl, test.convertTo, newS.index.Levels[test.lvl].DataType, test.convertTo)
// 		}
// 		// excludes Int because the original test Index is int
// 		if test.convertTo != datatypes.Int {
// 			if s.index.Levels[test.lvl].DataType == newS.index.Levels[test.lvl].DataType {
// 				t.Errorf("Conversion to %v occurred in place, want copy only", test.convertTo)
// 			}
// 		}
// 	}
// }

// [END Convert tests]
