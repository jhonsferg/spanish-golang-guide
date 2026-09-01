// Ejemplo del Capítulo 62: usar el código "generado por sqlc" (ver
// sqlcgen/db.go) - funciones típadas por consulta, sin ORM, con el SQL
// crudo visible y versionado en query.sql.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jhonsferg/spanish-golang-guide/examples/62-sqlc/sqlcgen"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.ExecContext(ctx, `CREATE TABLE autores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL,
		pais TEXT NOT NULL
	)`); err != nil {
		log.Fatal(err)
	}

	queries := sqlcgen.New(db)

	a1, err := queries.CrearAutor(ctx, sqlcgen.CrearAutorParams{Nombre: "Gabriel García Márquez", Pais: "Colombia"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("autor creado: %+v\n", a1)

	if _, err := queries.CrearAutor(ctx, sqlcgen.CrearAutorParams{Nombre: "Isabel Allende", Pais: "Chile"}); err != nil {
		log.Fatal(err)
	}

	autor, err := queries.ObtenerAutor(ctx, a1.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("autor obtenido por ID: %+v\n", autor)

	todos, err := queries.ListarAutores(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("todos los autores, ordenados por nombre:")
	for _, a := range todos {
		fmt.Printf("  %s (%s)\n", a.Nombre, a.Pais)
	}
}
