package series

import (
	"log"

	"github.com/ptiger10/pd/new/datatypes"
	"github.com/ptiger10/pd/new/internal/values"
)

func (s Series) Len() int {
	switch s.Kind {
	// case Float:
	// 	vals := s.Values.(floatValues)
	// 	return vals.count()
	// case Int:
	// 	vals := s.Values.(intValues)
	// 	return vals.count()
	// case Bool:
	// 	vals := s.Values.(boolValues)
	// 	return vals.count()
	case datatypes.String:
		vals := s.Values.(values.StringValues)
		return vals.Len()
	// case DateTime:
	// 	vals := s.Values.(dateTimeValues)
	// 	return vals.count()
	default:
		log.Printf("Count not supported for Series type %v", s.Kind)
		return 0
	}
}
