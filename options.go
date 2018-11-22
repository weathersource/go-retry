package retry

import (
	"time"
)

// Option is the type of constructor options for Do(...).
type Option func(*config)

// WithAttempts configures the number of retry attampts.
func WithAttempts(attempts int) Option {
	if attempts < 1 {
		attempts = DefaultAttempts
	}
	return func(c *config) {
		c.attempts = attempts
	}
}

// WithConstantBackoff configures a constant delay between retries.
func WithConstantBackoff(delay time.Duration) Option {
	if delay < 0 {
		delay = DefaultInitDelay
	}
	return func(c *config) {
		c.initDelay = delay
		c.maxDelay = delay
		c.factor = 1
	}
}

// WithExponentialBackoff configures an exponential delay between retries.
func WithExponentialBackoff(initDelay time.Duration, maxDelay time.Duration, factor ...float64) Option {
	var fac float64
	if len(factor) > 0 {
		fac = factor[0]
	} else {
		fac = DefaultFactor
	}

	if initDelay < 0 {
		initDelay = DefaultInitDelay
	}
	if maxDelay < 0 {
		maxDelay = DefaultMaxDelay
	}
	if initDelay > maxDelay {
		initDelay = maxDelay
	}
	if fac <= 1 {
		fac = DefaultFactor
	}
	return func(c *config) {
		c.initDelay = initDelay
		c.maxDelay = maxDelay
		c.factor = fac
	}
}

// ConditionalRetryFunc is the function signature for a conditional retry
// function passed to WithConditionalRetry(...).
type ConditionalRetryFunc func(error) bool

// WithConditionalRetry configures conditions upon which a retry should be attempted.
func WithConditionalRetry(conditionalRetry ConditionalRetryFunc) Option {
	return func(c *config) {
		c.conditionalRetry = conditionalRetry
	}
}
