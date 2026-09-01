// Mini-proyecto de la Parte I (Fundamentos): un conversor de unidades de
// línea de comandos. Usa solo lo visto en los capítulos 1-10: variables,
// constantes, funciones, control de flujo, maps y strings — sin structs
// ni interfaces todavía (eso llega en la Parte II).
package main

import (
	"fmt"
	"strconv"
	"strings"
)

type conversion struct {
	desde, hacia string
	factor       func(v float64) float64
}

var tablaLongitud = []conversion{
	{"km", "mi", func(v float64) float64 { return v * 0.621371 }},
	{"mi", "km", func(v float64) float64 { return v / 0.621371 }},
	{"m", "ft", func(v float64) float64 { return v * 3.28084 }},
	{"ft", "m", func(v float64) float64 { return v / 3.28084 }},
}

var tablaPeso = []conversion{
	{"kg", "lb", func(v float64) float64 { return v * 2.20462 }},
	{"lb", "kg", func(v float64) float64 { return v / 2.20462 }},
}

func convertirTemperatura(desde, hacia string, v float64) (float64, bool) {
	switch {
	case desde == "c" && hacia == "f":
		return v*9/5 + 32, true
	case desde == "f" && hacia == "c":
		return (v - 32) * 5 / 9, true
	case desde == hacia:
		return v, true
	default:
		return 0, false
	}
}

// convertir busca en las tablas de longitud y peso, y cae a temperatura
// como último recurso (es el único caso que necesita fórmula, no factor).
func convertir(desde, hacia string, valor float64) (float64, error) {
	for _, tabla := range [][]conversion{tablaLongitud, tablaPeso} {
		for _, c := range tabla {
			if c.desde == desde && c.hacia == hacia {
				return c.factor(valor), nil
			}
		}
	}
	if resultado, ok := convertirTemperatura(desde, hacia, valor); ok {
		return resultado, nil
	}
	return 0, fmt.Errorf("no sé convertir de %q a %q", desde, hacia)
}

// procesarComando parsea líneas con el formato "10 km a mi".
func procesarComando(linea string) (string, error) {
	partes := strings.Fields(linea)
	if len(partes) != 4 || partes[2] != "a" {
		return "", fmt.Errorf("formato esperado: \"<valor> <unidad> a <unidad>\", recibí: %q", linea)
	}

	valor, err := strconv.ParseFloat(partes[0], 64)
	if err != nil {
		return "", fmt.Errorf("valor inválido %q: %w", partes[0], err)
	}

	resultado, err := convertir(partes[1], partes[3], valor)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s = %.4f %s", partes[0], partes[1], resultado, partes[3]), nil
}

func main() {
	comandos := []string{
		"10 km a mi",
		"5 lb a kg",
		"100 c a f",
		"32 f a c",
		"3 km a lb", // combinación inválida, a propósito
	}

	for _, cmd := range comandos {
		resultado, err := procesarComando(cmd)
		if err != nil {
			fmt.Printf("%-15s -> error: %v\n", cmd, err)
			continue
		}
		fmt.Printf("%-15s -> %s\n", cmd, resultado)
	}
}
