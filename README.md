average
[![TravisCI](https://travis-ci.org/prep/average.svg?branch=master)](https://travis-ci.org/prep/average.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/prep/average)](https://goreportcard.com/report/github.com/prep/average)
[![GoDoc](https://godoc.org/github.com/prep/average?status.svg)](https://godoc.org/github.com/prep/average)
=======
This stupidly named Go package contains a single struct that is used to implement counters on a sliding time window.

Usage
-----
```go

import (
    "fmt"

    "github.com/prep/average"
)

func main() {
    // Create a SlidingWindow that has a window of 15 minutes, with a
    // granularity of 1 minute.
    sw := average.MustNew(15*time.Second, time.Second)
    defer sw.Stop()

    sw.Add(15)
    // Do some more work.
    time.Sleep(time.Second)
    sw.Add(22)
    // Do even more work.
    time.Sleep(2 * time.Second)
    sw.Add(22)

    fmt.Printf("Average of last 1s: %f\n", sw.Average(time.Second))
    fmt.Printf("Average of last 2s: %f\n", sw.Average(2*time.Second))
    fmt.Printf("Average of last 15s: %f\n\n", sw.Average(15*time.Second))

    total, numSamples := sw.Total(15 * time.Second)
    fmt.Printf("Counter has a total of %d over %d samples\n", total, numSamples)
}
```

License
-------
This software is created for MessageBird B.V. and distributed under the BSD-style license found in the LICENSE file.
