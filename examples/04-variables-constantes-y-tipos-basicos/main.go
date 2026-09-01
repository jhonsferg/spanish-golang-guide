// Ejemplo del Capítulo 4: declaración de variables, constantes,
// inferencia de tipos e iota.
package main

import "fmt"

type ByteSize float64

const (
	_           = iota // salta el 0
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
)

func main() {
	// Declaración explícita vs inferencia de tipos.
	var edad int = 30
	nombre := "Ada" // tipo inferido: string
	var activo bool

	// Múltiples variables en una sola declaración.
	var x, y, z = 1, 2.5, "tres"

	// Constantes tipadas y sin tipar.
	const pi = 3.14159
	const maxIntentos int = 5

	fmt.Printf("nombre=%s edad=%d activo=%v\n", nombre, edad, activo)
	fmt.Printf("x=%v (%T) y=%v (%T) z=%v (%T)\n", x, x, y, y, z, z)
	fmt.Printf("pi=%v maxIntentos=%v\n", pi, maxIntentos)

	fmt.Printf("1 KB = %.0f bytes\n", float64(KB))
	fmt.Printf("1 MB = %.0f bytes\n", float64(MB))
	fmt.Printf("1 GB = %.0f bytes\n", float64(GB))
}
