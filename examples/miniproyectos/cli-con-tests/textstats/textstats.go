// Package textstats es la lógica de negocio del mini-proyecto de la
// Parte V (Producción y herramientas) - separada de main.go a propósito,
// para que sea testeable sin tocar entrada/salida de consola.
package textstats

import (
	"strings"
	"unicode"
)

// ContarPalabras normaliza (minúsculas, sin puntuación) y cuenta
// ocurrencias de cada palabra.
func ContarPalabras(texto string) map[string]int {
	conteo := make(map[string]int)
	for _, palabra := range strings.FieldsFunc(texto, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) {
		conteo[strings.ToLower(palabra)]++
	}
	return conteo
}

// PalabraMasFrecuente retorna la palabra con más ocurrencias; en caso de
// empate, la que aparece primero alfabéticamente (para ser determinista).
func PalabraMasFrecuente(conteo map[string]int) (string, int) {
	mejorPalabra, mejorConteo := "", -1
	for palabra, n := range conteo {
		if n > mejorConteo || (n == mejorConteo && palabra < mejorPalabra) {
			mejorPalabra, mejorConteo = palabra, n
		}
	}
	return mejorPalabra, mejorConteo
}

// LongitudPromedio calcula el promedio de caracteres por palabra.
func LongitudPromedio(conteo map[string]int) float64 {
	totalPalabras, totalCaracteres := 0, 0
	for palabra, n := range conteo {
		totalPalabras += n
		totalCaracteres += len([]rune(palabra)) * n
	}
	if totalPalabras == 0 {
		return 0
	}
	return float64(totalCaracteres) / float64(totalPalabras)
}
