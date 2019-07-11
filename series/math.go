package series

import (
	"math"
	"runtime"
	"sort"
	"sync"

	"github.com/montanaflynn/stats" // and stats package otherwise
	"github.com/ptiger10/pd/options"
)

// an interface of valid (non-null) values; appropriate for type assertion
func (s *Series) validVals() interface{} {
	valid := s.values.Subset(s.valid())
	return valid.Vals()
}

func asyncMath(data []float64, awaitFn func([]float64, chan<- interface{}, *sync.WaitGroup)) chan interface{} {
	var wg sync.WaitGroup
	numPartitions := runtime.GOMAXPROCS(0)
	// when partition is applied, residual values go in last partition
	// 8 values, 3 partitions -> [2 values], [2 values], [4 values]
	valsPerPartition := len(data) / numPartitions

	ch := make(chan interface{}, numPartitions)
	for i := 0; i < numPartitions; i++ {
		var sub []float64
		if i != numPartitions-1 {
			sub = data[i*valsPerPartition : (i+1)*(valsPerPartition)]
		} else {
			sub = data[i*valsPerPartition:] // residual values go in last partition
		}
		wg.Add(1)
		go awaitFn(sub, ch, &wg)
	}
	wg.Wait()
	close(ch)
	return ch
}

// Sum of non-null float64 or int64 Series values. For bool values, sum of true values. If inapplicable, defaults to math.Nan().
func (s *Series) Sum() float64 {
	var sum float64

	// null int values are represented as 0, but that's ok for sum
	sumFunc := func(data []float64) float64 {
		var sum float64
		for _, d := range data {
			if !math.IsNaN(d) {
				sum += d
			}
		}
		return sum
	}

	awaitSumFunc := func(data []float64, ch chan<- interface{}, wg *sync.WaitGroup) {
		ch <- sumFunc(data)
		wg.Done()
	}

	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(s.Vals())

		if !options.GetAsync() {
			return sumFunc(data)
		}
		ch := asyncMath(data, awaitSumFunc)
		for partitionSum := range ch {
			sum += partitionSum.(float64)
		}
		return sum

	case options.Bool:
		data := ensureBools(s.Vals())
		// null bool values are represented as false, but that's ok for sum
		for _, d := range data {
			if d {
				sum++
			}
		}
		return sum
	default:
		return math.NaN()
	}
}

// Mean of non-null series values. For bool values, mean of true values.
// Applies to float64 and int64. If inapplicable, defaults to math.Nan().
func (s *Series) Mean() float64 {
	meanFunc := func(data []float64) (sum float64, counter int) {
		for _, d := range data {
			if !math.IsNaN(d) {
				sum += d
				counter++
			}
		}
		return
	}
	awaitMeanFunc := func(data []float64, ch chan<- interface{}, wg *sync.WaitGroup) {
		sum, counter := meanFunc(data)
		ch <- []float64{sum, float64(counter)}
		wg.Done()
	}

	switch s.datatype {
	case options.Float64:
		data := ensureFloatFromNumerics(s.Vals())
		if !options.GetAsync() {
			sum, validCount := meanFunc(data)
			return sum / float64(validCount)
		}

		var sum float64
		var validCount float64
		ch := asyncMath(data, awaitMeanFunc)
		for partitionMean := range ch {
			// partition form: [sum, validCount]
			p := partitionMean.([]float64)
			sum += p[0]
			validCount += p[1]
		}
		return sum / validCount

	case options.Int64:
		return s.Sum() / float64(s.validCount())
	case options.Bool:
		return s.Sum() / float64(s.validCount())
	default:
		return math.NaN()
	}
}

// Median of a series. Applies to float64 and int64. If inapplicable, defaults to math.Nan().
func (s *Series) Median() float64 {
	calculateMedian := func(data []float64) float64 {
		if len(data) == 0 {
			return math.NaN()
		}
		// rounds down if there are even number of elements
		mNumber := len(data) / 2

		// odd number of elements
		if len(data)%2 != 0 {
			return data[mNumber]
		}
		// even number of elements
		return (data[mNumber-1] + data[mNumber]) / 2
	}

	switch s.datatype {
	case options.Float64:
		// sort then slice out the nulls from the beginning
		data := ensureFloatFromNumerics(s.Vals())
		sort.Float64s(data)
		numValids := s.validCount()
		firstValid := len(data) - numValids
		data = data[firstValid:]
		return calculateMedian(data)

	case options.Int64:
		vals := s.validVals()
		data := ensureFloatFromNumerics(vals)
		sort.Float64s(data)
		return calculateMedian(data)

	default:
		return math.NaN()
	}
}

