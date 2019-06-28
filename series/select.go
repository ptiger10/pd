package series

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/options"
)

// Element returns information about the value and index labels at this position but panics if an out-of-range position is provided.
func (s *Series) Element(position int) Element {
	elem := s.values.Element(position)
	idxElems := s.index.Elements(position)
	return Element{elem.Value, elem.Null, idxElems.Labels, idxElems.DataTypes}
}

// selectByRows copies a Series then subsets it to include only index items and values at the positions supplied
func (s *Series) selectByRows(positions []int) (*Series, error) {
	if err := s.ensureAlignment(); err != nil {
		return newEmptySeries(), fmt.Errorf("series internal alignment error: %v", err)
	}
	if err := s.ensureRowPositions(positions); err != nil {
		return newEmptySeries(), fmt.Errorf("s.selectByRows(): %v", err)
	}

	s = s.Copy()
	s.values = s.values.Subset(positions)
	s.index = s.index.Subset(positions)
	return s, nil
}

func (s *Series) mustSelectRows(positions []int) *Series {
	s, err := s.selectByRows(positions)
	if err != nil {
		log.Printf("Internal error: %v\n", err)
		return newEmptySeries()
	}
	return s
}

// Subset returns a subset of a Series based on the supplied integer positions.
func (s *Series) Subset(rowPositions []int) (*Series, error) {
	if rowPositions == nil {
		return newEmptySeries(), fmt.Errorf("series.Subset(): rowPositions cannot be nil")
	}
	if len(rowPositions) == 0 {
		return newEmptySeries(), fmt.Errorf("series.Subset(): no valid rows provided")
	}

	sub, err := s.selectByRows(rowPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("series.Subset(): %v", err)
	}
	return sub, nil
}

// At returns the value at a single integer position, bur returns nil if the position is out of range.
func (s *Series) At(position int) interface{} {
	if position >= s.Len() {
		if options.GetLogWarnings() {
			log.Printf("s.Index.At(): invalid position: %d (max: %v)", position, s.Len()-1)
		}
		return nil
	}
	elem := s.Element(position)
	return elem.Value
}

// [END Series methods]

// [START Selection]

// XS returns a new Series with only the rows and index levels at the specified positions.
func (s *Series) XS(rowPositions []int, levelPositions []int) (*Series, error) {
	var err error
	s, err = s.Subset(rowPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.XS() rows: %v", err)
	}
	s, err = s.Index.Subset(levelPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.XS() index levels: %v", err)
	}
	return s, nil
}

// SelectLabel returns the integer location of the first row in index level 0 with the supplied label, or -1 if the label does not exist.
func (s *Series) SelectLabel(label string) int {
	if s.index.NumLevels() == 0 {
		if options.GetLogWarnings() {
			log.Println("Series.SelectLabel(): index has no length")
		}
		return -1
	}
	val, ok := s.index.Levels[0].LabelMap[label]
	if !ok {
		if options.GetLogWarnings() {
			log.Printf("Series.SelectLabel(): %v not in label map\n", label)
		}
		return -1
	}
	return val[0]
}

// SelectLabels returns the integer locations of all rows with the supplied labels within the supplied level.
// If an error is encountered, returns a new slice of 0 length.
func (s *Series) SelectLabels(labels []string, level int) []int {
	empty := make([]int, 0)
	err := s.ensureLevelPositions([]int{level})
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("Series.SelectLabels(): %v", err)
		}
		return empty
	}
	include := make([]int, 0)
	for _, label := range labels {
		val, ok := s.index.Levels[level].LabelMap[label]
		if !ok {
			if options.GetLogWarnings() {
				log.Printf("Series.SelectLabels(): %v not in label map", label)
			}
			return empty
		}
		include = append(include, val...)
	}
	return include
}
