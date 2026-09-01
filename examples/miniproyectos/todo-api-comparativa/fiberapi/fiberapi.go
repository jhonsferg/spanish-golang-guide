// Package fiberapi: misma API de TODOs que ginapi, echoapi y chiapi,
// implementada con Fiber.
package fiberapi

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
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

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	store := &almacen{datos: make(map[int]Todo), siguienteID: 1}

	app.Get("/todos", func(c *fiber.Ctx) error {
		store.mu.Lock()
		defer store.mu.Unlock()
		lista := make([]Todo, 0, len(store.datos))
		for _, t := range store.datos {
			lista = append(lista, t)
		}
		return c.JSON(lista)
	})

	app.Post("/todos", func(c *fiber.Ctx) error {
		var t Todo
		if err := c.BodyParser(&t); err != nil || t.Titulo == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "titulo es requerido"})
		}
		store.mu.Lock()
		t.ID = store.siguienteID
		store.datos[t.ID] = t
		store.siguienteID++
		store.mu.Unlock()
		return c.Status(http.StatusCreated).JSON(t)
	})

	app.Put("/todos/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		store.mu.Lock()
		defer store.mu.Unlock()
		t, ok := store.datos[id]
		if !ok {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "no encontrado"})
		}
		t.Hecho = true
		store.datos[id] = t
		return c.JSON(t)
	})

	app.Delete("/todos/:id", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		store.mu.Lock()
		delete(store.datos, id)
		store.mu.Unlock()
		return c.SendStatus(http.StatusNoContent)
	})

	return app
}
