package metrics

import (
	"sync"

	"github.com/nats-io/nats"
)

var globalEnv *environment

var initLock = sync.Mutex{}

// Init will setup the global context
func Init(nc *nats.Conn, subject string) error {
	initLock.Lock()
	defer initLock.Unlock()

	if globalEnv == nil {
		var err error
		globalEnv, err = newEnvironment(nc, subject)
		if err != nil {
			return err
		}
	} else {
		return DoubleInitError{errString{"double init attempted"}}
	}

	return globalEnv.isReady()
}

// AddDimension will let you store a dimension in the global space
func AddDimension(key string, value interface{}) {
	if globalEnv != nil {
		globalEnv.AddDimension(key, value)
	}
}

// NewCounter creates a named counter with these dimensions
func NewCounter(name string, metricDims *DimMap) Counter {
	return globalEnv.newCounter(name, metricDims)
}

// NewGauge creates a named gauge with these dimensions
func NewGauge(name string, metricDims *DimMap) Gauge {
	return globalEnv.newGauge(name, metricDims)
}

// NewTimer creates a named timer with these dimensions
func NewTimer(name string, metricDims *DimMap) Timer {
	timer := globalEnv.newTimer(name, metricDims)
	timer.Start()
	return timer
}
