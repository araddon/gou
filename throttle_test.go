package gou

import (
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

func TestThrottleer(t *testing.T) {
	th := NewThrottler(10, 10*time.Second)
	for i := 0; i < 10; i++ {
		assert.Tf(t, th.Throttle() == false, "Should not throttle %v", i)
		time.Sleep(time.Millisecond * 10)
	}
	throttled := 0
	th = NewThrottler(10, 1*time.Second)
	// We are going to loop 20 times, first 10 should make it, next 10 throttled
	for i := 0; i < 20; i++ {
		LogThrottleKey(WARN, 10, "throttle", "hello %v", i)
		if th.Throttle() {
			throttled += 1
		}
	}
	assert.Tf(t, throttled == 10, "Should throttle 10 of 20 requests: %v", throttled)
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
