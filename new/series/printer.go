package series

import (
	"fmt"
	"reflect"
	"strings"
)

func (s Series) String() string {
	switch s.Kind {
	// case DateTime:
	// 	var printer string
	// 	vals := s.Values.(dateTimeValues)
	// 	for _, val := range vals {
	// 		printer += fmt.Sprintln(val.v.Format("01/02/2006"))
	// 	}
	// 	return printer
	default:
		return print(s)
	}
}

// expects to receive a slice of typed value structs (eg pd.floatValues)
// each struct must contain a boolean field called "v"
func print(s Series) string {

	vals := reflect.ValueOf(s.Values)
	idx := reflect.ValueOf(s.Index.Levels)
	var printer string
	for i := 0; i < s.Len(); i++ {
		var newLine string
		for j := 0; j < idx.Len(); j++ {
			// s.Index.Levels[j].Labels[i]
			newLine += fmt.Sprintf("%v ", idx.Index(j).FieldByName("Labels").Elem().Index(i))
		}
		// s.Values[i].V
		newLine += fmt.Sprint(vals.Index(i).FieldByName("V"))
		// Null strings must not return a trailing space
		newLine = strings.TrimSpace(newLine)
		printer += fmt.Sprintln(newLine)
	}
	return printer
}
