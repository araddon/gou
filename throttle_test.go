package gou

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestThrottleer(t *testing.T) {
	th := NewThrottler(10, 10)
	for i := 0; i < 10; i++ {
		assert.Tf(t, th.Throttle() == false, "Should not throttle %v", i)
		time.Sleep(time.Millisecond * 10)
	}
	throttled := 0
	th = NewThrottler(10, 1)
	// We are going to loop 20 times, first 10 should make it, next 10 throttled
	for i := 0; i < 20; i++ {
		if th.Throttle() {
			throttled += 1
		}
	}
	assert.Tf(t, throttled == 10, "Should throttle 10 of 20 requests: %v", throttled)
}
