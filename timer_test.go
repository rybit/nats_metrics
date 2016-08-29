package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeIt(t *testing.T) {
	sub, env, msgs := listenForOne(t)
	defer sub.Unsubscribe()

	timer := env.newTimer("something", nil)
	start := timer.Start()
	<-time.After(time.Millisecond * 100)
	stop := time.Now()
	_, err := timer.Stop(nil)
	assert.Nil(t, err)

	thisOrTimeout(t, msgs, func(m *metric) {
		measured := start.Add(time.Duration(m.Value))
		assert.WithinDuration(t, stop, measured, time.Millisecond*10)
	})
}
