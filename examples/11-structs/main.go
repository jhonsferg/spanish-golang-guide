// Ejemplo del Capítulo 11: structs, structs anónimos y structs embebidos.
package main

import "fmt"

type Direccion struct {
	Ciudad string
	Pais   string
}

type Persona struct {
	Nombre    string
	Edad      int
	Direccion // embebido: promueve Ciudad y Pais a Persona
}

func main() {
	p := Persona{
		Nombre: "Ada Lovelace",
		Edad:   36,
		Direccion: Direccion{
			Ciudad: "Londres",
			Pais:   "Reino Unido",
		},
	}

	// Acceso directo a campos promovidos por embedding.
	fmt.Printf("%s vive en %s, %s\n", p.Nombre, p.Ciudad, p.Pais)

	// Struct anónimo: útil para datos de un solo uso.
	config := struct {
		Debug   bool
		Puerto  int
	}{
		Debug:  true,
		Puerto: 8080,
	}
	fmt.Printf("config: %+v\n", config)

	// Comparación de structs: campo a campo, si todos los campos son comparables.
	otra := p
	otra.Direccion.Ciudad = "Manchester"
	fmt.Println("¿mismas direcciones?", p.Direccion == otra.Direccion)
}
