package main

import (
	"database/sql"
	"errors"
	"testing"

	_ "modernc.org/sqlite"
)

func repoDePrueba(t *testing.T) *RepositorioTareas {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if err := migrar(db); err != nil {
		t.Fatal(err)
	}
	return &RepositorioTareas{db: db}
}

func TestRepositorioTareas_CrearYListar(t *testing.T) {
	repo := repoDePrueba(t)

	if _, err := repo.Crear("primera tarea"); err != nil {
		t.Fatalf("Crear() error = %v", err)
	}
	if _, err := repo.Crear("segunda tarea"); err != nil {
		t.Fatalf("Crear() error = %v", err)
	}

	tareas, err := repo.Listar()
	if err != nil {
		t.Fatalf("Listar() error = %v", err)
	}
	if len(tareas) != 2 {
		t.Fatalf("esperaba 2 tareas, obtuve %d", len(tareas))
	}
}

func TestRepositorioTareas_MarcarHecha(t *testing.T) {
	repo := repoDePrueba(t)

	tarea, err := repo.Crear("tarea a completar")
	if err != nil {
		t.Fatalf("Crear() error = %v", err)
	}

	if err := repo.MarcarHecha(tarea.ID); err != nil {
		t.Fatalf("MarcarHecha() error = %v", err)
	}

	tareas, _ := repo.Listar()
	if !tareas[0].Hecha {
		t.Error("esperaba que la tarea estuviera marcada como hecha")
	}
}

func TestRepositorioTareas_MarcarHechaInexistente(t *testing.T) {
	repo := repoDePrueba(t)

	err := repo.MarcarHecha(999)
	if !errors.Is(err, ErrNoEncontrada) {
		t.Errorf("esperaba ErrNoEncontrada, obtuve: %v", err)
	}
}
