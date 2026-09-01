// Ejemplo del Capítulo 57: una API de productos con Gin — binding y
// validación automática de JSON, grupos de rutas y middleware.
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func bodyJSON(s string) io.Reader { return strings.NewReader(s) }

type Producto struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre" binding:"required"`
	Precio float64 `json:"precio" binding:"required,gt=0"`
}

func middlewareLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		inicio := time.Now()
		c.Next()
		fmt.Printf("[gin] %s %s -> %d (%v)\n", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(inicio))
	}
}

func nuevoRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), middlewareLogging())

	productos := map[int]Producto{1: {1, "Monitor", 199.99}}
	siguienteID := 2

	api := r.Group("/api/v1")
	{
		api.GET("/productos", func(c *gin.Context) {
			lista := make([]Producto, 0, len(productos))
			for _, p := range productos {
				lista = append(lista, p)
			}
			c.JSON(http.StatusOK, lista)
		})

		api.POST("/productos", func(c *gin.Context) {
			var p Producto
			// ShouldBindJSON valida automáticamente las tags `binding`.
			if err := c.ShouldBindJSON(&p); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			p.ID = siguienteID
			productos[p.ID] = p
			siguienteID++
			c.JSON(http.StatusCreated, p)
		})

		api.GET("/productos/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
				return
			}
			p, ok := productos[id]
			if !ok {
				c.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"})
				return
			}
			c.JSON(http.StatusOK, p)
		})
	}

	return r
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server := &http.Server{Handler: nuevoRouter()}
	go server.Serve(ln)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	base := "http://" + ln.Addr().String()
	cliente := &http.Client{}

	resp, _ := cliente.Post(base+"/api/v1/productos", "application/json",
		bodyJSON(`{"nombre":"Teclado","precio":49.9}`))
	fmt.Println("POST producto ->", resp.StatusCode)
	resp.Body.Close()

	respInvalido, _ := cliente.Post(base+"/api/v1/productos", "application/json",
		bodyJSON(`{"nombre":"Sin precio"}`))
	fmt.Println("POST producto inválido ->", respInvalido.StatusCode)
	respInvalido.Body.Close()

	respLista, _ := cliente.Get(base + "/api/v1/productos")
	fmt.Println("GET lista ->", respLista.StatusCode)
	respLista.Body.Close()
}
