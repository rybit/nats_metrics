package metrics

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/nats-io/nats"
)

func NewEnvironment(nc *nats.Conn, subject string) (*Environment, error) {
	env := &Environment{
		subject:    subject,
		dimlock:    sync.Mutex{},
		globalDims: DimMap{},
	}

	if nc != nil {
		econn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
		if err != nil {
			return nil, err
		}
		env.nc = econn
	}

	if err := env.isReady(); err != nil {
		return nil, err
	}

	return env, nil
}

type Environment struct {
	subject    string
	globalDims DimMap
	dimlock    sync.Mutex
	nc         *nats.EncodedConn

	tracer func(m *RawMetric)

	// some metrics stuff
	timersSent   int64
	countersSent int64
	gaugesSent   int64
}

func (e *Environment) send(m *RawMetric) error {
	if err := e.isReady(); err != nil {
		return err
	}

	switch m.Type {
	case CounterType:
		atomic.AddInt64(&e.countersSent, 1)
	case TimerType:
		atomic.AddInt64(&e.timersSent, 1)
	case GaugeType:
		atomic.AddInt64(&e.gaugesSent, 1)
	default:
		return UnknownMetricTypeError{errString{fmt.Sprintf("unknown metric type: %s", m.Type)}}
	}

	if e.tracer != nil {
		go e.tracer(m)
	}

	if e.nc == nil {
		return nil
	}

	return e.nc.Publish(e.subject, &m)
}

func (e *Environment) AddDimension(k string, v interface{}) {
	e.dimlock.Lock()
	defer e.dimlock.Unlock()
	e.globalDims[k] = v
}

func (e *Environment) isReady() error {
	if e.subject == "" {
		return InitError{errString{"No subject provided"}}
	}
	return nil
}

func addAll(into DimMap, from DimMap) {
	if into != nil {
		if from != nil {
			for k, v := range from {
				into[k] = v
			}
		}
	}
}

func (e *Environment) newMetric(name string, t MetricType, dims DimMap) *metric {
	m := &metric{
		RawMetric: RawMetric{
			Name: name,
			Type: t,
			Dims: make(DimMap),
		},
		env:     e,
		dimlock: sync.Mutex{},
	}

	if dims != nil {
		for k, v := range dims {
			m.AddDimension(k, v)
		}
	}
	return m
}
