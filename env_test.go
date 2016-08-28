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
	env, err := newEnvironment(nc, metricsSubject)
	if !assert.NoError(t, err) {
		assert.FailNow(t, "can't create env")
	}

	// start listening for the metric
	msgs := make(chan *nats.Msg)
	sub, err := nc.Subscribe(metricsSubject, func(msg *nats.Msg) {
		msgs <- msg
		close(msgs)
	})
	defer sub.Unsubscribe()

	// create the metric
	m := env.newMetric("something", CounterType, nil)
	m.Value = 123
	err = m.send(nil)
	assert.Nil(t, err)

	select {
	case msg := <-msgs:
		m := new(metric)
		err = json.Unmarshal(msg.Data, m)
		assert.Nil(t, err)

		assert.Equal(t, "something", m.Name)
		assert.EqualValues(t, m.Value, 123)
		assert.Equal(t, m.Type, CounterType)
		assert.NotNil(t, m.Dims)
		assert.Len(t, m.Dims, 0)
	case <-time.After(time.Second):
		assert.FailNow(t, "didn't get the message in time")
	}

	// validate counts
}
