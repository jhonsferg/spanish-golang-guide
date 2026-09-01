# Guía de Go en Español

Una guía completa de Go (Golang), en español, pensada para acompañarte desde tu
primer `go run main.go` hasta diseñar microservicios, elegir un framework web o
un ORM en producción.

No es una traducción de la documentación oficial: es un recorrido con
explicaciones extendidas, comparativas, ejercicios y mini-proyectos, organizado
como un libro de 66 capítulos agrupados en 9 partes.

## ¿Para quién es esta guía?

- **Si vienes de otro lenguaje**: los primeros capítulos explican no solo la
  sintaxis de Go sino *por qué* Go decide las cosas distinto (sin herencia, con
  goroutines en vez de threads, con `error` en vez de excepciones).
- **Si ya escribes Go**: las partes de concurrencia, librería estándar,
  arquitectura y frameworks están pensadas para profundizar y comparar
  alternativas reales (Gin vs Echo vs Fiber vs Chi, GORM vs sqlc vs Ent).
- **Si necesitas resolver algo puntual**: usa la búsqueda del sitio o el
  [glosario](11-recursos/glosario.md).

## Cómo está organizada

| Parte | Contenido | Capítulos |
|---|---|---|
| I. Fundamentos | Sintaxis, tipos básicos, control de flujo, funciones, colecciones | 1–10 |
| II. Tipos y composición | Structs, métodos, interfaces, embedding, punteros, errores | 11–19 |
| III. Concurrencia | Goroutines, channels, select, sync, patrones | 20–25 |
| IV. Librería estándar | io, os, strings, json, net, time, math, reflect, sort, hash | 26–35 |
| V. Producción y herramientas | Testing, logging, HTTP avanzado, SQL, build y performance | 36–40 |
| VI. Arquitectura y sistemas distribuidos | Microservicios, GraphQL, WebSockets, patrones de diseño, Docker, Kubernetes, observabilidad | 41–53 |
| VII. Proyectos integrados | Dos proyectos completos de punta a punta | 54–55 |
| VIII. Frameworks web | Gin, Echo, Fiber, Chi | 56–60 |
| IX. Bases de datos y ORMs | GORM, sqlc, Ent, drivers, MongoDB, caching | 61–66 |

Además:

- **[Versiones de Go](10-versiones-de-go/index.md)** — qué trajo cada versión
  desde Go 1.18 (generics) hasta Go 1.27, la más reciente.
- **[Recursos](11-recursos/buenas-practicas.md)** — buenas prácticas
  consolidadas, glosario y enlaces oficiales curados.

## Ejemplos ejecutables

Cada capítulo con código enlaza a un ejemplo real y verificable en
[`examples/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples),
compilable con:

```bash
cd examples
go build ./...
go vet ./...
```

## Por dónde empezar

¿Primera vez con Go? Empieza por
[Capítulo 1 — Introducción a Go](01-fundamentos/01-introduccion-a-go.md).

¿Ya conoces lo básico? Salta directo a la parte que necesites en la tabla de
arriba, o revisa qué trajo [Go 1.27](10-versiones-de-go/go-1-27.md).
