package metrics

import "time"

func NewStats() Stats {
	return Stats{
		Time:       time.Now().UnixNano(),
		Counters:   map[string]int64{},
		Gauges:     map[string]int64{},
		Histograms: map[string]HistogramData{},
	}
}

type Stats struct {
	Time       int64
	Counters   map[string]int64
	Gauges     map[string]int64
	Histograms map[string]HistogramData
}

// HistogramData is a snapshot of an actual histogram
type HistogramData struct {
	Max         int64
	Min         int64
	Mean        float64
	Stddev      float64
	Variance    float64
	Percentiles []float64
}
