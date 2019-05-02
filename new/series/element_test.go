package series

import (
	"fmt"
	"testing"
)

func TestElement(t *testing.T) {
	s, err := New([]int{1, 2, 3}, Index([]string{"A", "B", "C"}), Index([]int{1, 2, 3}))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)
}
