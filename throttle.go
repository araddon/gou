package gou

import (
	"time"
)

type Throttler struct {

	// Limit to this events/second
	maxPerSec int

	// Last Event
	last time.Time

	// How many events are allowed left to happen?
	// Starts at limit, decrements down
	allowance float64

	per float64
}

// new Throttler that will tell you to limit or not based
// on given max events per second input @limit
func NewThrottler(maxPerSecond, per int) *Throttler {
	return &Throttler{
		maxPerSec: maxPerSecond,
		allowance: float64(maxPerSecond),
		last:      time.Now(),
		per:       float64(per),
	}
}

// Should we limit this because we are above rate?
func (r *Throttler) Throttle() bool {

	if r.maxPerSec == 0 {
		return false
	}

	// http://stackoverflow.com/questions/667508/whats-a-good-rate-limiting-algorithm
	rate := float64(r.maxPerSec)
	now := time.Now()
	elapsed := float64(now.Sub(r.last).Nanoseconds()) / 1e9 // nano Seconds
	r.last = now
	r.allowance += elapsed * (rate / r.per)

	//Infof("maxRate: %v  cur: %v elapsed:%-6.6f  incr: %v", r.maxPerSec, int(r.allowance), elapsed, elapsed*float64(r.maxPerSec))
	if r.allowance > rate {
		r.allowance = float64(r.maxPerSec)
	}

	if r.allowance <= 1.0 {
		return true // do throttle/limit
	}

	r.allowance -= 1.0
	return false // dont throttle
}

/*
type Ratelimiter struct {

    rate  int    // conn/sec
    last  time.Time  // last time we were polled/asked

    allowance float64
}

// Create new rate limiter that limits at rate/sec
func NewRateLimiter(rate int) (*Ratelimiter, error) {

    r := Ratelimiter{rate:rate, last:time.Now()}

    r.allowance = float64(r.rate)
    return &r, nil
}

// Return true if the current call exceeds the set rate, false
// otherwise
func (r* Ratelimiter) Limit() bool {

    // handle cases where rate in config file is unset - defaulting
    // to "0" (unlimited)
    if r.rate == 0 {
        return false
    }

    rate        := float64(r.rate)
    now         := time.Now()
    elapsed     := now.Sub(r.last)
    r.last       = now
    r.allowance += float64(elapsed) * rate

    // Clamp number of tokens in the bucket. Don't let it get
    // unboundedly large
    if r.allowance > rate {
        r.allowance = rate
    }

    var ret bool

    if r.allowance < 1.0 {
        ret = true
    } else {
        r.allowance -= 1.0
        ret = false
    }

    return ret
}
*/
