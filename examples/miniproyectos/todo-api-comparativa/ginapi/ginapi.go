// Package ginapi implementa la misma API de TODOs que echoapi, fiberapi
// y chiapi - para comparar la sintaxis de cada framework sobre el
// idéntico caso de uso. Ver el mini-proyecto de la Parte VIII.
package ginapi

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo" binding:"required"`
	Hecho  bool   `json:"hecho"`
}

type almacen struct {
	mu          sync.Mutex
	datos       map[int]Todo
	siguienteID int
}

func NewRouter() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	store := &almacen{datos: make(map[int]Todo), siguienteID: 1}

	r.GET("/todos", func(c *gin.Context) {
		store.mu.Lock()
		defer store.mu.Unlock()
		lista := make([]Todo, 0, len(store.datos))
		for _, t := range store.datos {
			lista = append(lista, t)
		}
		c.JSON(http.StatusOK, lista)
	})

	r.POST("/todos", func(c *gin.Context) {
		var t Todo
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		store.mu.Lock()
		t.ID = store.siguienteID
		store.datos[t.ID] = t
		store.siguienteID++
		store.mu.Unlock()
		c.JSON(http.StatusCreated, t)
	})

	r.PUT("/todos/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		store.mu.Lock()
		defer store.mu.Unlock()
		t, ok := store.datos[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"})
			return
		}
		t.Hecho = true
		store.datos[id] = t
		c.JSON(http.StatusOK, t)
	})

	r.DELETE("/todos/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		store.mu.Lock()
		delete(store.datos, id)
		store.mu.Unlock()
		c.Status(http.StatusNoContent)
	})

	return r
}
