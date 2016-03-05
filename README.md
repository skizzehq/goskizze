# goskizze ![](https://travis-ci.org/skizzehq/goskizze.svg?branch=master)

goskizze is a [Go](https://golang.org) client for the [Skizze](https://github.com/skizzehq/skizze) database.


### Documentation
 * [API Reference](https://godoc.org/github.com/skizzehq/goskizze/skizze)
 * [Skizze Server](https://github.com/skizzehq/skizze)
 

### Installation
Install goskizze with the `go get` command:

```
go get gopkg.in/skizzehq/goskizze.v1/skizze
```

### Example

```go
package main

import (
	"fmt"
	"gopkg.in/skizzehq/goskizze.v1/skizze"
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

Output:


```
MEMB: alvin is in testdomain: true
MEMB: simon is in testdomain: true
MEMB: theodore is in testdomain: true
MEMB: gary is in testdomain: false

FREQ: alvin appears in testdomain 1 times
FREQ: simon appears in testdomain 1 times
FREQ: theodore appears in testdomain 1 times
FREQ: gary appears in testdomain 0 times

RANK: #0 = alvin (count=1)
RANK: #1 = simon (count=1)
RANK: #2 = theodore (count=1)

CARD: There are 3 items in the testdomain domain
```


Note: Error checking has been removed for readability, but should be done in production code.


### TODO
 * [x] Support customized domain/sketch creation (with properties)
 * [ ] Benchmarking
 * [ ] Reduce allocations


### License
goskizze is available under the Apache License, Version 2.0.


### Authors
- [Neil Jagdish Patel](https://twitter.com/njpatel)
- [Seif Lotfy](https://twitter.com/seiflotfy)
