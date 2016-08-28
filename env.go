package metrics

import (
	"sync"

	"github.com/nats-io/nats"
)

func newEnvironment(nc *nats.Conn, subject string) (*environment, error) {
	econn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	env := &environment{
		subject:    subject,
		nc:         econn,
		dimlock:    sync.Mutex{},
		globalDims: make(map[string]interface{}),
	}
	if err := env.isReady(); err != nil {
		return nil, err
	}

	return env, nil
}

type environment struct {
	subject    string
	globalDims map[string]interface{}
	dimlock    sync.Mutex
	nc         *nats.EncodedConn

	// some metrics stuff
	timersSent   int64
	countersSent int64
	gaugesSent   int64
}

func (e *environment) send(m *metric, instanceDims *map[string]interface{}) error {
	if err := e.isReady(); err != nil {
		return err
	}

	// copy it so we don't mess it up
	metricToSend := metric{
		Type:  m.Type,
		Value: m.Value,
		Name:  m.Name,
		Dims:  make(map[string]interface{}),
	}

	// global
	e.dimlock.Lock()
	addAll(&metricToSend.Dims, &e.globalDims)
	e.dimlock.Unlock()

	// metric
	m.dimlock.Lock()
	addAll(&metricToSend.Dims, &m.Dims)
	m.dimlock.Unlock()

	// instance
	addAll(&metricToSend.Dims, instanceDims)

	// TODO count it

	return e.nc.Publish(e.subject, &metricToSend)
}

func (e *environment) isReady() error {
	if e.nc == nil {
		return InitError{errString{"Nil nats connection provided"}}
	}
	if e.subject == "" {
		return InitError{errString{"No subject provided"}}
	}
	return nil
}

func addAll(into *map[string]interface{}, from *map[string]interface{}) {
	if into != nil {
		if from != nil {
			for k, v := range *from {
				(*into)[k] = v
			}
		}
	}
}

func (e *environment) newMetric(name string, t MetricType, dims *map[string]interface{}) *metric {
	m := &metric{
		Name: name,
		Type: t,
		Dims: make(map[string]interface{}),

		env:     e,
		dimlock: sync.Mutex{},
	}

	if dims != nil {
		for k, v := range *dims {
			m.AddDimension(k, v)
		}
	}
	return m
}
