// Ejemplo del Capítulo 8: arrays de tamaño fijo vs. slices dinámicos,
// append, copy y el peligro de compartir el array subyacente.
package main

import "fmt"

func main() {
	// Array: tamaño parte del tipo, valor (no referencia).
	var arr [3]int = [3]int{1, 2, 3}
	arrCopia := arr
	arrCopia[0] = 99
	fmt.Println("array original:", arr, "| copia modificada:", arrCopia)

	// Slice: vista dinámica sobre un array subyacente.
	slice := []int{10, 20, 30}
	fmt.Println("slice:", slice, "len:", len(slice), "cap:", cap(slice))

	slice = append(slice, 40, 50)
	fmt.Println("tras append:", slice, "len:", len(slice), "cap:", cap(slice))

	// Dos slices pueden compartir el mismo array subyacente.
	base := []int{1, 2, 3, 4, 5}
	sub := base[1:3] // [2 3], comparte memoria con base
	sub[0] = 999
	fmt.Println("base tras modificar sub:", base)

	// copy() para desacoplar realmente los datos.
	destino := make([]int, len(base))
	n := copy(destino, base)
	destino[0] = -1
	fmt.Printf("copiados: %d | base: %v | destino: %v\n", n, base, destino)

	// Slice de slices (matriz).
	matriz := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	for _, fila := range matriz {
		fmt.Println("fila:", fila)
	}
}
