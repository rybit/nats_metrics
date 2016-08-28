package metrics

import "sync"

// Gauge keeps a running measure of the value at that moment
type Gauge interface {
	Increment(*map[string]interface{}) error
	Decrement(*map[string]interface{}) error
	Set(int64, *map[string]interface{}) error
}

type gauge struct {
	metric
	valueLock sync.Mutex
}

func (e *environment) newGauge(name string, metricDims *map[string]interface{}) Gauge {
	m := e.newMetric(name, GaugeType, metricDims)
	return &gauge{
		metric:    *m,
		valueLock: sync.Mutex{},
	}
}

func (m *gauge) Increment(instanceDims *map[string]interface{}) error {
	m.valueLock.Lock()
	defer m.valueLock.Unlock()
	m.Value++
	return m.send(instanceDims)
}

func (m *gauge) Decrement(instanceDims *map[string]interface{}) error {
	m.valueLock.Lock()
	defer m.valueLock.Unlock()
	m.Value--
	return m.send(instanceDims)
}

func (m *gauge) Set(val int64, instanceDims *map[string]interface{}) error {
	m.valueLock.Lock()
	defer m.valueLock.Unlock()
	m.Value = val
	return m.send(instanceDims)
}
