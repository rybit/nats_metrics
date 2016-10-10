package metrics

import "sync"

const (
	// TimerType is the type for timers
	TimerType MetricType = "timer"
	// CounterType is the type for timers
	CounterType MetricType = "counter"
	// GaugeType is the type for gauges
	GaugeType MetricType = "gauge"
)

// MetricType describes what type the metric is
type MetricType string

type metric struct {
	Name  string     `json:"name"`
	Type  MetricType `json:"type"`
	Value int64      `json:"value"`
	Dims  DimMap     `json:"dimensions"`

	dimlock sync.Mutex
	env     *environment
}

// AddDimension will add this dimension with locking
func (m *metric) AddDimension(key string, value interface{}) *metric {
	m.dimlock.Lock()
	defer m.dimlock.Unlock()
	m.Dims[key] = value
	return m
}

func (m *metric) send(instanceDims *DimMap) error {
	if m.env == nil {
		return InitError{"Environment not initialized"}
	}
	return m.env.send(m, instanceDims)
}

// DimMap is a map of dimensions
type DimMap map[string]interface{}
