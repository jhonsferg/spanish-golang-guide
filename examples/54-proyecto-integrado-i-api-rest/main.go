// Ejemplo del Capítulo 54 (versión condensada del proyecto integrado):
// una API REST de tareas con persistencia real en SQLite, capa de
// repositorio separada del handler HTTP, y autenticación simple por
// API key — las piezas mínimas de una API en producción, en un archivo
// para que sea fácil de leer de punta a punta.
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	_ "modernc.org/sqlite"
)

func jsonBody(s string) io.Reader { return strings.NewReader(s) }

// --- Modelo ---

type Tarea struct {
	ID     int64  `json:"id"`
	Titulo string `json:"titulo"`
	Hecha  bool   `json:"hecha"`
}

// --- Repositorio: aísla el resto de la app de SQL concreto ---

type RepositorioTareas struct{ db *sql.DB }

func (r *RepositorioTareas) Crear(titulo string) (Tarea, error) {
	res, err := r.db.Exec("INSERT INTO tareas (titulo, hecha) VALUES (?, 0)", titulo)
	if err != nil {
		return Tarea{}, err
	}
	id, _ := res.LastInsertId()
	return Tarea{ID: id, Titulo: titulo}, nil
}

func (r *RepositorioTareas) Listar() ([]Tarea, error) {
	rows, err := r.db.Query("SELECT id, titulo, hecha FROM tareas ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tareas []Tarea
	for rows.Next() {
		var t Tarea
		if err := rows.Scan(&t.ID, &t.Titulo, &t.Hecha); err != nil {
			return nil, err
		}
		tareas = append(tareas, t)
	}
	return tareas, rows.Err()
}

var ErrNoEncontrada = errors.New("tarea no encontrada")

func (r *RepositorioTareas) MarcarHecha(id int64) error {
	res, err := r.db.Exec("UPDATE tareas SET hecha = 1 WHERE id = ?", id)
	if err != nil {
		return err
	}
	filas, _ := res.RowsAffected()
	if filas == 0 {
		return ErrNoEncontrada
	}
	return nil
}

// --- Middleware de autenticación por API key ---

const apiKeyValida = "clave-de-ejemplo-123"

func requerirAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != apiKeyValida {
			http.Error(w, `{"error":"API key inválida o ausente"}`, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- Handlers HTTP ---

func nuevoMux(repo *RepositorioTareas) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tareas", func(w http.ResponseWriter, r *http.Request) {
		var entrada struct {
			Titulo string `json:"titulo"`
		}
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil || entrada.Titulo == "" {
			http.Error(w, `{"error":"titulo es requerido"}`, http.StatusBadRequest)
			return
		}
		tarea, err := repo.Crear(entrada.Titulo)
		if err != nil {
			http.Error(w, `{"error":"error interno"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tarea)
	})

	mux.HandleFunc("GET /tareas", func(w http.ResponseWriter, r *http.Request) {
		tareas, err := repo.Listar()
		if err != nil {
			http.Error(w, `{"error":"error interno"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tareas)
	})

	mux.HandleFunc("PATCH /tareas/{id}/completar", func(w http.ResponseWriter, r *http.Request) {
		var id int64
		fmt.Sscanf(r.PathValue("id"), "%d", &id)
		if err := repo.MarcarHecha(id); err != nil {
			if errors.Is(err, ErrNoEncontrada) {
				http.Error(w, `{"error":"tarea no encontrada"}`, http.StatusNotFound)
				return
			}
			http.Error(w, `{"error":"error interno"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return requerirAPIKey(mux)
}

func migrar(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE tareas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		titulo TEXT NOT NULL,
		hecha BOOLEAN NOT NULL DEFAULT 0
	)`)
	return err
}

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := migrar(db); err != nil {
		log.Fatal(err)
	}

	repo := &RepositorioTareas{db: db}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Handler: nuevoMux(repo)}
	go server.Serve(ln)
	defer server.Close()

	base := "http://" + ln.Addr().String()
	cliente := &http.Client{}

	hacer := func(metodo, ruta, cuerpo string) *http.Response {
		var r *http.Request
		if cuerpo != "" {
			r, _ = http.NewRequest(metodo, base+ruta, jsonBody(cuerpo))
		} else {
			r, _ = http.NewRequest(metodo, base+ruta, nil)
		}
		r.Header.Set("X-API-Key", apiKeyValida)
		resp, err := cliente.Do(r)
		if err != nil {
			log.Fatal(err)
		}
		return resp
	}

	resp := hacer("POST", "/tareas", `{"titulo":"Escribir el proyecto integrado"}`)
	fmt.Println("POST /tareas ->", resp.StatusCode)
	resp.Body.Close()

	resp2 := hacer("PATCH", "/tareas/1/completar", "")
	fmt.Println("PATCH /tareas/1/completar ->", resp2.StatusCode)
	resp2.Body.Close()

	resp3 := hacer("GET", "/tareas", "")
	var tareas []Tarea
	json.NewDecoder(resp3.Body).Decode(&tareas)
	resp3.Body.Close()
	fmt.Printf("GET /tareas -> %+v\n", tareas)

	respSinKey, _ := http.Get(base + "/tareas")
	fmt.Println("GET /tareas sin API key ->", respSinKey.StatusCode)
	respSinKey.Body.Close()
}
