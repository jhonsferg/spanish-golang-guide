package calculadora

import (
	"errors"
	"testing"
)

// Test de tabla: el patrón idiomático de Go para cubrir varios casos
// sin duplicar la estructura del test.
func TestSumar(t *testing.T) {
	casos := []struct {
		nombre   string
		a, b     int
		esperado int
	}{
		{"positivos", 2, 3, 5},
		{"con negativo", -1, 5, 4},
		{"ambos cero", 0, 0, 0},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			t.Parallel() // este subtest no comparte estado con los demás
			got := Sumar(c.a, c.b)
			if got != c.esperado {
				t.Errorf("Sumar(%d, %d) = %d; esperado %d", c.a, c.b, got, c.esperado)
			}
		})
	}
}

func TestDividir(t *testing.T) {
	if _, err := Dividir(10, 0); !errors.Is(err, ErrDivisionPorCero) {
		t.Errorf("esperaba ErrDivisionPorCero, obtuve: %v", err)
	}

	resultado, err := Dividir(10, 2)
	if err != nil {
		t.Fatalf("no esperaba error, obtuve: %v", err)
	}
	if resultado != 5 {
		t.Errorf("Dividir(10, 2) = %v; esperado 5", resultado)
	}
}

func TestEsPrimo(t *testing.T) {
	primos := map[int]bool{1: false, 2: true, 3: true, 4: false, 17: true, 100: false}
	for n, esperado := range primos {
		if got := EsPrimo(n); got != esperado {
			t.Errorf("EsPrimo(%d) = %v; esperado %v", n, got, esperado)
		}
	}
}

// BenchmarkEsPrimo usa b.Loop() (Go 1.24+): más simple y preciso que el
// patrón manual `for i := 0; i < b.N; i++`.
func BenchmarkEsPrimo(b *testing.B) {
	for b.Loop() {
		EsPrimo(7919)
	}
}
