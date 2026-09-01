// Mini-proyecto de la Parte V (Producción y herramientas): un CLI de
// estadísticas de texto, con la lógica separada en el paquete
// textstats/ (ver textstats_test.go para los tests de tabla y el
// benchmark — corre `go test -bench=. -cpuprofile=cpu.prof` desde este
// directorio para perfilarlo con pprof).
package main

import (
	"fmt"

	"github.com/jhonsferg/spanish-golang-guide/examples/miniproyectos/cli-con-tests/textstats"
)

func main() {
	texto := `Go es un lenguaje de programación compilado, con tipado
	estático y concurrencia nativa. Go fue diseñado en Google para
	resolver problemas reales de infraestructura a gran escala.`

	conteo := textstats.ContarPalabras(texto)
	palabra, n := textstats.PalabraMasFrecuente(conteo)
	promedio := textstats.LongitudPromedio(conteo)

	fmt.Printf("palabras únicas:        %d\n", len(conteo))
	fmt.Printf("palabra más frecuente:  %q (%d veces)\n", palabra, n)
	fmt.Printf("longitud promedio:      %.2f caracteres\n", promedio)
}
