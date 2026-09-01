// Package calculadora es el código bajo prueba del ejemplo del
// Capítulo 36. Ver calculadora_test.go para los tests de tabla y el
// benchmark.
package calculadora

import "errors"

// ErrDivisionPorCero se retorna cuando el divisor es 0.
var ErrDivisionPorCero = errors.New("división por cero")

func Sumar(a, b int) int { return a + b }

func Dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionPorCero
	}
	return a / b, nil
}

func EsPrimo(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
