// Ejemplo del Capítulo 3: estructura mínima de un programa Go y las
// piezas básicas de entrada/salida.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hola, mundo desde Go 1.27")

	nombre := leerNombre(os.Stdin)
	if nombre == "" {
		nombre = "desconocido"
	}
	fmt.Printf("Encantado de conocerte, %s.\n", nombre)
}

// leerNombre lee una sola línea de entrada; si no hay entrada disponible
// (como al correr el binario sin stdin interactivo), retorna "".
func leerNombre(r *os.File) string {
	scanner := bufio.NewScanner(r)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}
