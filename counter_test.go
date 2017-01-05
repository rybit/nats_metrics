package metrics

import (
	"testing"

	"time"

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

func TestCountMultipleTimes(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	ts := time.Now()

	c := env.newCounter("thingy", nil)
	c.SetTimestamp(ts)
	c.Count(nil)

	first := readOne(t, sub)
	if assert.NotNil(t, first) {
		assert.EqualValues(t, 1, first.Value)
		assert.Equal(t, ts.UnixNano(), first.Timestamp.UnixNano())
	}

	c.SetTimestamp(time.Time{})
	c.Count(nil)
	second := readOne(t, sub)
	if assert.NotNil(t, second) {
		assert.EqualValues(t, 1, second.Value)
		assert.True(t, ts.UnixNano() < second.Timestamp.UnixNano())
	}

	c.Count(nil)
	third := readOne(t, sub)
	if assert.NotNil(t, third) {
		assert.EqualValues(t, 1, third.Value)
		assert.True(t, second.Timestamp.UnixNano() < third.Timestamp.UnixNano())
	}
}
