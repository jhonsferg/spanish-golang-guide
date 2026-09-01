// Ejemplo del Capítulo 58: una API de notas con Echo - bind + validación
// manual, grupos de rutas y middleware estándar (Recover, request ID).
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Nota struct {
	ID       int    `json:"id"`
	Contenido string `json:"contenido"`
}

func (n Nota) Validar() error {
	if strings.TrimSpace(n.Contenido) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "contenido es requerido")
	}
	return nil
}

func nuevoEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover(), middleware.RequestID())

	notas := map[int]Nota{1: {1, "primera nota"}}
	siguienteID := 2

	api := e.Group("/api")

	api.GET("/notas", func(c echo.Context) error {
		lista := make([]Nota, 0, len(notas))
		for _, n := range notas {
			lista = append(lista, n)
		}
		return c.JSON(http.StatusOK, lista)
	})

	api.POST("/notas", func(c echo.Context) error {
		var n Nota
		if err := c.Bind(&n); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "JSON inválido"})
		}
		if err := n.Validar(); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
		n.ID = siguienteID
		notas[n.ID] = n
		siguienteID++
		return c.JSON(http.StatusCreated, n)
	})

	api.GET("/notas/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "id inválido"})
		}
		n, ok := notas[id]
		if !ok {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "no encontrada"})
		}
		return c.JSON(http.StatusOK, n)
	})

	return e
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	e := nuevoEcho()
	server := &http.Server{Handler: e}
	go server.Serve(ln)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	base := "http://" + ln.Addr().String()
	cliente := &http.Client{}

	resp, _ := cliente.Post(base+"/api/notas", "application/json", strings.NewReader(`{"contenido":"aprender Echo"}`))
	fmt.Println("POST nota ->", resp.StatusCode)
	resp.Body.Close()

	respVacia, _ := cliente.Post(base+"/api/notas", "application/json", strings.NewReader(`{"contenido":""}`))
	fmt.Println("POST nota vacía ->", respVacia.StatusCode)
	respVacia.Body.Close()

	respGet, _ := cliente.Get(base + "/api/notas/1")
	fmt.Println("GET /api/notas/1 ->", respGet.StatusCode)
	respGet.Body.Close()
}
