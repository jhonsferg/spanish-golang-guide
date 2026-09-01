// Ejemplo del Capítulo 66: una interfaz Cache pequeña con dos
// implementaciones — Redis (go-redis) para producción, y en memoria para
// tests o para cuando no hay Redis disponible. El resto de la app solo
// depende de la interfaz, nunca de un cliente concreto.
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, clave, valor string, ttl time.Duration) error
	Get(ctx context.Context, clave string) (string, bool, error)
}

// --- Implementación en memoria ---

type CacheEnMemoria struct {
	mu   sync.Mutex
	data map[string]entrada
}

type entrada struct {
	valor     string
	expiraEn  time.Time
}

func NuevaCacheEnMemoria() *CacheEnMemoria {
	return &CacheEnMemoria{data: make(map[string]entrada)}
}

func (c *CacheEnMemoria) Set(ctx context.Context, clave, valor string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[clave] = entrada{valor: valor, expiraEn: time.Now().Add(ttl)}
	return nil
}

func (c *CacheEnMemoria) Get(ctx context.Context, clave string) (string, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.data[clave]
	if !ok || time.Now().After(e.expiraEn) {
		return "", false, nil
	}
	return e.valor, true, nil
}

// --- Implementación con Redis ---

type CacheRedis struct {
	cliente *redis.Client
}

func (c *CacheRedis) Set(ctx context.Context, clave, valor string, ttl time.Duration) error {
	return c.cliente.Set(ctx, clave, valor, ttl).Err()
}

func (c *CacheRedis) Get(ctx context.Context, clave string) (string, bool, error) {
	valor, err := c.cliente.Get(ctx, clave).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return valor, true, nil
}

// demostrar ejercita cualquier implementación de Cache de forma idéntica.
func demostrar(nombre string, c Cache) {
	ctx := context.Background()
	c.Set(ctx, "usuario:1:nombre", "Ada Lovelace", time.Minute)

	valor, ok, err := c.Get(ctx, "usuario:1:nombre")
	if err != nil {
		fmt.Printf("[%s] error: %v\n", nombre, err)
		return
	}
	fmt.Printf("[%s] usuario:1:nombre = %q (presente=%v)\n", nombre, valor, ok)

	_, ok, _ = c.Get(ctx, "clave-inexistente")
	fmt.Printf("[%s] clave-inexistente presente=%v\n", nombre, ok)
}

func main() {
	demostrar("memoria", NuevaCacheEnMemoria())

	cliente := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer cliente.Close()

	ctxPing, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()

	if err := cliente.Ping(ctxPing).Err(); err != nil {
		fmt.Println("no hay Redis disponible en localhost:6379; para probar esta")
		fmt.Println("implementación real, levanta uno con: docker run -p 6379:6379 redis")
		return
	}

	demostrar("redis", &CacheRedis{cliente: cliente})
}
