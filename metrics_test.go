package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDimensionalOverride(t *testing.T) {
	sub, env, msgs := listenForOne(t)
	defer sub.Unsubscribe()

	env.AddDimension("global-val", 12)
	env.AddDimension("metric-overide", "global-level")
	env.AddDimension("instance-overide", "global-level")
	m := env.newMetric("thing.one", CounterType, &DimMap{
		"metric-val":       456,
		"metric-overide":   "metric-level",
		"instance-overide": "metric-level",
	})

	m.send(&DimMap{
		"instance-overide": "instance-level",
		"instance-val":     789,
	})
	thisOrTimeout(t, msgs, func(m *metric) {
		assert.EqualValues(t, 5, len(m.Dims))
		expected := DimMap{
			"global-val":       12,
			"metric-overide":   "metric-level",
			"instance-overide": "instance-level",
			"metric-val":       456,
			"instance-val":     789,
		}
		for k, v := range expected {
			dimVal, exists := m.Dims[k]
			assert.True(t, exists)
			assert.EqualValues(t, v, dimVal)
		}
	})
}
