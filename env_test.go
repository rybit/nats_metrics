package metrics

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nats-io/gnatsd/server"
	"github.com/nats-io/nats"
	"github.com/nats-io/nats/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var nc *nats.Conn
var s *server.Server

const metricsSubject = "test.metrics"

func TestMain(m *testing.M) {

	s = test.RunDefaultServer()
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

func TestSendMetricWithNilConn(t *testing.T) {
	env, err := NewEnvironment(nil, metricsSubject)
	if assert.NoError(t, err) {
		sender := env.newMetric("something", CounterType, nil)
		sender.Value = 123
		err := sender.send(nil)
		assert.Nil(t, err)
	}
}

func TestSendWithTracer(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	called := false
	env.tracer = func(m *RawMetric) {
		if assert.NotNil(t, m) {
			assert.Equal(t, "something", m.Name)
			assert.EqualValues(t, m.Value, 123)
			assert.Equal(t, m.Type, CounterType)
			assert.NotNil(t, m.Dims)
			assert.Len(t, m.Dims, 0)
		}
		called = true
	}

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
	assert.True(t, called)
	// validate counts
	checkCounters(t, 1, 0, 0, env)
}

func TestSeparateEnv(t *testing.T) {
	f1, err := NewEnvironment(nc, "first-env")
	require.NoError(t, err)
	f2, err := NewEnvironment(nc, "second-env")
	require.NoError(t, err)

	sub1, err := nc.SubscribeSync("first-env")
	require.NoError(t, err)
	sub2, err := nc.SubscribeSync("second-env")
	require.NoError(t, err)

	require.NoError(t, f1.NewCounter("c1", nil).Count(nil))
	require.NoError(t, f2.NewCounter("c2", nil).Count(nil))

	raw1, err := sub1.NextMsg(time.Second)
	require.NoError(t, err)
	raw2, err := sub2.NextMsg(time.Second)
	require.NoError(t, err)

	m1 := new(RawMetric)
	m2 := new(RawMetric)
	require.NoError(t, json.Unmarshal(raw1.Data, m1))
	require.NoError(t, json.Unmarshal(raw2.Data, m2))

	assert.Equal(t, "c1", m1.Name)
	assert.Equal(t, "c2", m2.Name)
}

func subscribe(t *testing.T) (*nats.Subscription, *Environment) {
	sub, err := nc.SubscribeSync(metricsSubject)
	if err != nil {
		assert.FailNow(t, "Failed to subscribe: "+err.Error())
	}

	env, err := NewEnvironment(nc, metricsSubject)
	if err != nil {
		assert.FailNow(t, "Failed to create test env")
	}

	return sub, env
}

func checkCounters(t *testing.T, counters, timers, gauges int, env *Environment) {
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
