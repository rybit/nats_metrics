package metrics

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats"
	"github.com/nats-io/nats/test"
	"github.com/stretchr/testify/assert"
)

var nc *nats.Conn

const metricsSubject = "test.metrics"

func TestMain(m *testing.M) {

	s := test.RunDefaultServer()
	defer s.Shutdown()

	var err error
	nc, err = nats.Connect("nats://" + s.Addr().String())
	if err != nil {
		log.Fatal("failed to connect to server: " + err.Error())
	}
	defer nc.Close()

	os.Exit(m.Run())
}

func TestSendMetric(t *testing.T) {
	// start listening for the metric
	sub, env, msgs := listenForOne(t)
	defer sub.Unsubscribe()

	// create the metric
	m := env.newMetric("something", CounterType, nil)
	m.Value = 123
	err := m.send(nil)
	assert.Nil(t, err)

	thisOrTimeout(t, msgs, func(m *metric) {
		assert.Equal(t, "something", m.Name)
		assert.EqualValues(t, m.Value, 123)
		assert.Equal(t, m.Type, CounterType)
		assert.NotNil(t, m.Dims)
		assert.Len(t, m.Dims, 0)
	})

	// validate counts
	checkCounters(t, 1, 0, 0, env)
}

func listenForOne(t *testing.T) (*nats.Subscription, *environment, chan *nats.Msg) {
	msgs := make(chan *nats.Msg)
	sub, err := nc.Subscribe(metricsSubject, func(msg *nats.Msg) {
		msgs <- msg
		close(msgs)
	})
	if err != nil {
		assert.FailNow(t, "Failed to subscribe")
	}

	env, err := newEnvironment(nc, metricsSubject)
	if err != nil {
		assert.FailNow(t, "Failed to create test env")
	}

	return sub, env, msgs
}

func checkCounters(t *testing.T, counters, timers, gauges int, env *environment) {
	assert.EqualValues(t, counters, env.countersSent)
	assert.EqualValues(t, timers, env.timersSent)
	assert.EqualValues(t, gauges, env.gaugesSent)
}

func thisOrTimeout(t *testing.T, msgs chan *nats.Msg, f func(m *metric)) {
	select {
	case msg := <-msgs:
		m := new(metric)
		err := json.Unmarshal(msg.Data, m)
		assert.Nil(t, err)
	case <-time.After(time.Second):
		assert.FailNow(t, "didn't get the message in time")
	}
}
