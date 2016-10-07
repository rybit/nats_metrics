package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	c := env.newCounter("thingy", nil)
	c.Count(nil)

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
		assert.EqualValues(t, 1, m.Value)
	}
}

func TestCountN(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	c := env.newCounter("thingy", nil)
	c.CountN(100, nil)

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
		assert.EqualValues(t, 100, m.Value)
	}
}
