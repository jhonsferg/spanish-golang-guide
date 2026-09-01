// Package echoapi: misma API de TODOs que ginapi, fiberapi y chiapi,
// implementada con Echo.
package echoapi

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
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

func NewRouter() http.Handler {
	e := echo.New()
	e.HideBanner = true

	store := &almacen{datos: make(map[int]Todo), siguienteID: 1}

	e.GET("/todos", func(c echo.Context) error {
		store.mu.Lock()
		defer store.mu.Unlock()
		lista := make([]Todo, 0, len(store.datos))
		for _, t := range store.datos {
			lista = append(lista, t)
		}
		return c.JSON(http.StatusOK, lista)
	})

	e.POST("/todos", func(c echo.Context) error {
		var t Todo
		if err := c.Bind(&t); err != nil || t.Titulo == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "titulo es requerido"})
		}
		store.mu.Lock()
		t.ID = store.siguienteID
		store.datos[t.ID] = t
		store.siguienteID++
		store.mu.Unlock()
		return c.JSON(http.StatusCreated, t)
	})

	e.PUT("/todos/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		store.mu.Lock()
		defer store.mu.Unlock()
		t, ok := store.datos[id]
		if !ok {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "no encontrado"})
		}
		t.Hecho = true
		store.datos[id] = t
		return c.JSON(http.StatusOK, t)
	})

	e.DELETE("/todos/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		store.mu.Lock()
		delete(store.datos, id)
		store.mu.Unlock()
		return c.NoContent(http.StatusNoContent)
	})

	return e
}
