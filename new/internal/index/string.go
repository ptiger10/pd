package index

import "log"

type stringLabels []string

func (labels stringLabels) In(positions []int) interface{} {
	var ret stringLabels
	for _, position := range positions {
		if position >= len(labels) {
			log.Panicf("Unable to get index label(s): index out of range: %d", position)
		}
		ret = append(ret, labels[position])
	}
	return ret
}
