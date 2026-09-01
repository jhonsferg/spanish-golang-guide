// Package chiapi: misma API de TODOs que ginapi, echoapi y fiberapi,
// implementada con Chi (handlers 100% net/http estándar).
package chiapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Todo struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Hecho  bool   `json:"hecho"`
}

type almacen struct {
	mu          sync.Mutex
	datos       map[int]Todo
	siguienteID int
}

func responder(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func NewRouter() http.Handler {
	r := chi.NewRouter()
	store := &almacen{datos: make(map[int]Todo), siguienteID: 1}

	r.Get("/todos", func(w http.ResponseWriter, r *http.Request) {
		store.mu.Lock()
		defer store.mu.Unlock()
		lista := make([]Todo, 0, len(store.datos))
		for _, t := range store.datos {
			lista = append(lista, t)
		}
		responder(w, http.StatusOK, lista)
	})

	r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
		var t Todo
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil || t.Titulo == "" {
			responder(w, http.StatusBadRequest, map[string]string{"error": "titulo es requerido"})
			return
		}
		store.mu.Lock()
		t.ID = store.siguienteID
		store.datos[t.ID] = t
		store.siguienteID++
		store.mu.Unlock()
		responder(w, http.StatusCreated, t)
	})

	r.Put("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		store.mu.Lock()
		defer store.mu.Unlock()
		t, ok := store.datos[id]
		if !ok {
			responder(w, http.StatusNotFound, map[string]string{"error": "no encontrado"})
			return
		}
		t.Hecho = true
		store.datos[id] = t
		responder(w, http.StatusOK, t)
	})

	r.Delete("/todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		store.mu.Lock()
		delete(store.datos, id)
		store.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	})

	return r
}
