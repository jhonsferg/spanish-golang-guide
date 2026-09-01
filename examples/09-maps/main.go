// Ejemplo del Capítulo 9: creación, el "comma ok idiom", borrado y por
// qué el orden de iteración de un map NO está garantizado.
package main

import (
	"fmt"
	"sort"
)

func main() {
	inventario := map[string]int{
		"manzanas": 10,
		"peras":    5,
	}

	inventario["naranjas"] = 8
	inventario["manzanas"]++

	// El "comma ok idiom": distingue "existe con valor cero" de "no existe".
	if cantidad, existe := inventario["peras"]; existe {
		fmt.Printf("hay %d peras\n", cantidad)
	}
	if _, existe := inventario["sandías"]; !existe {
		fmt.Println("no hay sandías registradas")
	}

	delete(inventario, "peras")

	// Para iterar en orden determinista, ordena las claves explícitamente:
	// el orden de `for k := range map` varía entre ejecuciones a propósito.
	claves := make([]string, 0, len(inventario))
	for k := range inventario {
		claves = append(claves, k)
	}
	sort.Strings(claves)

	for _, k := range claves {
		fmt.Printf("%s: %d\n", k, inventario[k])
	}

	fmt.Println("total de productos:", len(inventario))
}
