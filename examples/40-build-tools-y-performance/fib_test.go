package fib

import "testing"

func TestMemoizadoIgualQueNaive(t *testing.T) {
	for n := 0; n <= 15; n++ {
		if got, want := Memoizado(n), Naive(n); got != want {
			t.Errorf("Memoizado(%d) = %d; Naive(%d) = %d", n, got, n, want)
		}
	}
}

// Compara ambas implementaciones con:
//
//	go test ./40-build-tools-y-performance/ -bench=. -benchmem
//
// Para profiling real (CPU y memoria):
//
//	go test ./40-build-tools-y-performance/ -bench=Naive -cpuprofile=cpu.prof -memprofile=mem.prof
//	go tool pprof cpu.prof
func BenchmarkNaive(b *testing.B) {
	for b.Loop() {
		Naive(20)
	}
}

func BenchmarkMemoizado(b *testing.B) {
	for b.Loop() {
		Memoizado(20)
	}
}
