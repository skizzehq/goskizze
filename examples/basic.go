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

	name := "testdomain"
	_, err = client.CreateDomain(name)
	if err != nil {
		fmt.Printf("Error creating domain: %s\n", err)
	}

	// NOTE: Error checking is removed for readibility, however errors should always be
	// checked, especially with distributed systems (can't trust anything!)
	client.AddToDomain(name, "alvin", "simon", "theodore")
	printSketches(name, client)

	// Let's add some more values to the domain
	client.AddToDomain(name, "alvin", "alvin", "simon", "claire", "patrick", "rajendra")
	printSketches(name, client)
}

func printSketches(name string, client *skizze.Client) {
	fmt.Println("")

	membs, _ := client.GetMembership(name, "alvin", "simon", "theodore", "gary")
	for _, m := range membs {
		fmt.Printf("MEMB: %s is in %s: %v\n", m.Value, name, m.IsMember)
	}
	fmt.Println("")

	freqs, _ := client.GetFrequency(name, "alvin", "simon", "theodore", "gary")
	for _, f := range freqs {
		fmt.Printf("FREQ: %s appears in %s %v times\n", f.Value, name, f.Count)
	}
	fmt.Println("")

	ranks, _ := client.GetRankings(name)
	for i, r := range ranks {
		fmt.Printf("RANK: #%v = %s (count=%v)\n", i, r.Value, r.Count)
	}
	fmt.Println("")

	card, _ := client.GetCardinality(name)
	fmt.Printf("CARD: There are %v items in the %s domain\n\n", card, name)
}
