package gou

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThrottleer(t *testing.T) {
	th := NewThrottler(10, 10*time.Second)
	for i := 0; i < 10; i++ {
		thb, tc := th.Throttle()
		assert.True(t, thb == false, "Should not throttle %v", i)
		assert.True(t, tc < 10, "Throttle count should remain below 10 %v", tc)
		time.Sleep(time.Millisecond * 10)
	}

	throttled := 0
	th = NewThrottler(10, 1*time.Second)
	// We are going to loop 20 times, first 10 should make it, next 10 throttled
	for i := 0; i < 20; i++ {
		LogThrottleKey(WARN, 10, "throttle", "hello %v", i)
		thb, tc := th.Throttle()

		if thb {
			throttled += 1
			assert.True(t, int(tc) == i-9, "Throttle count should rise %v, i: %d", tc, i)
		}
	}
	assert.True(t, throttled == 10, "Should throttle 10 of 20 requests: %v", throttled)

	// Now sleep for 1 second so that we should
	// no longer be throttled
	time.Sleep(time.Second * 2)
	thb, _ := th.Throttle()
	assert.True(t, thb == false, "We should not have been throttled")
}

func TestThrottler2(t *testing.T) {

	th := NewThrottler(10, 1*time.Second)

	tkey := "throttle2"
	throttleMu.Lock()
	logThrottles[tkey] = th
	throttleMu.Unlock()

	th, ok := logThrottles[tkey]
	if !ok {
		t.Errorf("Throttle key %s not created!", tkey)
	}

	// We are going to loop 20 times, first 10 should make it, next 10 throttled
	for i := 0; i < 20; i++ {

		LogThrottleKey(WARN, 10, tkey, "hello %v", i)

	}

	throttleMu.Lock()
	th = logThrottles[tkey]
	tcount := th.ThrottleCount()
	assert.True(t, tcount == 10, "Should throttle 10 of 20 requests: %v", tcount)
	throttleMu.Unlock()

	// Now sleep for 1 second so that we should
	// no longer be throttled
	time.Sleep(time.Second * 1)
	LogThrottleKey(WARN, 10, tkey, "hello again %v", 20)

}

func TestThrottle(t *testing.T) {
	throttled := 0
	th := NewThrottler(10, 1*time.Second)
	for i := 0; i < 20; i++ {
		LogThrottle(WARN, 10, "throttle", "hihi %v", i)
		if th.Throttle() {
			throttled++
		}
	}
	assert.Tf(t, throttled == 10, "Should throttle 10 of 20 requests: %v", throttled)
}

func TestThrottleLow(t *testing.T) {
	throttled := 0
	th := NewThrottler(100, 1*time.Second)
	start := time.Now()
	for i := 0; i < 200; i++ {
		LogThrottleKey(WARN, 100, "throttle", "hello %v", i)
		//LogThrottle(WARN, 100, "throttle", "hihi %v", i)
		if th.Throttle() {
			throttled++
		}
	}
	Infof("\n\nTime taken Low: %v\n", time.Since(start))
	assert.Tf(t, throttled == 100, "Should throttle 100 of 200 requests: %v", throttled)
}

func TestThrottleMed(t *testing.T) {
	throttled := 0
	th := NewThrottler(1000, 1*time.Second)
	start := time.Now()
	for i := 0; i < 2000; i++ {
		LogThrottleKey(WARN, 1000, "throttle", "hello %v", i)
		//LogThrottle(WARN, 1000, "throttle", "hihi %v", i)
		if th.Throttle() {
			throttled++
		}
	}
	Infof("\n\nTime taken Med: %v\n", time.Since(start))
	assert.Tf(t, throttled == 1000, "Should throttle 1000 of 2000 requests: %v", throttled)
}

func TestThrottleMed5(t *testing.T) {
	throttled := 0
	th := NewThrottler(1000, 1*time.Second)
	start := time.Now()
	for i := 0; i < 2000; i++ {
		LogThrottle(WARN, 1000, "throttle", "hihi %v", i)
		if th.Throttle() {
			throttled++
		}
	}
	Infof("\n\nTime taken Med5: %v\n", time.Since(start))
	assert.Tf(t, throttled == 5000, "Should throttle 5000 of 10000 requests: %v", throttled)
}

func TestThrottleBig(t *testing.T) {
	throttled := 0
	start := time.Now()
	th := NewThrottler(10000, 1*time.Second)
	for i := 0; i < 20000; i++ {
		LogThrottle(WARN, 10000, "throttle", "hihi %v", i)
		if th.Throttle() {
			throttled++
		}
	}
	Infof("\n\nTime taken Big: %v\n", time.Since(start))
	assert.Tf(t, throttled == 10000, "Should throttle 10000 of 20000 requests: %v", throttled)
}

func TestThrottleBigger(t *testing.T) {
	throttled := 0
	th := NewThrottler(100000, 1*time.Second)
	start := time.Now()
	for i := 0; i < 200000; i++ {
		LogThrottle(WARN, 100000, "throttle", "hihi %v", i)
		if th.Throttle() {
			throttled++
		}
	}
	Infof("\n\nTime taken: %v\n", time.Since(start))
	assert.Tf(t, throttled == 100000, "Should throttle 100000 of 200000 requests: %v", throttled)
}

func TestThrottleAbsurd(t *testing.T) {
	throttled := 0
	start := time.Now()
	th := NewThrottler(1000000, 1*time.Second)
	for i := 0; i < 2000000; i++ {
		LogThrottle(WARN, 1000000, "throttle", "hihihihi %v", i)
		//fmt.Printf("Wat: %v\n", throttled)
		if th.Throttle() {
			throttled++
		}
	}
	Infof("\n\nTime taken absurd: %v\n", time.Since(start))
	assert.Tf(t, throttled == 1000000, "Should throttle 1000000 of 2000000 requests, but seems to break limiting algorithm: %v", throttled)
}
