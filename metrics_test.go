package metrics

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats"
	"github.com/stretchr/testify/assert"
)

func TestDimensionalOverride(t *testing.T) {
	msgs := make(chan *nats.Msg)
	sub, env := listenUntil(t, func(msg *nats.Msg) {
		msgs <- msg
	})
	defer sub.Unsubscribe()

	env.AddDimension("global-val", 12)
	env.AddDimension("metric-overide", "global-level")
	env.AddDimension("instance-overide", "global-level")
	m := env.newMetric("thing.one", CounterType, &map[string]interface{}{
		"metric-val":       456,
		"metric-overide":   "metric-level",
		"instance-overide": "metric-level",
	})

	m.send(&map[string]interface{}{
		"instance-overide": "instance-level",
		"instance-val":     789,
	})
	select {
	case msg := <-msgs:
		m := new(metric)
		err := json.Unmarshal(msg.Data, m)
		assert.Nil(t, err)

		// check the dimensions
		assert.EqualValues(t, 5, len(m.Dims))
		expected := map[string]interface{}{
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
	case <-time.After(time.Second):
		assert.FailNow(t, "failed to get messages in time")
	}
}
