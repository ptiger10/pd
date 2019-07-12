// +build benchmarks

package benchmarks

// func TestPythonProfiler(t *testing.T) {
// 	got := RunPythonProfiler()
// 	fmt.Println(got)
// }

// func TestCompareBenchmarks(t *testing.T) {
// 	type args struct {
// 		goBenchmarks Results
// 		pyBenchmarks Results
// 		sampleSizes  []string
// 		descs        map[string]desc
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{name: "normal", args: args{
// 			goBenchmarks: Results{"100k": {
// 				"sum": []interface{}{"50ms", 50.0}, "mean": []interface{}{"50ms", 50.0}}},
// 			pyBenchmarks: Results{"100k": {"sum": []interface{}{"100ms", 100.0}}},
// 			descs:        map[string]desc{"sum": desc{1, "Simple sum"}, "mean": desc{2, "Simple mean"}},
// 			sampleSizes:  []string{"100k", "200k"}}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := CompareBenchmarks(tt.args.goBenchmarks, tt.args.pyBenchmarks, tt.args.sampleSizes, tt.args.descs)
// 			print(got)
// 		})
// 	}
// }

// func TestProfileGo(t *testing.T) {
// 	ProfileGo(benchmarkMeanFloat64_100000)

// }
