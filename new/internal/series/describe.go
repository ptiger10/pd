package series

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/new/internal/values"
	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
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
	var values values.Values
	var idx Option

	// common data
	l := s.Len()
	valids := len(s.values.Valid())
	nulls := len(s.values.Null())
	len := fmt.Sprintln(l)
	valid := fmt.Sprintln(valids)
	null := fmt.Sprintln(nulls)
	// type-specific data
	switch s.Kind {
	case kinds.Float, kinds.Int:
		precision := 4
		mean := fmt.Sprintf("%.*f\n", precision, s.Mean())
		min := fmt.Sprintf("%.*f\n", precision, s.Min())
		q1 := fmt.Sprintf("%.*f\n", precision, s.Quartile(1))
		q2 := fmt.Sprintf("%.*f\n", precision, s.Quartile(2))
		q3 := fmt.Sprintf("%.*f\n", precision, s.Quartile(3))
		max := fmt.Sprintf("%.*f\n", precision, s.Max())

		idx = Index([]string{"len", "valid", "null", "mean", "min", "25%", "50%", "75%", "max"}, Name("statistic"))
		values = constructVal.SliceString([]string{
			len, valid, null, mean, min, q1, q2, q3, max})
	case kinds.String:
		// value counts
		idx = Index([]string{"len", "valid", "null"}, Name("statistic"))
		values = constructVal.SliceString([]string{len, valid, null})
	case kinds.DateTime:
		// min and max
		idx = Index([]string{"len", "valid", "null"}, Name("statistic"))
		values = constructVal.SliceString([]string{len, valid, null})
	default:
		idx = Index([]string{"len", "valid", "null"}, Name("statistic"))
		values = constructVal.SliceString([]string{len, valid, null})
	}
	s, err := New(values, idx)
	if err != nil {
		log.Printf("Internal error: s.Describe() could not construct Series: %v\n", err)
		return
	}
	fmt.Print(s)

}
