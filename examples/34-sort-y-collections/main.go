// Ejemplo del Capítulo 34: ordenar slices con sort.Slice y con el
// paquete genérico `slices` (Go 1.21+).
package main

import (
	"fmt"
	"slices"
	"sort"
)

type Empleado struct {
	Nombre  string
	Salario int
}

func main() {
	empleados := []Empleado{
		{"Carla", 3200},
		{"Beto", 2800},
		{"Ana", 4100},
	}

	// sort.Slice: ordenamiento con función de comparación custom.
	sort.Slice(empleados, func(i, j int) bool {
		return empleados[i].Salario > empleados[j].Salario
	})
	fmt.Println("empleados por salario descendente:")
	for _, e := range empleados {
		fmt.Printf("  %s: %d\n", e.Nombre, e.Salario)
	}

	// El paquete genérico `slices` para tipos simples: más directo que sort.Ints.
	numeros := []int{5, 2, 8, 1, 9, 3}
	slices.Sort(numeros)
	fmt.Println("números ordenados:", numeros)

	fmt.Println("¿está ordenado?", slices.IsSorted(numeros))
	fmt.Println("índice de 8:", slices.Index(numeros, 8))
	fmt.Println("contiene 100:", slices.Contains(numeros, 100))

	// slices.SortFunc para tipos custom, alternativa moderna a sort.Slice.
	slices.SortFunc(empleados, func(a, b Empleado) int {
		return len(a.Nombre) - len(b.Nombre)
	})
	fmt.Println("empleados por longitud de nombre:", empleados)
}
