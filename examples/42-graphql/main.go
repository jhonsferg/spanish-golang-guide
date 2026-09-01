// Ejemplo del Capítulo 42: un schema GraphQL mínimo con
// github.com/graphql-go/graphql — un solo endpoint que responde
// exactamente los campos que el cliente pide, a diferencia de REST.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type libro struct {
	Titulo string
	Autor  string
	Anio   int
}

var libros = []libro{
	{"Cien años de soledad", "Gabriel García Márquez", 1967},
	{"Ficciones", "Jorge Luis Borges", 1944},
}

func main() {
	libroType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Libro",
		Fields: graphql.Fields{
			"titulo": &graphql.Field{Type: graphql.String},
			"autor":  &graphql.Field{Type: graphql.String},
			"anio":   &graphql.Field{Type: graphql.Int},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"libros": &graphql.Field{
				Type: graphql.NewList(libroType),
				Resolve: func(p graphql.ResolveParams) (any, error) {
					return libros, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: queryType})
	if err != nil {
		log.Fatal(err)
	}

	// El cliente pide solo "titulo" y "autor" — GraphQL no devuelve "anio"
	// aunque el resolver lo tenga disponible. Esa es la diferencia clave
	// frente a un endpoint REST de forma fija.
	consulta := `{ libros { titulo autor } }`

	resultado := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: consulta,
	})
	if len(resultado.Errors) > 0 {
		log.Fatalf("errores GraphQL: %v", resultado.Errors)
	}

	salida, _ := json.MarshalIndent(resultado.Data, "", "  ")
	fmt.Println(string(salida))
}
