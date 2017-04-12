package metrics

import (
	"testing"
	"time"

	"errors"

	"github.com/stretchr/testify/assert"
)

func TestTimeIt(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	timer := env.NewTimer("something", nil)
	start := timer.Start()
	<-time.After(time.Millisecond * 100)
	stop := time.Now()
	_, err := timer.Stop(nil)
	assert.Nil(t, err)

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
		measured := start.Add(time.Duration(m.Value))
		assert.WithinDuration(t, stop, measured, time.Millisecond*10)
	}
}

func TestTimeBlock(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	wasCalled := false
	env.timeBlock("something", DimMap{"pokemon": "pikachu"}, func() {
		wasCalled = true
	})

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
		assert.True(t, wasCalled)
		assert.Equal(t, "pikachu", m.Dims["pokemon"])
		assert.NotZero(t, m.Value)
	}
}

func TestTimeBlockErr(t *testing.T) {
	sub, env := subscribe(t)
	defer sub.Unsubscribe()

	wasCalled := false
	madeErr := errors.New("garbage error")
	_, err := env.timeBlockErr("something", DimMap{"pokemon": "pikachu"}, func() error {
		wasCalled = true
		return madeErr
	})

	m := readOne(t, sub)
	if assert.NotNil(t, m) {
		assert.True(t, wasCalled)
		assert.Equal(t, madeErr, err)
		assert.Equal(t, "pikachu", m.Dims["pokemon"])
		assert.NotZero(t, m.Value)
	}
}
