package metrics

import (
	"sync"
	"time"
)

const (
	// TimerType is the type for timers
	TimerType MetricType = "timer"
	// CounterType is the type for timers
	CounterType MetricType = "counter"
	// GaugeType is the type for gauges
	GaugeType MetricType = "gauge"
)

var zeroTime time.Time

// MetricType describes what type the metric is
type MetricType string

type RawMetric struct {
	Name      string     `json:"name"`
	Type      MetricType `json:"type"`
	Value     int64      `json:"value"`
	Dims      DimMap     `json:"dimensions"`
	Timestamp time.Time  `json:"timestamp"`
}

type metric struct {
	RawMetric

	dimlock sync.Mutex
	env     *environment
}

func (m *metric) SetTimestamp(t time.Time) {
	m.Timestamp = t
}

// AddDimension will add this dimension with locking
func (m *metric) AddDimension(key string, value interface{}) *metric {
	m.dimlock.Lock()
	defer m.dimlock.Unlock()
	m.Dims[key] = value
	return m
}

func (m *metric) send(instanceDims DimMap) error {
	if m.env == nil {
		return InitError{errString{"Environment not initialized"}}
	}
	metricToSend := &RawMetric{
		Type:      m.Type,
		Value:     m.Value,
		Name:      m.Name,
		Timestamp: m.Timestamp,
		Dims:      DimMap{},
	}

	// global
	m.env.dimlock.Lock()
	addAll(metricToSend.Dims, m.env.globalDims)
	m.env.dimlock.Unlock()

	// metric
	m.dimlock.Lock()
	addAll(metricToSend.Dims, m.Dims)
	m.dimlock.Unlock()

	// instance
	addAll(metricToSend.Dims, instanceDims)

	if metricToSend.Timestamp == zeroTime {
		metricToSend.Timestamp = time.Now()
	}

	return m.env.send(metricToSend)
}

// DimMap is a map of dimensions
type DimMap map[string]interface{}
