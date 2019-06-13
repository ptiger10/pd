package series

import (
	"fmt"
	"time"

	"github.com/ptiger10/pd/options"
)

// Describe the key details of the Series.
func (s *Series) Describe() {
	var values interface{}
	var idx interface{}
	// dt preserves the datatype of the original Series
	dt := s.datatype

	// shared data
	l := s.Len()
	valids := len(s.valid())
	nulls := len(s.null())
	length := fmt.Sprint(l)
	valid := fmt.Sprint(valids)
	null := fmt.Sprint(nulls)
	// type-specific data
	switch s.datatype {
	case options.Float64, options.Int64:
		precision := options.GetDisplayFloatPrecision()
		mean := fmt.Sprintf("%.*f", precision, s.Math.Mean())
		min := fmt.Sprintf("%.*f", precision, s.Math.Min())
		q1 := fmt.Sprintf("%.*f", precision, s.Math.Quartile(1))
		q2 := fmt.Sprintf("%.*f", precision, s.Math.Quartile(2))
		q3 := fmt.Sprintf("%.*f", precision, s.Math.Quartile(3))
		max := fmt.Sprintf("%.*f", precision, s.Math.Max())

		values = []string{length, valid, null, mean, min, q1, q2, q3, max}
		idx = []string{"len", "valid", "null", "mean", "min", "25%", "50%", "75%", "max"}

	case options.String:
		unique := fmt.Sprint(len(s.UniqueVals()))
		values = []string{length, valid, null, unique}
		idx = []string{"len", "valid", "null", "unique"}

	case options.Bool:
		precision := options.GetDisplayFloatPrecision()
		sum := fmt.Sprintf("%.*f", precision, s.Math.Sum())
		mean := fmt.Sprintf("%.*f", precision, s.Math.Mean())
		values = []string{length, valid, null, sum, mean}
		idx = []string{"len", "valid", "null", "sum", "mean"}

	case options.DateTime:
		unique := fmt.Sprint(len(s.UniqueVals()))
		earliest := fmt.Sprint(s.Earliest())
		latest := fmt.Sprint(s.Latest())
		values = []string{length, valid, null, unique, earliest, latest}
		idx = []string{"len", "valid", "null", "unique", "earliest", "latest"}

	// Interface or None
	default:
		values = []string{length, valid, null}
		idx = []string{"len", "valid", "null"}
	}

	s, err := NewWithConfig(Config{Name: s.Name}, values, Idx(idx))
	if err != nil {
		fmt.Printf("series.Describe(): %v", err)
	}
	s.datatype = dt
	fmt.Println(s)
	return
}

// ValueCounts returns a map of non-null value labels to number of occurrences in the Series.
//
// Applies to: All
func (s *Series) ValueCounts() map[string]int {
	valid, _ := s.in(s.valid())
	vals := valid.all()
	counter := make(map[string]int)
	for _, val := range vals {
		counter[fmt.Sprint(val)]++
	}
	return counter
}

// UniqueVals returns a de-duplicated list all element values (as strings) that appear in the Series.
//
// Applies to: All
func (s *Series) UniqueVals() []string {
	var ret []string
	counter := s.ValueCounts()
	for k := range counter {
		ret = append(ret, k)
	}
	return ret
}

// Earliest returns the earliest non-null time.Time{} in the Series
//
// Applies to: time.Time. If inapplicable, defaults to time.Time{}.
func (s *Series) Earliest() time.Time {
	earliest := time.Time{}
	vals := s.validVals()
	switch s.datatype {
	case options.DateTime:
		data := ensureDateTime(vals)
		for _, t := range data {
			if earliest == (time.Time{}) || t.Before(earliest) {
				earliest = t
			}
		}
		return earliest
	default:
		return earliest

	}
}

// Latest returns the latest non-null time.Time{} in the Series
//
// Applies to: time.Time. If inapplicable, defaults to time.Time{}.
func (s *Series) Latest() time.Time {
	latest := time.Time{}
	vals := s.validVals()
	switch s.datatype {
	case options.DateTime:
		data := ensureDateTime(vals)
		for _, t := range data {
			if latest == (time.Time{}) || t.After(latest) {
				latest = t
			}
		}
		return latest
	default:
		return latest

	}
}
