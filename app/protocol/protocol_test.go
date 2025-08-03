package protocol

import (
	"math/big"
	"runtime"
	"strconv"
	"testing"
)

var p = []byte("23")

func BenchmarkBigMath(b *testing.B) {
	var x int

	for b.Loop() {
		x = int(big.NewInt(0).SetBytes(p).Int64())
	}

	runtime.KeepAlive(x)
}

func BenchmarkString(b *testing.B) {
	var x int

	for b.Loop() {
		x, _ = strconv.Atoi(string(p))
	}

	runtime.KeepAlive(x)
}
