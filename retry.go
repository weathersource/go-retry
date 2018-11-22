package retry

import (
	"math"
	"math/rand"
	"time"

	errors "github.com/weathersource/go-errors"
)

const maxInt64 = float64(math.MaxInt64 - 512)

// DefaultAttempts is the default number of retry attempts
const DefaultAttempts = 5

// DefaultInitDelay is the default initial delay between retries
const DefaultInitDelay = 1 * time.Second

// DefaultMaxDelay is the default maximum delay between retries
const DefaultMaxDelay = 10 * time.Second

// DefaultFactor is the default exponential factor for delay growth between retries
const DefaultFactor = float64(2)

// RetryableFunc is a function signature for a retryable function passed to the
// Do function.
type RetryableFunc func() error

// Do is the main function of the retry package. Do retries the retryable
// function subject to the passed options.
func Do(retryableFunc RetryableFunc, opts ...Option) error {

	// default
	c := &config{
		attempts:         DefaultAttempts,
		initDelay:        DefaultInitDelay,
		maxDelay:         DefaultMaxDelay,
		factor:           DefaultFactor,
		conditionalRetry: func(err error) bool { return true },
	}

	// apply options
	for _, opt := range opts {
		opt(c)
	}

	// configure needed variables
	var n int
	errs := make([]error, 0, c.attempts)

	for n < c.attempts {

		err := retryableFunc()

		if err != nil {

			errs = append(errs, err)

			// if this is last attempt - don't wait
			if n >= c.attempts-1 {
				break
			}

			if !c.conditionalRetry(err) {
				break
			}

			time.Sleep(c.duration(n))

		} else {
			return nil
		}

		n++
	}

	return errors.Errors(errs)
}

type config struct {
	attempts         int
	initDelay        time.Duration
	maxDelay         time.Duration
	factor           float64
	conditionalRetry ConditionalRetryFunc
}

func (c *config) duration(attempt int) time.Duration {

	//calculate this duration
	initf := float64(c.initDelay)

	// add jitter
	durf := initf * math.Pow(c.factor, float64(attempt))
	durf = rand.Float64()*(durf-initf) + initf

	//ensure float64 wont overflow int64
	if durf > maxInt64 {
		return c.maxDelay
	}

	dur := time.Duration(durf)

	//keep within bounds
	if dur > c.maxDelay {
		return c.maxDelay
	}
	return dur
}
