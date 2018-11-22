package retry

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
)

func TestWithAttempts(t *testing.T) {
	tests := []struct {
		attempts, expected int
	}{
		{
			attempts: -1,
			expected: DefaultAttempts,
		},
		{
			attempts: 0,
			expected: DefaultAttempts,
		},
		{
			attempts: 1,
			expected: 1,
		},
		{
			attempts: 2,
			expected: 2,
		},
	}
	c := &config{}
	for _, test := range tests {
		f := WithAttempts(test.attempts)
		f(c)
		assert.Equal(t, test.expected, c.attempts)
	}
}

func TestWithConstantBackoff(t *testing.T) {
	tests := []struct {
		delay             time.Duration
		expectedInitDelay time.Duration
		expectedMaxDelay  time.Duration
		expectedFactor    float64
	}{
		{
			delay:             1 * time.Second,
			expectedInitDelay: 1 * time.Second,
			expectedMaxDelay:  1 * time.Second,
			expectedFactor:    1,
		},
		{
			delay:             -1 * time.Second,
			expectedInitDelay: DefaultInitDelay,
			expectedMaxDelay:  DefaultInitDelay,
			expectedFactor:    1,
		},
		{
			delay:             0,
			expectedInitDelay: 0,
			expectedMaxDelay:  0,
			expectedFactor:    1,
		},
	}
	c := &config{}
	for _, test := range tests {
		f := WithConstantBackoff(test.delay)
		f(c)
		assert.Equal(t, test.expectedInitDelay, c.initDelay)
		assert.Equal(t, test.expectedMaxDelay, c.maxDelay)
		assert.Equal(t, test.expectedFactor, c.factor)
	}
}

func TestWithExponentialBackoff(t *testing.T) {

	tests := []struct {
		initDelay         time.Duration
		maxDelay          time.Duration
		factor            []float64
		expectedInitDelay time.Duration
		expectedMaxDelay  time.Duration
		expectedFactor    float64
	}{
		// missing input factor
		{
			initDelay:         1 * time.Second,
			maxDelay:          2 * time.Second,
			expectedInitDelay: 1 * time.Second,
			expectedMaxDelay:  2 * time.Second,
			expectedFactor:    DefaultFactor,
		},
		// all negative input values
		{
			initDelay:         -1 * time.Second,
			maxDelay:          -1 * time.Second,
			factor:            []float64{-1},
			expectedInitDelay: DefaultInitDelay,
			expectedMaxDelay:  DefaultMaxDelay,
			expectedFactor:    DefaultFactor,
		},
		// out of order init and max delays, zero factor
		{
			initDelay:         2 * time.Second,
			maxDelay:          1 * time.Second,
			factor:            []float64{0},
			expectedInitDelay: 1 * time.Second,
			expectedMaxDelay:  1 * time.Second,
			expectedFactor:    DefaultFactor,
		},
		// zero factor
		{
			initDelay:         1 * time.Second,
			maxDelay:          1 * time.Second,
			factor:            []float64{1},
			expectedInitDelay: 1 * time.Second,
			expectedMaxDelay:  1 * time.Second,
			expectedFactor:    DefaultFactor,
		},
		// correct included factor
		{
			initDelay:         1 * time.Second,
			maxDelay:          1 * time.Second,
			factor:            []float64{2},
			expectedInitDelay: 1 * time.Second,
			expectedMaxDelay:  1 * time.Second,
			expectedFactor:    2,
		},
	}
	c := &config{}
	for _, test := range tests {
		f := WithExponentialBackoff(test.initDelay, test.maxDelay, test.factor...)
		f(c)
		assert.Equal(t, test.expectedInitDelay, c.initDelay)
		assert.Equal(t, test.expectedMaxDelay, c.maxDelay)
		assert.Equal(t, test.expectedFactor, c.factor)
	}
}

func TestWithConditionalRetry(t *testing.T) {
	tests := []struct {
		retryIf  ConditionalRetryFunc
		expected bool
	}{
		{
			retryIf:  func(err error) bool { return true },
			expected: true,
		},
		{
			retryIf:  func(err error) bool { return false },
			expected: false,
		},
	}
	c := &config{}
	for _, test := range tests {
		f := WithConditionalRetry(test.retryIf)
		f(c)
		assert.Equal(t, test.expected, c.conditionalRetry(nil))
	}
}
