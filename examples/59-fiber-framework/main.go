// Ejemplo del Capítulo 59: una API de comentarios con Fiber - la API
// inspirada en Express.js, construida sobre fasthttp en vez de net/http.
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Comentario struct {
	ID    int    `json:"id"`
	Autor string `json:"autor"`
	Texto string `json:"texto"`
}

func nuevaApp() *fiber.App {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	comentarios := map[int]Comentario{1: {1, "Ada", "primer comentario"}}
	siguienteID := 2

	api := app.Group("/api")

	api.Get("/comentarios", func(c *fiber.Ctx) error {
		lista := make([]Comentario, 0, len(comentarios))
		for _, com := range comentarios {
			lista = append(lista, com)
		}
		return c.JSON(lista)
	})

	api.Post("/comentarios", func(c *fiber.Ctx) error {
		var com Comentario
		if err := c.BodyParser(&com); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
		}
		if strings.TrimSpace(com.Texto) == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "texto es requerido"})
		}
		com.ID = siguienteID
		comentarios[com.ID] = com
		siguienteID++
		return c.Status(http.StatusCreated).JSON(com)
	})

	api.Get("/comentarios/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "id inválido"})
		}
		com, ok := comentarios[id]
		if !ok {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "no encontrado"})
		}
		return c.JSON(com)
	})

	return app
}

func probar(app *fiber.App, metodo, ruta, cuerpo string) *http.Response {
	var req *http.Request
	if cuerpo != "" {
		req, _ = http.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(metodo, ruta, nil)
	}
	resp, err := app.Test(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func main() {
	app := nuevaApp()

	resp := probar(app, http.MethodPost, "/api/comentarios", `{"autor":"Grace","texto":"aprendiendo Fiber"}`)
	fmt.Println("POST comentario ->", resp.StatusCode)

	respVacio := probar(app, http.MethodPost, "/api/comentarios", `{"autor":"Grace","texto":""}`)
	fmt.Println("POST comentario vacío ->", respVacio.StatusCode)

	respGet := probar(app, http.MethodGet, "/api/comentarios", "")
	fmt.Println("GET lista ->", respGet.StatusCode)

	respNotFound := probar(app, http.MethodGet, "/api/comentarios/999", "")
	fmt.Println("GET inexistente ->", respNotFound.StatusCode)
}
