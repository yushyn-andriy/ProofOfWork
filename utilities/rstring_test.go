package utilities

import (
	"math/rand"

	"testing"
	"time"
)

const (
	STRING_LEN = 8
)

func BenchmarkRandStringBytesMaskImprSrcUnsafe(b *testing.B) {
	// run the Fib function b.N times
	src := rand.NewSource(time.Now().UnixNano())
	for n := 0; n < b.N; n++ {
		RandStringBytesMaskImprSrcUnsafe(STRING_LEN, src)
	}
}

func BenchmarkRandUnicodeString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandUnicodeString(STRING_LEN)
	}
}

func BenchmarkRandStringRunes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandStringRunes(STRING_LEN)
	}
}

func BenchmarkRandBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandBytes(STRING_LEN)
	}
}

func BenchmarkIntn(b *testing.B) {
	for n := 0; n < b.N; n++ {
		rand.Intn(10000000) // 15.83 ns/op
	}
}

func BenchmarkInt63(b *testing.B) {
	src := rand.NewSource(time.Now().UnixNano())
	for n := 0; n < b.N; n++ {
		src.Int63() // 1.848 ns/op
	}
}
