// Ejemplo del Capítulo 32: funciones de math y generación de números
// aleatorios con math/rand/v2 (Go 1.22+).
package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

func main() {
	fmt.Println("--- math ---")
	fmt.Printf("Sqrt(2) = %.4f\n", math.Sqrt(2))
	fmt.Printf("Pow(2, 10) = %.0f\n", math.Pow(2, 10))
	fmt.Printf("Max(3, 7) = %.0f | Min(3, 7) = %.0f\n", math.Max(3, 7), math.Min(3, 7))
	fmt.Printf("Ceil(4.2) = %.0f | Floor(4.8) = %.0f | Round(4.5) = %.0f\n",
		math.Ceil(4.2), math.Floor(4.8), math.Round(4.5))
	fmt.Printf("Abs(-9.5) = %.1f\n", math.Abs(-9.5))

	fmt.Println("--- math/rand/v2 (sin necesidad de sembrar manualmente) ---")
	fmt.Println("entero en [0, 100):", rand.IntN(100))
	fmt.Println("float64 en [0, 1):", rand.Float64() < 1.0)

	dados := make([]int, 5)
	for i := range dados {
		dados[i] = rand.IntN(6) + 1 // 1..6
	}
	fmt.Println("5 tiradas de dado:", dados)

	// Generador determinista con semilla fija, útil para tests reproducibles.
	fuente := rand.NewPCG(42, 1024)
	determinista := rand.New(fuente)
	fmt.Println("valor determinista con semilla fija:", determinista.IntN(1000) >= 0)
}
