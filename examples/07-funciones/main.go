// Ejemplo del Capítulo 7: múltiples valores de retorno, funciones
// variádicas, funciones anónimas/closures y defer.
package main

import (
	"errors"
	"fmt"
)

// dividir retorna dos valores: el resultado y un posible error.
// Este es EL patrón idiomático de Go para manejo de errores.
func dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("división por cero")
	}
	return a / b, nil
}

// sumar es variádica: acepta cero o más enteros.
func sumar(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// contador retorna una closure que recuerda su propio estado entre llamadas.
func contador() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

func main() {
	if resultado, err := dividir(10, 2); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("10 / 2 =", resultado)
	}

	if _, err := dividir(10, 0); err != nil {
		fmt.Println("error esperado:", err)
	}

	fmt.Println("sumar(1,2,3) =", sumar(1, 2, 3))
	fmt.Println("sumar() =", sumar())

	siguiente := contador()
	fmt.Println(siguiente(), siguiente(), siguiente())

	defer fmt.Println("esto se imprime al final, por el defer")
	fmt.Println("esto se imprime primero")
}
