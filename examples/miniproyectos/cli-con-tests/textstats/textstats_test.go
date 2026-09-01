package textstats

import (
	"strings"
	"testing"
)

func TestContarPalabras(t *testing.T) {
	casos := []struct {
		nombre   string
		texto    string
		esperado map[string]int
	}{
		{"texto simple", "Go es simple. Go es rápido.", map[string]int{"go": 2, "es": 2, "simple": 1, "rápido": 1}},
		{"texto vacío", "", map[string]int{}},
		{"con números", "Go 1.27 es la versión 27", map[string]int{"go": 1, "1": 1, "27": 2, "es": 1, "la": 1, "versión": 1}},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := ContarPalabras(c.texto)
			if len(got) != len(c.esperado) {
				t.Fatalf("ContarPalabras(%q) = %v; esperaba %v", c.texto, got, c.esperado)
			}
			for palabra, n := range c.esperado {
				if got[palabra] != n {
					t.Errorf("conteo de %q = %d; esperaba %d", palabra, got[palabra], n)
				}
			}
		})
	}
}

func TestPalabraMasFrecuente(t *testing.T) {
	conteo := ContarPalabras("go go go es genial y es rápido")
	palabra, n := PalabraMasFrecuente(conteo)
	if palabra != "go" || n != 3 {
		t.Errorf("PalabraMasFrecuente() = (%q, %d); esperaba (\"go\", 3)", palabra, n)
	}
}

func TestPalabraMasFrecuente_MapaVacio(t *testing.T) {
	palabra, n := PalabraMasFrecuente(map[string]int{})
	if palabra != "" || n != -1 {
		t.Errorf("con mapa vacío esperaba (\"\", -1), obtuve (%q, %d)", palabra, n)
	}
}

func TestLongitudPromedio(t *testing.T) {
	conteo := ContarPalabras("ab ab cd") // "ab" x2 (len 2), "cd" x1 (len 2)
	got := LongitudPromedio(conteo)
	if got != 2.0 {
		t.Errorf("LongitudPromedio() = %v; esperaba 2.0", got)
	}
}

func BenchmarkContarPalabras(b *testing.B) {
	texto := strings.Repeat("el veloz murciélago hindú comía feliz cardillo y kiwi ", 100)
	for b.Loop() {
		ContarPalabras(texto)
	}
}