// Min of a series. Applies to float64 and int64. If inapplicable, defaults to math.Nan().
func (s *Series) Min() float64 {
	minFunc := func(data []float64) float64 {
		min := math.NaN()
		for _, d := range data {
			if !math.IsNaN(d) {
				if math.IsNaN(min) {
					min = d
				} else if d < min {
					min = d
				}
			}
		}
		return min
	}
	awaitMinFunc := func(data []float64, ch chan<- interface{}, wg *sync.WaitGroup) {
		ch <- minFunc(data)
		wg.Done()
	}
	min := math.NaN()
	switch s.datatype {
	case options.Float64:
		data := ensureFloatFromNumerics(s.Vals())
		if !options.GetAsync() {
			return minFunc(data)
		}
		ch := asyncMath(data, awaitMinFunc)
		var results []float64
		for partitionMin := range ch {
			results = append(results, partitionMin.(float64))
		}
		return minFunc(results)

	case options.Int64:
		data := ensureFloatFromNumerics(s.Vals())
		for i, d := range data {
			if !s.values.Null(i) {
				if math.IsNaN(min) {
					min = d
				} else if d < min {
					min = d
				}
			}
		}
		return min
	default:
		return math.NaN()
	}
}

// Max of a series. Applies to float64 and int64. If inapplicable, defaults to math.Nan().
func (s *Series) Max() float64 {
	maxFunc := func(data []float64) float64 {
		max := math.NaN()
		for _, d := range data {
			if !math.IsNaN(d) {
				if math.IsNaN(max) {
					max = d
				} else if d > max {
					max = d
				}
			}
		}
		return max
	}
	awaitMaxFunc := func(data []float64, ch chan<- interface{}, wg *sync.WaitGroup) {
		ch <- maxFunc(data)
		wg.Done()
	}
	max := math.NaN()
	switch s.datatype {
	case options.Float64:
		data := ensureFloatFromNumerics(s.Vals())
		if !options.GetAsync() {
			return maxFunc(data)
		}
		ch := asyncMath(data, awaitMaxFunc)
		var results []float64
		for partitionMax := range ch {
			results = append(results, partitionMax.(float64))
		}
		return maxFunc(results)

	case options.Int64:
		data := ensureFloatFromNumerics(s.Vals())
		for i, d := range data {
			if !s.values.Null(i) {
				if math.IsNaN(max) {
					max = d
				} else if d > max {
					max = d
				}
			}
		}
		return max
	default:
		return math.NaN()
	}
}

// Quartile i (must be 1, 2, 3)
//
// Applies to: Float, Int. If inapplicable, defaults to math.Nan().
func (s *Series) Quartile(i int) float64 {
	vals := s.validVals()
	switch s.datatype {
	case options.Float64, options.Int64:
		data := ensureFloatFromNumerics(vals)
		val, err := stats.Quartile(data)
		if err != nil {
			return math.NaN()
		}
		switch i {
		case 1:
			return val.Q1
		case 2:
			return val.Q2
		case 3:
			return val.Q3
		default:
			return math.NaN()
		}
	default:
		return math.NaN()
	}
}

// Std returns the Standard Deviation of a series.
//
// Applies to: Float, Int. If inapplicable, defaults to math.Nan().
func (s *Series) Std() float64 {
	stdFunc := func(data []float64) (variance float64, counter int) {
		mean := s.Mean()
		for _, d := range data {
			if !math.IsNaN(d) {
				variance += (d - mean) * (d - mean)
				counter++
			}
		}
		return
	}
	awaitStdFunc := func(data []float64, ch chan<- interface{}, wg *sync.WaitGroup) {
		sum, counter := stdFunc(data)
		ch <- []float64{sum, float64(counter)}
		wg.Done()
	}
	switch s.datatype {
	case options.Float64:
		data := ensureFloatFromNumerics(s.Vals())
		if !options.GetAsync() {
			variance, validCount := stdFunc(data)
			return math.Pow(variance/float64(validCount), 0.5)
		}
		var variance float64
		var validCount float64
		ch := asyncMath(data, awaitStdFunc)
		for partitionMean := range ch {
			// partition form: [sum, validCount]
			p := partitionMean.([]float64)
			variance += p[0]
			validCount += p[1]
		}
		return math.Pow(variance/validCount, 0.5)

	case options.Int64:
		data := ensureFloatFromNumerics(s.Vals())
		mean := s.Mean()
		var variance float64
		var counter int
		for i, d := range data {
			if !s.values.Null(i) {
				variance += (d - mean) * (d - mean)
				counter++
			}
		}
		return math.Pow(variance/float64(counter), 0.5)
	default:
		return math.NaN()
	}
}
