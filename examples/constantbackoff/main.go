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
		retry.WithConstantBackoff(1*time.Second),
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("i =", i)
}
