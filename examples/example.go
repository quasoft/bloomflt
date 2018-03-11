package main

import (
	"fmt"

	"github.com/quasoft/bloomflt"
)

func main() {
	// We expect the set to contains up to 100 elements and accept a false-positive rate of 0.01%
	b := bloomflt.New(100, 0.01)

	b.AddString("value1")
	if b.ContainsString("value1") {
		fmt.Println("The set now has 'value1'.")
	}

	someID := uint64(123)
	b.AddUInt64(someID)
	if b.ContainsUInt64(someID) {
		fmt.Printf("The set now has ID %v.", someID)
	}
}
