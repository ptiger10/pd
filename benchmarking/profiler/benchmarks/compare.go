// +build benchmarks

package benchmarks

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// CompareBenchmarks creates a comparison table of GoPandas <> Pandas for equivalent operations
func CompareBenchmarks(
	goBenchmarks, pyBenchmarks Results,
	sampleSizes []string,
	descs map[string]desc,
) string {

	var printer string
	printer += "GoPandas vs Pandas speed comparison\n"
	printer += time.Now().In(time.Local).Format(time.RFC1123) + "\n"
	// model
	// +-----+-----+
	// | foo | bar |
	// +-----+-----+
	spacerChar := "-"
	sepChar := "+"
	vChar := "|"

	// Sections
	type section struct {
		name  string
		width int
	}
	num := section{name: "#", width: 4}
	desc := section{name: "DESCRIPTION", width: 40}
	sample := section{name: "N", width: 6}

	goBenchmark := section{name: "GOPANDAS", width: 11}
	pyBenchmark := section{name: "PANDAS", width: 11}
	comparison := section{name: "SPEED Δ", width: 9}
	sections := []section{num, desc, sample, goBenchmark, pyBenchmark, comparison}

	// Break Line
	breakLineComponents := make([]string, len(sections))
	for i := 0; i < len(sections); i++ {
		breakLineComponents[i] = strings.Repeat(spacerChar, sections[i].width)
	}
	breakLine := sepChar + strings.Join(breakLineComponents, sepChar) + sepChar + "\n"

	// Header
	headerComponents := make([]string, len(sections))
	for i := 0; i < len(sections); i++ {
		headerComponents[i] = fmt.Sprintf(" %-*v", sections[i].width-1, sections[i].name)
	}
	header := vChar + strings.Join(headerComponents, vChar) + vChar + "\n"
	printer += breakLine + header + breakLine

	// Rows
	var i int
	type orderedDesc struct {
		n     int
		label string
	}
	var orderedDescs []orderedDesc
	for k, v := range descs {
		orderedDescs = append(orderedDescs, orderedDesc{v.order, k})
	}
	sort.Slice(orderedDescs, func(i, j int) bool {
		if orderedDescs[i].n < orderedDescs[j].n {
			return true
		}
		return false
	})
	for _, sample := range sampleSizes {
		results, ok := goBenchmarks[sample]
		if !ok {
			log.Printf("sample size %v not in %v", sample, goBenchmarks)
			continue
		}
		for _, desc := range orderedDescs {
			testName := desc.label
			goResult, ok := results[desc.label]
			if !ok {
				continue
			}
			i++
			gospeed, gons := goResult[0], goResult[1]
			goSpeed := gospeed.(string)
			goNS := gons.(float64)
			pySpeed := "n/a"
			comparison := "n/a"
			py, ok := pyBenchmarks[sample]
			if ok {
				pyResult, ok := py[testName]
				if ok {
					pyspeed, pyns := pyResult[0], pyResult[1]
					pySpeed = pyspeed.(string)
					pyNS := pyns.(float64)
					comparison = fmt.Sprintf("%.2fx", pyNS/goNS)
				}
			}

			rowComponents := []string{
				strconv.Itoa(i), descs[desc.label].str, sample, goSpeed, pySpeed, comparison,
			}
			for i := range rowComponents {
				rowComponents[i] = fmt.Sprintf(
					" %-*v", sections[i].width-1, rowComponents[i])
			}
			printer += vChar + strings.Join(rowComponents, vChar) + vChar + "\n"
			printer += breakLine

		}
	}
	return printer
}
