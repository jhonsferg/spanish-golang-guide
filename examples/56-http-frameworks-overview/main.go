// Ejemplo del Capítulo 56: el mismo endpoint mínimo — GET /ping ->
// {"mensaje":"pong"} — implementado con los cuatro frameworks que
// cubre esta parte de la guía, para comparar su sintaxis lado a lado.
// Ver los capítulos 57-60 para un ejemplo más completo de cada uno.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
)

type respuestaPing struct {
	Mensaje string `json:"mensaje"`
}

func probarConGin() string {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, respuestaPing{Mensaje: "pong"})
	})
	return probarHandlerHTTP(r)
}

func probarConEcho() string {
	e := echo.New()
	e.HideBanner = true
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, respuestaPing{Mensaje: "pong"})
	})
	return probarHandlerHTTP(e)
}

func probarConChi() string {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respuestaPing{Mensaje: "pong"})
	})
	return probarHandlerHTTP(r)
}

// probarHandlerHTTP sirve para Gin/Echo/Chi: los tres implementan
// http.Handler, así que se pueden montar sobre el mismo net/http.Server.
func probarHandlerHTTP(h http.Handler) string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "error: " + err.Error()
	}
	server := &http.Server{Handler: h}
	go server.Serve(ln)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	resp, err := http.Get("http://" + ln.Addr().String() + "/ping")
	if err != nil {
		return "error: " + err.Error()
	}
	defer resp.Body.Close()
	var r respuestaPing
	json.NewDecoder(resp.Body).Decode(&r)
	return fmt.Sprintf("status=%d body=%+v", resp.StatusCode, r)
}

// Fiber no usa net/http por dentro (corre sobre fasthttp), así que se
// prueba con su propio método Test() en vez de un net.Listener.
func probarConFiber() string {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(respuestaPing{Mensaje: "pong"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	resp, err := app.Test(req)
	if err != nil {
		return "error: " + err.Error()
	}
	defer resp.Body.Close()
	var r respuestaPing
	json.NewDecoder(resp.Body).Decode(&r)
	return fmt.Sprintf("status=%d body=%+v", resp.StatusCode, r)
}

func main() {
	fmt.Println("Gin:   ", probarConGin())
	fmt.Println("Echo:  ", probarConEcho())
	fmt.Println("Fiber: ", probarConFiber())
	fmt.Println("Chi:   ", probarConChi())
}
