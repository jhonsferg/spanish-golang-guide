// Ejemplo del Capítulo 20: importar y usar un paquete propio.
// Ver examples/20-paquetes/saludo/saludo.go para el paquete importado.
package main

import (
	"fmt"

	"github.com/jhonsferg/spanish-golang-guide/examples/20-paquetes/saludo"
)

func main() {
	fmt.Println(saludo.Saludar("Gopher"))
	fmt.Println(saludo.SaludoFormal("Gopher"))

	// saludo.formatearPrivado no compilaría aquí: no está exportada.
}
