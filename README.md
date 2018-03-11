# bloomflt
[![GoDoc](https://godoc.org/github.com/quasoft/bloomflt?status.svg)](https://godoc.org/github.com/quasoft/bloomflt) [![Build Status](https://travis-ci.org/quasoft/bloomflt.png?branch=master)](https://travis-ci.org/quasoft/bloomflt) [![Coverage Status](https://coveralls.io/repos/github/quasoft/bloomflt/badge.svg?branch=master)](https://coveralls.io/github/quasoft/bloomflt?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/quasoft/bloomflt)](https://goreportcard.com/report/github.com/quasoft/bloomflt)

Implementation of bloom filter in Golang with no 3rd party dependencies.

Uses the [`big.Int`](https://golang.org/pkg/math/big/#Int) type as bitset storage and two hash functions from the builtin [`hash`](https://golang.org/pkg/hash/) package:
- [FNV-1a (Fowler–Noll–Vo)](https://golang.org/pkg/hash/fnv/)
- [CRC32](https://golang.org/pkg/hash/crc32/)

Those functions are used to simulate an arbitrary number of hash functions using the "Double Hashing Scheme" as described by Kirsch and Mitzenmacher
in "Less Hashing, Same Performance: Building a Better Bloom Filter".

# What is a bloom filter?

BloomFilter is an efficient data structure, used to test whether an element is a member of a set.

Bloom filters are probabilistic, which means that false positives are tolerated, but it is guaranteed
there will be no false negatives.

In other words, if the bloom filter says an element is NOT in the set, then it is guaranteed the element
is not in the set.

If the bloom filter says an element IS in the set, then that means the element is probably there,
but it is not guaranteed.

# How to use?

```go
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
```