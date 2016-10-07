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
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	// create the metric
	sender := env.newMetric("something", CounterType, nil)
	sender.Value = 123
	err := sender.send(nil)
	assert.Nil(t, err)

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
		assert.Equal(t, "something", m.Name)
		assert.EqualValues(t, m.Value, 123)
		assert.Equal(t, m.Type, CounterType)
		assert.NotNil(t, m.Dims)
		assert.Len(t, m.Dims, 0)
	}

	// validate counts
	checkCounters(t, 1, 0, 0, env)
}

func subscribe(t *testing.T) (*nats.Subscription, *environment) {
	sub, err := nc.SubscribeSync(metricsSubject)
	if err != nil {
		assert.FailNow(t, "Failed to subscribe: "+err.Error())
	}

	env, err := newEnvironment(nc, metricsSubject)
	if err != nil {
		assert.FailNow(t, "Failed to create test env")
	}

	return sub, env
}

func checkCounters(t *testing.T, counters, timers, gauges int, env *environment) {
	assert.EqualValues(t, counters, env.countersSent)
	assert.EqualValues(t, timers, env.timersSent)
	assert.EqualValues(t, gauges, env.gaugesSent)
}

func readOne(t *testing.T, sub *nats.Subscription) *metric {
	msg, err := sub.NextMsg(time.Second)
	if err != nil {
		assert.Fail(t, "Failed waiting for a message: "+err.Error())
		return nil
	}
	m := new(metric)
	err = json.Unmarshal(msg.Data, m)
	assert.Nil(t, err)
	return m
}
