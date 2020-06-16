# go-retry

[![CircleCI](https://circleci.com/gh/weathersource/go-retry.svg?style=shield)](https://circleci.com/gh/weathersource/go-retry)
[![GoDoc](https://img.shields.io/badge/godoc-ref-blue.svg)](https://godoc.org/github.com/weathersource/go-retry)

Package retry enables automatic retry of code. This package is goroutine-safe and supports:
constant backoff, exponential backoff, conditional retries.

This package started with [avast/retry-go](https://github.com/avast/retry-go) and [jpillora/backoff](https://github.com/jpillora/backoff) and grew from there.
