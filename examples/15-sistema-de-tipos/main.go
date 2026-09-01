// Ejemplo del Capítulo 15: tipos definidos (named types), type aliases y
// type assertions.
package main

import "fmt"

// Tipo definido: Celsius es un tipo NUEVO basado en float64, no
// intercambiable directamente con float64 sin conversión explícita.
type Celsius float64
type Fahrenheit float64

func (c Celsius) AFahrenheit() Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

// Type alias: EnteroGrande es exactamente el mismo tipo que int64,
// intercambiable sin conversión.
type EnteroGrande = int64

func main() {
	agua := Celsius(100)
	fmt.Printf("%.0f°C = %.0f°F\n", float64(agua), float64(agua.AFahrenheit()))

	var x EnteroGrande = 10
	var y int64 = x // sin conversión, porque es un alias, no un tipo nuevo
	fmt.Println("alias sin conversión:", y)

	// Type assertion sobre una interfaz.
	var i any = "texto"
	s, ok := i.(string)
	fmt.Println("assertion exitosa:", s, ok)

	n, ok := i.(int)
	fmt.Println("assertion fallida, valor cero:", n, ok)
}
