package main

import (
	"fmt"

	retry "github.com/weathersource/go-retry"
)

func main() {
	var i int

	retryFunc := func() error {
		// include retriable logic here
		i++
		return nil
	}

	err := retry.Do(retryFunc)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("i =", i)
}
