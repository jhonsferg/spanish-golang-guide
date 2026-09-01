// Ejemplo del Capítulo 10: por qué len("café") no es lo que esperas,
// y cómo iterar texto correctamente en presencia de UTF-8.
package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	texto := "café ☕"

	fmt.Println("len() en bytes:", len(texto))
	fmt.Println("RuneCountInString:", utf8.RuneCountInString(texto))

	// Iterar con range da (índice en bytes, rune) — NO (índice, byte).
	fmt.Println("--- iterando con range (runas) ---")
	for i, r := range texto {
		fmt.Printf("byte %d: %q (%d bytes)\n", i, r, utf8.RuneLen(r))
	}

	fmt.Println("--- iterando bytes crudos ---")
	for i := 0; i < len(texto); i++ {
		fmt.Printf("byte %d: 0x%02x\n", i, texto[i])
	}

	// strings.Builder: forma eficiente de construir strings sin
	// reallocaciones repetidas por concatenación con +.
	var sb strings.Builder
	for _, palabra := range []string{"Go", "es", "rápido"} {
		sb.WriteString(palabra)
		sb.WriteByte(' ')
	}
	fmt.Println("construido:", strings.TrimSpace(sb.String()))

	// Conversión []rune para operar por carácter (p.ej. invertir texto).
	runas := []rune("café")
	for i, j := 0, len(runas)-1; i < j; i, j = i+1, j-1 {
		runas[i], runas[j] = runas[j], runas[i]
	}
	fmt.Println("invertido:", string(runas))
}
