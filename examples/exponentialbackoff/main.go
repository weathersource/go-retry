package main

import (
	"fmt"
	"time"

	retry "github.com/weathersource/go-retry"
)

func main() {
	var i int

	retryFunc := func() error {
		// include retriable logic here
		i++
		return nil
	}

	err := retry.Do(
		retryFunc,
		retry.WithAttempts(5),
		retry.WithExponentialBackoff(1*time.Second, 10*time.Second, 2),
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("i =", i)
}
