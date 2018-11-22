# retry

[![GoDoc](https://godoc.org/github.com/weathersource/go-retry?status.svg)](https://godoc.org/github.com/weathersource/go-retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/weathersource/go-retry)](https://goreportcard.com/report/github.com/weathersource/go-retry)
[![Build Status](https://travis-ci.org/weathersource/go-retry.svg)](https://travis-ci.org/weathersource/go-retry)
[![Codevov](https://codecov.io/gh/weathersource/go-retry/branch/master/graphs/badge.svg)](https://codecov.io/gh/weathersource/go-retry)

Package retry enables automatic retry of code. This package is goroutine-safe and supports:
constant backoff, exponential backoff, conditional retries.

This package started with [avast/retry-go](https://github.com/avast/retry-go) and [jpillora/backoff](https://github.com/jpillora/backoff) and grew from there.
