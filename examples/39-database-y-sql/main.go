// Ejemplo del Capítulo 39: database/sql con un driver puro-Go de SQLite
// (modernc.org/sqlite, sin cgo) - placeholders parametrizados, no
// concatenación de strings, y manejo correcto de sql.ErrNoRows.
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // se registra como driver "sqlite" vía init()
)

type Tarea struct {
	ID          int
	Titulo      string
	Completada  bool
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(`
		CREATE TABLE tareas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			titulo TEXT NOT NULL,
			completada BOOLEAN NOT NULL DEFAULT 0
		)`); err != nil {
		log.Fatal(err)
	}

	// SIEMPRE con placeholders (?), nunca con fmt.Sprintf: evita SQL injection.
	titulos := []string{"Escribir la guía", "Revisar ejemplos", "Publicar el sitio"}
	for _, t := range titulos {
		if _, err := db.Exec("INSERT INTO tareas (titulo) VALUES (?)", t); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := db.Exec("UPDATE tareas SET completada = 1 WHERE titulo = ?", "Escribir la guía"); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, titulo, completada FROM tareas ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tareas []Tarea
	for rows.Next() {
		var t Tarea
		if err := rows.Scan(&t.ID, &t.Titulo, &t.Completada); err != nil {
			log.Fatal(err)
		}
		tareas = append(tareas, t)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, t := range tareas {
		fmt.Printf("[%d] %-20s completada=%v\n", t.ID, t.Titulo, t.Completada)
	}

	// QueryRow + sql.ErrNoRows: el patrón correcto para "cero o un resultado".
	var titulo string
	err = db.QueryRow("SELECT titulo FROM tareas WHERE id = ?", 999).Scan(&titulo)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		fmt.Println("tarea 999 no existe (esperado)")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Println("título:", titulo)
	}
}
