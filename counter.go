package metrics

// Counter will send when an event occurs
type Counter interface {
	Count(*map[string]interface{}) error
	CountN(int64, *map[string]interface{}) error
}

func (e *environment) newCounter(name string, metricDims *map[string]interface{}) Counter {
	return e.newMetric(name, CounterType, metricDims)
}

// Count will count 1 occurrence of an event
func (m *metric) Count(instanceDims *map[string]interface{}) error {
	return m.CountN(1, instanceDims)
}

//CountN will count N occurrences of an event
func (m *metric) CountN(val int64, instanceDims *map[string]interface{}) error {
	m.Value = val
	return m.send(instanceDims)
}
