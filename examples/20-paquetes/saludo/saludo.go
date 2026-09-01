// Package saludo es un paquete de apoyo para el ejemplo del Capítulo 20:
// cómo se organiza y se importa código propio entre paquetes.
package saludo

import "fmt"

// Saludar está exportada (empieza con mayúscula): visible fuera del paquete.
func Saludar(nombre string) string {
	return fmt.Sprintf("¡Hola, %s! (desde el paquete saludo)", nombre)
}

// formatearPrivado NO está exportada: solo visible dentro de este paquete.
func formatearPrivado(s string) string {
	return "[" + s + "]"
}

// SaludoFormal usa la función privada del propio paquete.
func SaludoFormal(nombre string) string {
	return formatearPrivado(Saludar(nombre))
}
