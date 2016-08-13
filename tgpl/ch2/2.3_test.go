package main

import (
	"testing"
	"ch2/popcount"
	"math/rand"
	"time"
	"fmt"
)

var value uint64

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	value = uint64(rand.Int63())
}



func BenchmarkUnloopedPopcount(b *testing.B) {
	var result int
	for i := 0; i < b.N; i++ {
		result = popcount.PopCount(value)
	}
	fmt.Printf("unlooped: %d -> %d\n", value, result)
}

func BenchmarkLoopedPopcount(b *testing.B) {
	var result int
	for i := 0; i < b.N; i++ {
		result = popcount.PopCountLooped(value)
	}
	fmt.Printf("looped: %d -> %d\n", value, result)
}

func BenchmarkSlowPopcount(b *testing.B) {
	var result int
	for i := 0; i < b.N; i++ {
		result = popcount.PopCountSlow(value)
	}
	fmt.Printf("slow: %d -> %d\n", value, result)
}

func BenchmarkBitTrickPopcount(b *testing.B) {
	var result int
	for i := 0; i < b.N; i++ {
		result = popcount.PopCountBitTrick(value)
	}
	fmt.Printf("bit trick: %d -> %d\n", value, result)
}


