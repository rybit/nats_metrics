package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrement(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	g := env.newGauge("something", nil)
	g.Increment(nil)

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
		assert.EqualValues(t, 1, m.Value)
	}
}

func TestDecrement(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	g := env.newGauge("something", nil)
	g.Decrement(nil)

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
		assert.EqualValues(t, -1, m.Value)
	}
}

func TestSet(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	g := env.newGauge("something", nil)
	g.Set(123, nil)

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
		assert.EqualValues(t, 123, m.Value)
	}
}
