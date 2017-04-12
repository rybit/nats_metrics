package metrics

import (
	"sync"

	"time"

	"github.com/nats-io/nats"
)

var globalEnv *Environment

var initLock = sync.Mutex{}

// Init will setup the global context
func Init(nc *nats.Conn, subject string) error {
	initLock.Lock()
	defer initLock.Unlock()

	if globalEnv == nil {
		var err error
		globalEnv, err = NewEnvironment(nc, subject)
		if err != nil {
			return err
		}
	} else {
		return DoubleInitError{errString{"double init attempted"}}
	}

	return globalEnv.isReady()
}

func GlobalEnv() *Environment {
	return globalEnv
}

// AddDimension will let you store a dimension in the global space
func AddDimension(key string, value interface{}) {
	if globalEnv != nil {
		globalEnv.AddDimension(key, value)
	}
}

// NewCounter creates a named counter with these dimensions
func NewCounter(name string, metricDims DimMap) Counter {
	return globalEnv.NewCounter(name, metricDims)
}

// NewGauge creates a named gauge with these dimensions
func NewGauge(name string, metricDims DimMap) Gauge {
	return globalEnv.NewGauge(name, metricDims)
}

// NewTimer creates a named timer with these dimensions
func NewTimer(name string, metricDims DimMap) Timer {
	timer := globalEnv.NewTimer(name, metricDims)
	timer.Start()
	return timer
}

// TimeBlock will just time the block provided
func TimeBlock(name string, metricDims DimMap, f func()) time.Duration {
	return globalEnv.timeBlock(name, metricDims, f)
}

// TimeBlockErr will run the function and publish the time it took.
// It will add the dimension 'had_error' and return the error from the internal function
func TimeBlockErr(name string, metricDims DimMap, f func() error) (time.Duration, error) {
	return globalEnv.timeBlockErr(name, metricDims, f)
}

func Trace(tracer func(m *RawMetric)) {
	globalEnv.tracer = tracer
}

func Count(name string, metricDims DimMap) error {
	return globalEnv.NewCounter(name, nil).Count(metricDims)
}

func CountN(name string, val int64, metricDims DimMap) error {
	return globalEnv.NewCounter(name, nil).CountN(val, metricDims)
}
