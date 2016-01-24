# goskizze ![](https://travis-ci.org/skizzehq/goskizze.svg?branch=master)

goskizze is a [Go](https://golang.org) client for the [Skizze](https://github.com/skizzehq/skizze) database.

### Features
 * Simple API with 100% coverage
 * Thread-safe


### Documentation
 * [API Reference](https://godoc.org/github.com/skizzehq/goskizze/skizze)
 * [Skizze Server](https://github.com/skizzehq/skizze)
 

### Installation
Install goskizze with the `go get` command:

```
go get github.com/skizzehq/goskizze/skizze
```

### Example

```golang
package main

import (
	"fmt"
	"github.com/skizzehq/goskizze/skizze"
)

func main() {
	client, err := skizze.Dial("127.0.0.1:3596", skizze.Options{Insecure: true})
	if err != nil {
		fmt.Printf("Error connecting to Skizze: %s\n", err)
		return
	}
	
	// A domain is an easy way to use the same data set for multiple statistics
	name := "testdomain"
	client.CreateDomain(name)

	// Adding values to a domain will trigger statistics generation for each of
	// the supported Sketches in the domain
	client.AddToDomain(name, "alvin", "simon", "theodore")
	
	// The Membership sketch will test if a value resides in a data set, returning
	// true or false
	membs, _ := client.GetMembership(name, "alvin", "simon", "theodore", "gary")
	for _, m := range membs {
		fmt.Printf("MEMB: %s is in %s: %v\n", m.Value, name, m.IsMember)
	}


	// The Frequency sketch will return how many times a value occurs in a sketch
	freqs, _ := client.GetFrequency(name, "alvin", "simon", "theodore", "gary")
	for _, f := range freqs {
		fmt.Printf("FREQ: %s appears in %s %v times\n", f.Value, name, f.Count)
	}

	// The Rankings sketch will always keep the top N (configurable) rankings and
	// their occurrance counts
	ranks, _ := client.GetRankings(name)
	for i, r := range ranks {
		fmt.Printf("RANK: #%v = %s (count=%v)\n", i, r.Value, r.Count)
	}

	// Finally, the Cardinality sketch will keep a count of how many unique items
	// have been added to the data set
	card, _ := client.GetCardinality(name)
	fmt.Printf("CARD: There are %v items in the %s domain\n\n", card, name)
	
	client.DeleteDomain(name)
}

```

### TODO
 * [ ] Benchmarking
 * [ ] Reduce allocations


### License
goskizze is available under the Apache License, Version 2.0.


### Authors
- [Neil Jagdish Patel](https://twitter.com/njpatel)