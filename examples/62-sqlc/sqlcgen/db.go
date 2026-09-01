// Package sqlcgen reproduce a mano, con fines didácticos, la FORMA del
// código que `sqlc generate` produciría automáticamente a partir de
// ../schema.sql y ../query.sql: una struct Queries con un método típado
// por cada consulta con nombre, sin ORM ni reflection de por medio -
// solo database/sql con Scan explícito. En un proyecto real este
// archivo NO se escribe a mano, sqlc lo regenera en cada build.
package sqlcgen

import (
	"context"
	"database/sql"
)

type Autor struct {
	ID     int64
	Nombre string
	Pais   string
}

type Queries struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

const crearAutor = `INSERT INTO autores (nombre, pais) VALUES (?, ?) RETURNING id, nombre, pais`

type CrearAutorParams struct {
	Nombre string
	Pais   string
}

func (q *Queries) CrearAutor(ctx context.Context, arg CrearAutorParams) (Autor, error) {
	row := q.db.QueryRowContext(ctx, crearAutor, arg.Nombre, arg.Pais)
	var a Autor
	err := row.Scan(&a.ID, &a.Nombre, &a.Pais)
	return a, err
}

const obtenerAutor = `SELECT id, nombre, pais FROM autores WHERE id = ?`

func (q *Queries) ObtenerAutor(ctx context.Context, id int64) (Autor, error) {
	row := q.db.QueryRowContext(ctx, obtenerAutor, id)
	var a Autor
	err := row.Scan(&a.ID, &a.Nombre, &a.Pais)
	return a, err
}

const listarAutores = `SELECT id, nombre, pais FROM autores ORDER BY nombre`

func (q *Queries) ListarAutores(ctx context.Context) ([]Autor, error) {
	rows, err := q.db.QueryContext(ctx, listarAutores)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Autor
	for rows.Next() {
		var a Autor
		if err := rows.Scan(&a.ID, &a.Nombre, &a.Pais); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}
