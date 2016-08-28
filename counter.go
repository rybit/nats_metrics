package metrics

// NewCounter creates a named counter with these dimensions
func NewCounter(name string, metricDims *map[string]interface{}) (Counter, error) {
	if err := checkEnv(); err != nil {
		return nil, err
	}

	return globalEnv.newCounter(name, metricDims), nil
}

// Counter will send when an event occurs
type Counter interface {
	Count(*map[string]interface{}) error
	CountN(int64, *map[string]interface{}) error
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
