package index

import "log"

type IntLabels []int64

func (labels IntLabels) In(positions []int) interface{} {
	var ret IntLabels
	for _, position := range positions {
		if position >= len(labels) {
			log.Panicf("Unable to get index label(s): index out of range: %d", position)
		}
		ret = append(ret, labels[position])
	}
	return ret
}
