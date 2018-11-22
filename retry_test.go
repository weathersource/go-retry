package retry

import (
	"errors"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
	wxErrors "github.com/weathersource/go-errors"
)

func TestDo(t *testing.T) {
	tests := []struct {
		retryableFunc RetryableFunc
		opts          []Option
		result        error
	}{
		{
			retryableFunc: func() error { return nil },
			opts:          []Option{},
			result:        nil,
		},
		{
			retryableFunc: func() error { return errors.New("foo") },
			opts: []Option{
				WithConstantBackoff(1 * time.Millisecond),
				WithAttempts(2),
			},
			result: wxErrors.NewErrors(errors.New("foo"), errors.New("foo")),
		},
		{
			retryableFunc: func() error { return errors.New("foo") },
			opts: []Option{
				WithAttempts(2),
				WithExponentialBackoff(1*time.Millisecond, 10*time.Second),
				WithConditionalRetry(func(err error) bool { return true }),
			},
			result: wxErrors.NewErrors(errors.New("foo"), errors.New("foo")),
		},
		{
			retryableFunc: func() error { return errors.New("foo") },
			opts: []Option{
				WithConditionalRetry(func(err error) bool { return false }),
			},
			result: wxErrors.NewErrors(errors.New("foo")),
		},
	}
	for _, test := range tests {
		err := Do(test.retryableFunc, test.opts...)
		if test.result == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, test.result.Error(), err.Error())
		}
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		config   *config
		attempt  int
		expected time.Duration
	}{
		{
			config: &config{
				initDelay: 1 * time.Second,
				maxDelay:  10 * time.Second,
				factor:    2,
			},
			attempt:  3,
			expected: 5651920372,
		},
		{
			config: &config{
				initDelay: 1 * time.Second,
				maxDelay:  10 * time.Second,
				factor:    maxInt64,
			},
			attempt:  3,
			expected: 10 * time.Second,
		},
		{
			config: &config{
				initDelay: 1 * time.Second,
				maxDelay:  10 * time.Second,
				factor:    2,
			},
			attempt:  20,
			expected: 10 * time.Second,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, test.config.duration(test.attempt))
	}
}
