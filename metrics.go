package metrics

import "sync"

const (
	TimerType   MetricType = "timer"
	CounterType MetricType = "counter"
	GaugeType   MetricType = "gauge"
)

// MetricType describes what type the metric is
type MetricType string

type metric struct {
	Name  string                 `json:"name"`
	Type  MetricType             `json:"type"`
	Value int64                  `json:"value"`
	Dims  map[string]interface{} `json:"dimensions"`

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

func (m *metric) send(instanceDims *map[string]interface{}) error {
	return m.env.send(m, instanceDims)
}
