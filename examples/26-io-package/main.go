// Ejemplo del Capítulo 26: las interfaces io.Reader / io.Writer y por
// qué son la base de casi toda la librería estándar de I/O.
package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// contarBytes acepta CUALQUIER io.Reader: un string, un archivo, una
// conexión de red, un buffer... todos implementan la misma interfaz.
func contarBytes(r io.Reader) (int, error) {
	n, err := io.Copy(io.Discard, r)
	return int(n), err
}

func main() {
	origen := strings.NewReader("Go trata la I/O como interfaces pequeñas y componibles.")

	n, err := contarBytes(origen)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("bytes leídos:", n)

	// bytes.Buffer implementa tanto io.Reader como io.Writer a la vez.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "línea 1\n")
	fmt.Fprintf(&buf, "línea 2\n")

	// io.Copy mueve datos de un Reader a un Writer sin que ninguno de
	// los dos sepa qué hay del otro lado.
	if _, err := io.Copy(&buf, strings.NewReader("línea 3\n")); err != nil {
		fmt.Println("error copiando:", err)
	}

	fmt.Print(buf.String())

	// io.MultiReader encadena varios readers como si fueran uno solo.
	combinado := io.MultiReader(
		strings.NewReader("A"),
		strings.NewReader("B"),
		strings.NewReader("C"),
	)
	datos, _ := io.ReadAll(combinado)
	fmt.Println("combinado:", string(datos))
}
