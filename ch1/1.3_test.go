package main

import (
	"strings"
	"testing"
)

var tokens = []string{
	"aaa",
	"bbb",
	"ccc",
	"ddd",
	"eee",
	"aaa",
	"bbb",
	"ccc",
	"ddd",
	"eee",
}

func BenchmarkNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for _, token := range tokens {
			s += token
		}
	}
}

func BenchmarkJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Join(tokens, "")
	}
}
