// Package fib acompaña al Capítulo 40: dos implementaciones de Fibonacci
// con muy distinto perfil de performance, para comparar con benchmarks
// y pprof (ver fib_test.go).
package fib

// Naive es O(2^n): útil solo para ilustrar el problema.
func Naive(n int) int {
	if n < 2 {
		return n
	}
	return Naive(n-1) + Naive(n-2)
}

// Memoizado es O(n) gracias a cachear resultados ya calculados.
func Memoizado(n int) int {
	cache := make([]int, n+1)
	var calcular func(int) int
	calcular = func(n int) int {
		if n < 2 {
			return n
		}
		if cache[n] != 0 {
			return cache[n]
		}
		resultado := calcular(n-1) + calcular(n-2)
		cache[n] = resultado
		return resultado
	}
	return calcular(n)
}
