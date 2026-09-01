-- Cada comentario -- name: define una función que sqlc generará en db.go,
-- con el tipo de retorno indicado (:one, :many, :exec).

-- name: CrearAutor :one
INSERT INTO autores (nombre, pais) VALUES (?, ?) RETURNING *;

-- name: ObtenerAutor :one
SELECT * FROM autores WHERE id = ?;

-- name: ListarAutores :many
SELECT * FROM autores ORDER BY nombre;
