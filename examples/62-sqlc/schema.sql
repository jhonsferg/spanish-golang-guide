-- Esquema fuente para sqlc. En un proyecto real, `sqlc generate` lee este
-- archivo junto a query.sql y produce el código Go de db.go automáticamente.
CREATE TABLE autores (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    pais   TEXT NOT NULL
);
