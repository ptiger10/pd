package series

// func TestIndexTo_Float(t *testing.T) {
// 	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	newS := s.Index.ToFloat64()
// 	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]float64{1.5, 1.0, 1.0, 0, 1.5566688e+18}))
// 	if !Equal(newS, wantS) {
// 		t.Errorf("s.ToFloat64() returned %v, want %v", newS, wantS)
// 	}
// 	wantDataType := options.Float64
// 	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
// 		t.Errorf("s.ToFloat64() returned kind %v, want %v", gotDataType, wantDataType)
// 	}
// 	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
// 		t.Errorf("Conversion to float occurred in place, want copy only")
// 	}
// }

// func TestIndexTo_Int(t *testing.T) {
// 	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	newS := s.Index.ToInt64()
// 	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]int64{1, 1, 1, 0, 1.5566688e+18}))
// 	if !Equal(newS, wantS) {
// 		t.Errorf("s.IndexToInt64() returned %v, want %v", newS, wantS)
// 	}
// 	wantDataType := options.Int64
// 	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
// 		t.Errorf("s.IndexToInt64() returned kind %v, want %v", gotDataType, wantDataType)
// 	}
// 	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
// 		t.Errorf("Conversion to int occurred in place, want copy only")
// 	}
// }

// func TestIndexTo_String(t *testing.T) {
// 	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	newS := s.Index.ToString()
// 	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]string{"1.5", "1", "1", "false", "2019-05-01 00:00:00 +0000 UTC"}))
// 	if !Equal(newS, wantS) {
// 		t.Errorf("s.IndexToString() returned %v, want %v", newS, wantS)
// 	}
// 	wantDataType := options.String
// 	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
// 		t.Errorf("s.IndexToString() returned kind %v, want %v", gotDataType, wantDataType)
// 	}
// 	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
// 		t.Errorf("Conversion to string occurred in place, want copy only")
// 	}
// }

// func TestIndexTo_Bool(t *testing.T) {
// 	// testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	s, err := New([]int{0, 1, 2, 3}, Idx([]interface{}{1.5, 1, "1", false}))
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	newS := s.Index.ToBool()
// 	wantS, _ := New([]int{0, 1, 2, 3}, Idx([]bool{true, true, true, false}))
// 	if !Equal(newS, wantS) {
// 		t.Errorf("s.IndexToBool() returned %v, want %v", newS, wantS)
// 	}
// 	wantDataType := options.Bool
// 	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
// 		t.Errorf("s.IndexToBool() returned kind %v, want %v", gotDataType, wantDataType)
// 	}
// 	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
// 		t.Errorf("Conversion to bool occurred in place, want copy only")
// 	}
// }

// func TestIndexTo_DateTime(t *testing.T) {
// 	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	epochDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
// 	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	newS := s.Index.ToDateTime()
// 	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]time.Time{epochDate, epochDate, time.Time{}, epochDate, testDate}))
// 	if !Equal(newS, wantS) {
// 		t.Errorf("s.IndexToDateTime() returned %v, want %v", newS, wantS)
// 	}
// 	wantDataType := options.DateTime
// 	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
// 		t.Errorf("s.IndexToDateTime() returned kind %v, want %v", gotDataType, wantDataType)
// 	}
// 	if newS.index.Levels[0].DataType == s.index.Levels[0].DataType {
// 		t.Errorf("Conversion to DateTime occurred in place, want copy only")
// 	}
// }

// func TestIndexTo_Interface(t *testing.T) {
// 	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
// 	s, err := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	newS := s.Index.ToInterface()
// 	wantS, _ := New([]int{0, 1, 2, 3, 4}, Idx([]interface{}{1.5, 1, "1", false, testDate}))
// 	if !Equal(newS, wantS) {
// 		t.Errorf("s.IndexToInterface() returned %v, want %v", newS, wantS)
// 	}
// 	wantDataType := options.Interface
// 	if gotDataType := newS.index.Levels[0].DataType; gotDataType != wantDataType {
// 		t.Errorf("s.IndexToInterface() returned kind %v, want %v", gotDataType, wantDataType)
// 	}
// }

// func TestIndexAt(t *testing.T) {
// 	s, _ := New([]int{0, 1, 2})
// 	got, _ := s.Index.At(0, 0)
// 	want := int64(0)
// 	if got.(int64) != want {
// 		t.Errorf("IndexAt() got %v, want %v", got, want)
// 	}
// }

// func Test_Index_Sort(t *testing.T) {
// 	var tests = []struct {
// 		desc string
// 		orig *Series
// 		asc  bool
// 		want *Series
// 	}{
// 		{"float", MustNew([]float64{1, 3, 5}, Idx([]int{2, 0, 1})), true, MustNew([]float64{3, 5, 1}, Idx([]int{0, 1, 2}))},
// 		{"float reverse", MustNew([]float64{1, 3, 5}, Idx([]int{2, 0, 1})), false, MustNew([]float64{1, 5, 3}, Idx([]int{2, 1, 0}))},
// 	}
// 	for _, test := range tests {
// 		s := test.orig
// 		s.Index.Sort(test.asc)
// 		if !Equal(s, test.want) {
// 			t.Errorf("series.Index.Sort() test %v got %v, want %v", test.desc, s, test.want)
// 		}
// 	}
// }

// func TestConvertIndexMulti(t *testing.T) {
// 	var tests = []struct {
// 		convertTo options.DataType
// 		lvl       int
// 	}{
// 		{options.Float64, 0},
// 		{options.Float64, 1},
// 		{options.Int, 0},
// 		{options.Int, 1},
// 		{options.String, 0},
// 		{options.String, 1},
// 		{options.Bool, 0},
// 		{options.Bool, 1},
// 		{options.DateTime, 0},
// 		{options.DateTime, 1},
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
// 			t.Errorf("Conversion of Series with multiIndex level %v to %v got %v, want %v", test.lvl, test.convertTo, newS.index.Levels[test.lvl].DataType, test.convertTo)
// 		}
// 		// excludes Int because the original test Index is int
// 		if test.convertTo != options.Int {
// 			if s.index.Levels[test.lvl].DataType == newS.index.Levels[test.lvl].DataType {
// 				t.Errorf("Conversion to %v occurred in place, want copy only", test.convertTo)
// 			}
// 		}
// 	}
// }
