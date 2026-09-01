// Ejemplo del Capítulo 33: introspección básica con reflect — cómo
// bibliotecas como encoding/json inspeccionan structs en runtime.
package main

import (
	"fmt"
	"reflect"
)

type Producto struct {
	Nombre string `validar:"requerido"`
	Precio float64 `validar:"positivo"`
	Stock  int
}

// describirCampos usa reflection para listar nombre, tipo y tag de
// cada campo de CUALQUIER struct que se le pase.
func describirCampos(v any) {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	if t.Kind() != reflect.Struct {
		fmt.Println("no es un struct:", t.Kind())
		return
	}

	for i := 0; i < t.NumField(); i++ {
		campo := t.Field(i)
		valor := val.Field(i)
		tag := campo.Tag.Get("validar")
		fmt.Printf("  %-8s tipo=%-8s valor=%-10v tag=%q\n",
			campo.Name, campo.Type, valor, tag)
	}
}

func main() {
	p := Producto{Nombre: "Teclado", Precio: 49.99, Stock: 12}

	fmt.Println("campos de Producto:")
	describirCampos(p)

	// reflect.ValueOf + Kind() para inspeccionar cualquier valor genérico.
	valores := []any{42, "texto", 3.14, true}
	for _, v := range valores {
		fmt.Printf("valor=%v kind=%v\n", v, reflect.ValueOf(v).Kind())
	}
}
