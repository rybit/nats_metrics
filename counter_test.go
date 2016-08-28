package metrics

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats"
	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	msgs := make(chan *nats.Msg)
	sub, env := listenUntil(t, func(msg *nats.Msg) {
		msgs <- msg
	})
	defer sub.Unsubscribe()

	c := env.newCounter("thingy", nil)
	c.Count(nil)

	select {
	case msg := <-msgs:
		m := new(metric)
		err := json.Unmarshal(msg.Data, m)
		assert.Nil(t, err)
		assert.EqualValues(t, 1, m.Value)
	case <-time.After(time.Second):
		assert.FailNow(t, "failed to get message in time")
	}
}

func TestCountN(t *testing.T) {
	msgs := make(chan *nats.Msg)
	sub, env := listenUntil(t, func(msg *nats.Msg) {
		msgs <- msg
	})
	defer sub.Unsubscribe()

	c := env.newCounter("thingy", nil)
	c.CountN(100, nil)

	select {
	case msg := <-msgs:
		m := new(metric)
		err := json.Unmarshal(msg.Data, m)
		assert.Nil(t, err)
		assert.EqualValues(t, 100, m.Value)
	case <-time.After(time.Second):
		assert.FailNow(t, "failed to get message in time")
	}
}
