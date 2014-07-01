package gousim

import (
	"fmt"
	"time"
)

type Analyzer struct {
	Metrics map[string] []int64
	StartTime time.Time
	EndTime time.Time
}

func NewAnalyzer(session *RecordingSession) *Analyzer {
	a := Analyzer{ nil,session.StartTime,session.EndTime}
	a.Metrics = make(map[string] []int64, len(session.dataMap))
	for k, v := range session.dataMap {
    	a.Metrics[k] = make([]int64,len(v))
    	copy(a.Metrics[k],v)
	}
	return &a
} 

type MetricNotFound string

func (m MetricNotFound) Error() string {
	return fmt.Sprintf("Cannot find metric %s.", m)
}

func (a *Analyzer) Average(metric string) (float64, error) {
	if a.Metrics[metric] == nil {
		return 0.0, MetricNotFound(metric)
	}
	total := int64(0)
    for _, x := range a.Metrics[metric] {
        total += x
    }
    return float64(total) / float64(len(a.Metrics)),nil
}

func (a *Analyzer) Median(metric string) (int64, error) {
	if a.Metrics[metric] == nil {
		return 0,MetricNotFound(metric)
	}
	return a.Metrics[metric][len(a.Metrics[metric])/2],nil
}

func (a *Analyzer) Max(metric string) (int64, error) {
	if a.Metrics[metric] == nil {
		return 0,MetricNotFound(metric)
	}
	max := int64(-9223372036854775808)
	for _,v := range a.Metrics[metric] {
		if v > max {
			max = v
		}
	}
	return max,nil
}

func (a *Analyzer) Min(metric string) (int64, error) {
	if a.Metrics[metric] == nil {
		return 0,MetricNotFound(metric)
	}	
	min := int64(9223372036854775807)
	for _,v := range a.Metrics[metric] {
		if v < min {
			min = v
		}
	}
	return min, nil
}

func (a *Analyzer) Count(metric string) (int, error) {
	if a.Metrics[metric] == nil {
		return 0,MetricNotFound(metric)
	}
	return len(a.Metrics[metric]),nil		
}