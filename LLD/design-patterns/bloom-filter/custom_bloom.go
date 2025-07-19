package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"hash/fnv"
	"math"
)

// CustomBloomFilter implements a bloom filter from scratch
type CustomBloomFilter struct {
	bitArray  []bool
	size      uint
	hashFuncs []hash.Hash
	numHashes uint
}

// NewCustomBloomFilter creates a new bloom filter with optimal parameters
func NewCustomBloomFilter(expectedItems uint, falsePositiveRate float64) *CustomBloomFilter {
	// Calculate optimal size: m = -n * ln(p) / (ln(2)^2)
	size := uint(-float64(expectedItems) * math.Log(falsePositiveRate) / (math.Log(2) * math.Log(2)))
	
	// Calculate optimal number of hash functions: k = (m/n) * ln(2)
	numHashes := uint(float64(size) / float64(expectedItems) * math.Log(2))
	
	// Ensure we have at least 1 hash function
	if numHashes == 0 {
		numHashes = 1
	}
	
	// Create hash functions
	hashFuncs := []hash.Hash{
		fnv.New64a(),
		md5.New(),
		sha1.New(),
		sha256.New(),
	}
	
	// Use only as many hash functions as we calculated
	if numHashes > uint(len(hashFuncs)) {
		numHashes = uint(len(hashFuncs))
	}
	
	return &CustomBloomFilter{
		bitArray:  make([]bool, size),
		size:      size,
		hashFuncs: hashFuncs[:numHashes],
		numHashes: numHashes,
	}
}

// Add inserts an item into the bloom filter
func (bf *CustomBloomFilter) Add(data []byte) {
	for _, hashFunc := range bf.hashFuncs {
		hashFunc.Reset()
		hashFunc.Write(data)
		hashValue := bf.hashToUint(hashFunc.Sum(nil))
		index := hashValue % bf.size
		bf.bitArray[index] = true
	}
}

// Test checks if an item might be in the bloom filter
func (bf *CustomBloomFilter) Test(data []byte) bool {
	for _, hashFunc := range bf.hashFuncs {
		hashFunc.Reset()
		hashFunc.Write(data)
		hashValue := bf.hashToUint(hashFunc.Sum(nil))
		index := hashValue % bf.size
		if !bf.bitArray[index] {
			return false
		}
	}
	return true
}

// hashToUint converts hash bytes to uint
func (bf *CustomBloomFilter) hashToUint(hashBytes []byte) uint {
	var result uint
	for i, b := range hashBytes {
		if i >= 8 { // Limit to 8 bytes for uint64
			break
		}
		result = result*256 + uint(b)
	}
	return result
}

// GetStats returns statistics about the bloom filter
func (bf *CustomBloomFilter) GetStats() (uint, uint, float64, uint) {
	setBits := uint(0)
	for _, bit := range bf.bitArray {
		if bit {
			setBits++
		}
	}
	
	// Calculate current false positive rate
	// p = (1 - e^(-k*n/m))^k
	ratio := float64(setBits) / float64(bf.size)
	currentFPRate := math.Pow(ratio, float64(bf.numHashes))
	
	return bf.size, setBits, currentFPRate, bf.numHashes
}
