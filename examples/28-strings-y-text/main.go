// Ejemplo del Capítulo 28: las funciones de manipulación de texto más
// usadas de `strings`.
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	frase := "  Go es simple, rápido y confiable  "

	limpia := strings.TrimSpace(frase)
	fmt.Printf("original: %q\n", frase)
	fmt.Printf("trim:     %q\n", limpia)

	fmt.Println("mayúsculas:", strings.ToUpper(limpia))
	fmt.Println("contiene 'rápido':", strings.Contains(limpia, "rápido"))
	fmt.Println("empieza con 'Go':", strings.HasPrefix(limpia, "Go"))
	fmt.Println("reemplazado:", strings.ReplaceAll(limpia, "simple", "directo"))

	palabras := strings.Fields(limpia) // separa por cualquier whitespace
	fmt.Println("palabras:", palabras, "| total:", len(palabras))

	unidas := strings.Join(palabras, "_")
	fmt.Println("unidas con _:", unidas)

	partes := strings.Split("a,b,,c", ",")
	fmt.Printf("split de \"a,b,,c\": %q (nota el elemento vacío)\n", partes)

	// strconv: el puente entre strings y tipos numéricos.
	n, err := strconv.Atoi("42")
	if err == nil {
		fmt.Println("convertido a int:", n+8)
	}
	fmt.Println("de vuelta a string:", strconv.Itoa(100))

	if _, err := strconv.Atoi("no-es-un-número"); err != nil {
		fmt.Println("error esperado de conversión:", err)
	}
}
