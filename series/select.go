package series

import (
	"fmt"
	"log"

	"github.com/ptiger10/pd/internal/values"

	"github.com/ptiger10/pd/options"
)

// Element returns information about the value and index labels at this position but panics if an out-of-range position is provided.
func (s *Series) Element(position int) Element {
	elem := s.values.Element(position)
	idxElems := s.index.Elements(position)
	return Element{elem.Value, elem.Null, idxElems.Labels, idxElems.DataTypes}
}

// subsetRows subsets a Series to include only index items and values at the row positions supplied and modifies the Series in place.
func (ip InPlace) subsetRows(positions []int) {
	ip.s.values = ip.s.values.Subset(positions)
	ip.s.index.Subset(positions)
}

// subsetRows subsets a Series to include only index items and values at the row positions supplied and returns a new Series.
func (s *Series) subsetRows(positions []int) *Series {
	s = s.Copy()
	s.InPlace.subsetRows(positions)
	return s
}

// Subset subsets a Series in place.
func (ip InPlace) Subset(rowPositions []int) error {
	if err := ip.s.ensureAlignment(); err != nil {
		return fmt.Errorf("series.Subset: internal alignment error: %v", err)
	}
	if err := ip.s.ensureRowPositions(rowPositions); err != nil {
		return fmt.Errorf("series.Subset(): %v", err)
	}

	ip.s.InPlace.subsetRows(rowPositions)
	return nil
}

// Subset subsets a Series based on the supplied integer positions and returns a new Series.
func (s *Series) Subset(rowPositions []int) (*Series, error) {
	s = s.Copy()
	err := s.InPlace.Subset(rowPositions)
	if err != nil {
		return newEmptySeries(), err
	}
	return s, nil
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

// From subsets the Series from start to end (inclusive) and returns a new Series.
// If an invalid position is provided, returns empty Series.
func (s *Series) From(start int, end int) *Series {
	rowPositions := values.MakeIntRangeInclusive(start, end)
	if err := s.ensureRowPositions(rowPositions); err != nil {
		if options.GetLogWarnings() {
			log.Printf("s.From(): %v", err)
		}
		return newEmptySeries()
	}
	return s.subsetRows(rowPositions)
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
	err = s.Index.SubsetLevels(levelPositions)
	if err != nil {
		return newEmptySeries(), fmt.Errorf("s.XS() index levels: %v", err)
	}
	return s, nil
}

// SelectLabel returns the integer location of the first row in index level 0 with the supplied label, or -1 if the label does not exist.
func (s *Series) SelectLabel(label string) int {
	if s.NumLevels() == 0 {
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
