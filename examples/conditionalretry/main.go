package main

import (
	"fmt"
	"time"

	retry "github.com/weathersource/go-retry"
	context "golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func main() {
	var i int

	retryIf := func(err error) bool {

		// test err against target interfaces
		sErr, sOk := err.(interface{ GRPCStatus() *(status.Status) })
		tiErr, tiOk := err.(interface{ Timeout() bool })
		teErr, teOk := err.(interface{ Temporary() bool })

		// Never retry if cancelled or timeout
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

	retryFunc := func() error {
		// include retriable logic here
		i++
		return nil
	}

	err := retry.Do(
		retryFunc,
		retry.WithAttempts(5),
		retry.WithExponentialBackoff(1*time.Second, 10*time.Second, 2),
		retry.WithConditionalRetry(retryIf),
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("i =", i)
}
