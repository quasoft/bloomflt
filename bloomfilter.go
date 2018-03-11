// Package bloomflt is implementation of bloom filter with no 3rd party dependencies.
package bloomflt

import (
	"encoding/binary"
	"hash/crc32"
	"hash/fnv"
	"math"
	"math/big"
)

// BloomFilter is an efficient data structure, used to test whether an element is a member of a set.
//
// Bloom filters are probabilistic, which means that false positives are tolerated, but it is guaranteed
// there will be no false negatives.
//
// In other words, if the bloom filter says an element is NOT in the set, then it is guaranteed the element
// is not in the set.
//
// If the bloom filter says an element IS in the set, then that means the element is probably there,
// but it is not guaranteed.
//
// This implementation uses the big.Int type as bitset storage and FNV-1a (Fowler–Noll–Vo) and CRC32
// from the builtin hash package as the base hash functions. Additional hash functions are simulated
// with "Double Hashing Scheme" by Kirsch and Mitzenmacher as explained in "Less Hashing, Same Performance:
// Building a Better Bloom Filter".
type BloomFilter struct {
	m      int      // Number of elements in the set
	k      int      // Number of hash functions
	bucket *big.Int // Bit storage
}

// NewMK creates a new bloom filter with bucket size equal to m and number of hash functions equal to k.
func NewMK(m int, k int) *BloomFilter {
	filter := BloomFilter{m, k, big.NewInt(0)}

	return &filter
}

// CalcOptimalMK calculates optimal values of m and k for the specified number of elements in the set (n),
// and given acceptable false-positive rate (value from 0.0 to 1.0).
func CalcOptimalMK(n int, falsePositiveRate float64) (int, int) {
	m := -1 * float64(n) * math.Log(falsePositiveRate) / math.Pow(math.Log(2), 2)
	k := m / float64(n) * math.Log(2)
	return int(m + 0.5), int(k + 0.5)
}

// New creates a new bloom filter with optimal values of m and k for the given given acceptable false-positive
// rate (value from 0.0 to 1.0).
func New(n int, falsePositiveRate float64) *BloomFilter {
	m, k := CalcOptimalMK(n, falsePositiveRate)
	// Use at least one bit
	if m < 1 {
		m = 1
	}
	// Limit the number of bits to 2^31
	if m > math.MaxInt32 {
		m = math.MaxInt32
	}
	// Use at least one hash function
	if k < 1 {
		k = 1
	}
	return NewMK(m, k)
}

// FNV-1a (Fowler–Noll–Vo) is used as the first hash function in kiMiHash
func (b *BloomFilter) hash1(value []byte) uint32 {
	f := fnv.New32a()
	f.Write(value)
	hash := f.Sum32()
	return hash
}

// CRC32 is used as the second hash function in kiMiHash
func (b *BloomFilter) hash2(value []byte) uint32 {
	f := crc32.NewIEEE()
	f.Write(value)
	hash := f.Sum32()
	return hash
}

// kiMiHash simulates arbitrary number of hash functions with a "Double Hashing Scheme" by using only
// two hash functions (5.2. in "Less Hashing, Same Performance: Building a Better Bloom Filter" by Kirsch
// and Mitzenmacher). Simplified explanation at:
// http://willwhim.wpengine.com/2011/09/03/producing-n-hash-functions-by-hashing-only-once/
func (b *BloomFilter) kiMiHash(value []byte, hashIdx int) int {
	index := (b.hash1(value) + b.hash2(value)*uint32(hashIdx)) % uint32(b.m)
	return int(index)
}

// AddBytes inserts a bytes value to the set
func (b *BloomFilter) AddBytes(value []byte) {
	for h := 0; h < b.k; h++ {
		index := b.kiMiHash(value, h)
		b.bucket.SetBit(b.bucket, index, 1)
	}
}

// AddString inserts a string value to the set
func (b *BloomFilter) AddString(value string) {
	b.AddBytes([]byte(value))
}

// AddUInt32 inserts an int value to the set
func (b *BloomFilter) AddUInt32(value uint32) {
	bytes := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	b.AddBytes(bytes)
}

// AddUInt64 inserts an int value to the set
func (b *BloomFilter) AddUInt64(value uint64) {
	bytes := make([]byte, 8, 8)
	binary.LittleEndian.PutUint64(bytes, value)
	b.AddBytes(bytes)
}

// ContainsBytes tests if the set contains the given bytes value
func (b *BloomFilter) ContainsBytes(value []byte) bool {
	res := true
	for h := 0; h < b.k; h++ {
		index := b.kiMiHash(value, h)
		if b.bucket.Bit(index) == 0 {
			res = false
		}
	}
	return res
}

// ContainsString tests if the set contains the given string value
func (b *BloomFilter) ContainsString(value string) bool {
	return b.ContainsBytes([]byte(value))
}

// ContainsUInt32 tests if the set contains the given int value
func (b *BloomFilter) ContainsUInt32(value uint32) bool {
	bytes := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	return b.ContainsBytes(bytes)
}

// ContainsUInt64 tests if the set contains the given int value
func (b *BloomFilter) ContainsUInt64(value uint64) bool {
	bytes := make([]byte, 8, 8)
	binary.LittleEndian.PutUint64(bytes, value)
	return b.ContainsBytes(bytes)
}
