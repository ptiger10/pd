package series

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/new/kinds"
)

// Len returns the length of the Series (including null values)
//
// Applies to: All
func (s Series) Len() int {
	all := s.values.All()
	return len(all)
}

// Describe the key details of the Series
//
// Applies to: All
func (s Series) Describe() {
	var err error
	// shared data
	origKind := s.Kind
	l := s.Len()
	valids := len(s.values.Valid())
	nulls := len(s.values.Null())
	len := fmt.Sprint(l)
	valid := fmt.Sprint(valids)
	null := fmt.Sprint(nulls)
	// type-specific data
	switch s.Kind {
	case kinds.Float, kinds.Int:
		precision := 4
		mean := fmt.Sprintf("%.*f", precision, s.Mean())
		min := fmt.Sprintf("%.*f", precision, s.Min())
		q1 := fmt.Sprintf("%.*f", precision, s.Quartile(1))
		q2 := fmt.Sprintf("%.*f", precision, s.Quartile(2))
		q3 := fmt.Sprintf("%.*f", precision, s.Quartile(3))
		max := fmt.Sprintf("%.*f", precision, s.Max())

		values := []string{len, valid, null, mean, min, q1, q2, q3, max}
		idx := Index([]string{"len", "valid", "null", "mean", "min", "25%", "50%", "75%", "max"})
		s, err = New(values, idx, Name("description"))

	case kinds.String:
		// value counts
		values := []string{len, valid, null}
		idx := Index([]string{"len", "valid", "null"})
		s, err = New(values, idx, Name("description"))
	case kinds.DateTime:
		// min and max
		values := []string{len, valid, null}
		idx := Index([]string{"len", "valid", "null"})
		s, err = New(values, idx, Name("description"))
	default:
		values := []string{len, valid, null}
		idx := Index([]string{"len", "valid", "null"})
		s, err = New(values, idx, Name("description"))
	}
	if err != nil {
		log.Printf("Internal error: s.Describe() could not construct Series: %v\nPlease open a Github issue.\n", err)
		return
	}
	// reset to pre-transformation Kind
	s.Kind = origKind
	fmt.Println(s)
	return
}
