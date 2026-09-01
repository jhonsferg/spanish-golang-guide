// Ejemplo del Capítulo 6: if, switch y las variantes del for (el único
// loop que tiene Go).
package main

import "fmt"

func clasificar(n int) string {
	switch {
	case n < 0:
		return "negativo"
	case n == 0:
		return "cero"
	case n%2 == 0:
		return "positivo par"
	default:
		return "positivo impar"
	}
}

func main() {
	// if con inicializador.
	if n := 42; n > 0 {
		fmt.Println("42 es positivo")
	}

	for i := -2; i <= 3; i++ {
		fmt.Printf("clasificar(%d) = %s\n", i, clasificar(i))
	}

	// for como while.
	suma, i := 0, 1
	for i <= 5 {
		suma += i
		i++
	}
	fmt.Println("suma 1..5 =", suma)

	// range sobre un slice.
	frutas := []string{"manzana", "banana", "cereza"}
	for idx, fruta := range frutas {
		fmt.Printf("frutas[%d] = %s\n", idx, fruta)
	}

	// range sobre enteros (Go 1.22+).
	total := 0
	for i := range 5 {
		total += i
	}
	fmt.Println("suma de range 5 =", total)

	// switch con fallthrough explícito (nota: fallthrough NO evalúa la
	// condición del siguiente case, simplemente ejecuta su cuerpo).
	switch dia := "lunes"; dia {
	case "sábado", "domingo":
		fmt.Println("es fin de semana")
	case "lunes":
		fmt.Println("empieza la semana")
		fallthrough
	case "martes":
		fmt.Println("todavía es temprano en la semana")
	default:
		fmt.Println("día laboral cualquiera")
	}
}
