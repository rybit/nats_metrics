package metrics

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestDimensionalOverride(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	env.AddDimension("global-val", 12)
	env.AddDimension("metric-overide", "global-level")
	env.AddDimension("instance-overide", "global-level")
	sender := env.newMetric("thing.one", CounterType, &DimMap{
		"metric-val":       456,
		"metric-overide":   "metric-level",
		"instance-overide": "metric-level",
	})

	sender.send(&DimMap{
		"instance-overide": "instance-level",
		"instance-val":     789,
	})

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
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

		assert.NotEqual(t, time.Time{}, m.Timestamp)
	}
}
