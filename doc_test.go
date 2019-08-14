package retry_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	retry "github.com/weathersource/go-retry"
	context "golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func Example() {

	// A result placeholder with scope external to retryFunc.
	var res *http.Response

	// A retryable function of type RetryableFunc
	retryFunc := func() error {
		var err error
		res, err = http.Get("https://example.com")
		if err != nil {
			return err
		}
		return nil
	}

	// A conditional retry function of type ConditionalRetryFunc
	retryIf := func(err error) bool {
		return true
	}

	// Retry retryFunc with options.
	err := retry.Do(
		retryFunc,
		retry.WithAttempts(5),
		retry.WithExponentialBackoff(1*time.Second, 10*time.Second, 2),
		retry.WithConditionalRetry(retryIf),
	)

	// handle the output
	if err != nil {
		fmt.Println(err)
	} else {
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s", body)
	}
}

func ExampleDo() {
	// A retryable function of type RetryableFunc
	retryFunc := func() error {
		// retryable code here
		return nil
	}

	// Retry retryFunc.
	err := retry.Do(retryFunc)

	// Do something
	fmt.Println(err)
}

func ExampleConditionalRetryFunc() {
	// Check if err is recoverable
	retryIf := func(err error) bool {

		// test err against target interfaces
		sErr, sOk := err.(interface{ GRPCStatus() *(status.Status) })
		tiErr, tiOk := err.(interface{ Timeout() bool })
		teErr, teOk := err.(interface{ Temporary() bool })

		// Never retry if canceled or timeout
		switch {
		case err == context.DeadlineExceeded:
			fallthrough
		case err == context.Canceled:
			fallthrough
		case sOk && codes.DeadlineExceeded == sErr.GRPCStatus().Code():
			fallthrough
		case sOk && codes.Canceled == sErr.GRPCStatus().Code():
			fallthrough
		case tiOk && tiErr.Timeout():
			return false
		}

		// Do retry if temporary error
		switch {
		case sOk && codes.Unknown == sErr.GRPCStatus().Code():
			fallthrough
		case sOk && codes.Unavailable == sErr.GRPCStatus().Code():
			fallthrough
		case sOk && codes.ResourceExhausted == sErr.GRPCStatus().Code():
			fallthrough
		case teOk && teErr.Temporary():
			return true
		}

		// Otherwise, do not retry
		return false
	}

	// A retryable function of type RetryableFunc
	retryFunc := func() error {
		// retryable code here
		return nil
	}

	// Retry retryFunc.
	err := retry.Do(
		retryFunc,
		retry.WithConditionalRetry(retryIf),
	)

	// Do something
	fmt.Println(err)
}

func ExampleWithAttempts() {
	// A retryable function of type RetryableFunc
	retryFunc := func() error {
		// retryable code here
		return nil
	}

	// Retry retryFunc with WithAttempts option.
	err := retry.Do(
		retryFunc,
		retry.WithAttempts(5),
	)

	// Do something
	fmt.Println(err)
}

func ExampleWithConstantBackoff() {
	// A retryable function of type RetryableFunc
	retryFunc := func() error {
		// retryable code here
		return nil
	}

	// Retry retryFunc with WithConstantBackoff option.
	err := retry.Do(
		retryFunc,
		retry.WithConstantBackoff(1*time.Second),
	)

	// Do something
	fmt.Println(err)
}

func ExampleWithExponentialBackoff() {
	// A retryable function of type RetryableFunc
	retryFunc := func() error {
		// retryable code here
		return nil
	}

	// Retry retryFunc with WithExponentialBackoff option.
	err := retry.Do(
		retryFunc,
		retry.WithExponentialBackoff(1*time.Second, 10*time.Second, 2),
	)

	// Do something
	fmt.Println(err)
}

func ExampleWithConditionalRetry() {
	// A retryable function of type RetryableFunc
	retryFunc := func() error {
		// retryable code here
		return nil
	}

	// A conditional retry function of type ConditionalRetryFunc
	retryIf := func(err error) bool {
		return true
	}

	// Retry retryFunc with WithConditionalRetry option.
	err := retry.Do(
		retryFunc,
		retry.WithConditionalRetry(retryIf),
	)

	// Do something
	fmt.Println(err)
}
