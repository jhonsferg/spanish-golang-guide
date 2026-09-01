// Ejemplo del Capítulo 31: duraciones, formateo/parseo con el layout de
// referencia de Go, y aritmética de tiempos.
package main

import (
	"fmt"
	"time"
)

func main() {
	ahora := time.Date(2026, time.September, 1, 14, 30, 0, 0, time.UTC)

	// El layout de referencia de Go es una fecha específica: 2006-01-02 15:04:05.
	fmt.Println("formato RFC3339:", ahora.Format(time.RFC3339))
	fmt.Println("formato custom:  ", ahora.Format("02/01/2006 15:04"))
	fmt.Println("solo fecha:      ", ahora.Format(time.DateOnly)) // desde Go 1.20

	parseado, err := time.Parse("02/01/2006", "25/12/2026")
	if err != nil {
		fmt.Println("error al parsear:", err)
	} else {
		fmt.Println("parseado:", parseado.Format(time.RFC3339))
	}

	// Duraciones: literales y aritmética.
	d := 90 * time.Minute
	fmt.Println("90 minutos como duración:", d, "| en horas:", d.Hours())

	futuro := ahora.Add(48 * time.Hour)
	fmt.Println("48h después:", futuro.Format(time.DateTime))

	diferencia := futuro.Sub(ahora)
	fmt.Println("diferencia:", diferencia)

	fmt.Println("¿ahora es antes que futuro?", ahora.Before(futuro))

	// time.Since se usa mucho para medir cuánto tardó una operación.
	inicio := time.Now()
	time.Sleep(5 * time.Millisecond)
	fmt.Println("tiempo transcurrido:", time.Since(inicio) > 0)
}
