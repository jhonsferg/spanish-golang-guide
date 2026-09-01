// Ejemplo del Capítulo 17: errores como valores, wrapping con %w,
// errors.Is/As y errores personalizados.
package main

import (
	"errors"
	"fmt"
)

// ErrNoEncontrado es un error "sentinel": se compara por identidad con
// errors.Is, no por el texto del mensaje.
var ErrNoEncontrado = errors.New("recurso no encontrado")

// ErrorValidacion es un error personalizado con datos estructurados.
type ErrorValidacion struct {
	Campo string
	Razon string
}

func (e *ErrorValidacion) Error() string {
	return fmt.Sprintf("validación fallida en %q: %s", e.Campo, e.Razon)
}

func buscarUsuario(id int) error {
	if id <= 0 {
		return &ErrorValidacion{Campo: "id", Razon: "debe ser positivo"}
	}
	if id > 1000 {
		// %w envuelve el error, preservando la cadena para errors.Is/As.
		return fmt.Errorf("buscando usuario %d: %w", id, ErrNoEncontrado)
	}
	return nil
}

func main() {
	for _, id := range []int{-1, 5000, 42} {
		err := buscarUsuario(id)
		if err == nil {
			fmt.Printf("usuario %d: encontrado\n", id)
			continue
		}

		var errVal *ErrorValidacion
		switch {
		case errors.Is(err, ErrNoEncontrado):
			fmt.Printf("usuario %d: no existe (%v)\n", id, err)
		case errors.As(err, &errVal):
			fmt.Printf("usuario %d: error de validación en campo %q\n", id, errVal.Campo)
		default:
			fmt.Printf("usuario %d: error inesperado: %v\n", id, err)
		}
	}
}
