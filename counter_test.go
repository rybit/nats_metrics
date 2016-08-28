package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	sub, env, msgs := listenForOne(t)
	defer sub.Unsubscribe()

	c := env.newCounter("thingy", nil)
	c.Count(nil)

	thisOrTimeout(t, msgs, func(m *metric) {
		assert.EqualValues(t, 1, m.Value)
	})
}

func TestCountN(t *testing.T) {
	sub, env, msgs := listenForOne(t)
	defer sub.Unsubscribe()

	c := env.newCounter("thingy", nil)
	c.CountN(100, nil)

	thisOrTimeout(t, msgs, func(m *metric) {
		assert.EqualValues(t, 100, m.Value)
	})
}
