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

func checkEnv() error {
	if globalEnv == nil {
		return InitError{errString{"the global environment hasn't been configured"}}
	}

	return globalEnv.isReady()
}
