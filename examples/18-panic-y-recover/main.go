// Ejemplo del Capítulo 18: panic/recover, y por qué son para errores
// irrecuperables - no para control de flujo habitual (para eso están
// los errores como valores, ver Capítulo 17).
package main

import "fmt"

// dividirSeguro recupera de un panic (p.ej. división por cero en enteros)
// y lo convierte en un error normal.
func dividirSeguro(a, b int) (resultado int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recuperado de panic: %v", r)
		}
	}()
	resultado = a / b // división entera por cero hace panic en Go
	return resultado, nil
}

func procesarLista(items []int, idx int) (valor int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("índice fuera de rango recuperado: %v", r)
		}
	}()
	return items[idx], nil
}

func main() {
	if r, err := dividirSeguro(10, 2); err == nil {
		fmt.Println("10 / 2 =", r)
	}

	if _, err := dividirSeguro(10, 0); err != nil {
		fmt.Println("error:", err)
	}

	items := []int{1, 2, 3}
	if _, err := procesarLista(items, 10); err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("el programa sigue vivo después de recuperar los panics")
}
