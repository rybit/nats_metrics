package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrement(t *testing.T) {
	sub, env, msgs := listenForOne(t)
	defer sub.Unsubscribe()

	g := env.newGauge("something", nil)
	g.Increment(nil)

	thisOrTimeout(t, msgs, func(m *metric) {
		assert.EqualValues(t, 1, m.Value)
	})

}

func TestDecrement(t *testing.T) {
	sub, env, msgs := listenForOne(t)
	defer sub.Unsubscribe()

	g := env.newGauge("something", nil)
	g.Decrement(nil)
	thisOrTimeout(t, msgs, func(m *metric) {
		assert.EqualValues(t, -1, m.Value)
	})
}

func TestSet(t *testing.T) {
	sub, env, msgs := listenForOne(t)
	defer sub.Unsubscribe()

	g := env.newGauge("something", nil)
	g.Set(123, nil)
	thisOrTimeout(t, msgs, func(m *metric) {
		assert.EqualValues(t, 123, m.Value)
	})
}
