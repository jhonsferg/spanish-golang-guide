// Ejemplo del Capítulo 65: el driver oficial de MongoDB — documentos con
// bson, filtros, inserción y consulta. A diferencia de los ejemplos
// anteriores no hay una versión "en memoria" de MongoDB: si no hay un
// servidor corriendo en localhost:27017, el ejemplo lo detecta con un
// timeout corto y lo indica claramente en vez de fallar de forma confusa.
package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Articulo struct {
	Titulo    string   `bson:"titulo"`
	Etiquetas []string `bson:"etiquetas"`
	Vistas    int      `bson:"vistas"`
}

func main() {
	ctxConexion, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	cliente, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("no se pudo crear el cliente de Mongo:", err)
		return
	}
	defer cliente.Disconnect(context.Background())

	if err := cliente.Ping(ctxConexion, nil); err != nil {
		fmt.Println("no hay un servidor MongoDB disponible en localhost:27017.")
		fmt.Println("este ejemplo muestra la API real del driver; para ejecutarlo")
		fmt.Println("de punta a punta, levanta uno con: docker run -p 27017:27017 mongo")
		return
	}

	ctx := context.Background()
	coleccion := cliente.Database("guia_go_demo").Collection("articulos")
	defer coleccion.Drop(ctx) // limpieza al terminar el ejemplo

	articulos := []any{
		Articulo{Titulo: "Introducción a Go", Etiquetas: []string{"go", "básico"}, Vistas: 150},
		Articulo{Titulo: "Concurrencia con goroutines", Etiquetas: []string{"go", "concurrencia"}, Vistas: 320},
	}
	if _, err := coleccion.InsertMany(ctx, articulos); err != nil {
		fmt.Println("error insertando:", err)
		return
	}

	// bson.M{} es un filtro tipo mapa; aquí buscamos por etiqueta y
	// ordenamos por vistas descendente.
	cursor, err := coleccion.Find(ctx,
		bson.M{"etiquetas": "go"},
		options.Find().SetSort(bson.D{{Key: "vistas", Value: -1}}),
	)
	if err != nil {
		fmt.Println("error consultando:", err)
		return
	}
	defer cursor.Close(ctx)

	var resultados []Articulo
	if err := cursor.All(ctx, &resultados); err != nil {
		fmt.Println("error leyendo resultados:", err)
		return
	}

	fmt.Println("artículos con etiqueta 'go', por vistas:")
	for _, a := range resultados {
		fmt.Printf("  %-30s vistas=%d\n", a.Titulo, a.Vistas)
	}
}
