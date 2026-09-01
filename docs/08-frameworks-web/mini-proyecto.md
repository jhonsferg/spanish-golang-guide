# Mini-proyecto: la misma API de TODOs en Gin, Echo, Fiber y Chi

La misma API (crear, listar, completar, borrar tareas) implementada
**cuatro veces**, una por framework, para comparar directamente cómo
resuelve cada uno el mismo problema - la mejor forma de decidir cuál
usar en un proyecto real, más allá de benchmarks de "hola mundo".

## Qué comparar al leer el código

| Aspecto | Gin | Echo | Fiber | Chi |
|---|---|---|---|---|
| Binding + validación | `ShouldBindJSON` + tag `binding` | `Bind` manual | `BodyParser` manual | `json.Decode` manual (net/http puro) |
| Respuestas JSON | `c.JSON(status, gin.H{...})` | `c.JSON(status, echo.Map{...})` | `c.JSON(...)` / `c.Status(...).JSON(...)` | `json.NewEncoder(w).Encode(...)` a mano |
| Parámetros de ruta | `c.Param("id")` | `c.Param("id")` | `c.Params("id")` | `chi.URLParam(r, "id")` |
| Motor HTTP | net/http | net/http | fasthttp (no net/http) | net/http |
| Testing sin puerto real | requiere `httptest`/listener | requiere `httptest`/listener | `app.Test(req)` incluido | requiere `httptest`/listener |

La fila más importante es la última: **Fiber no es un `http.Handler`**
porque corre sobre `fasthttp`, no sobre `net/http` - por eso
`main.go` lo prueba distinto al resto (`app.Test()` en vez de un
`net.Listener` + `http.Server`). Es la decisión de diseño que más
implicaciones tiene si necesitas integrar con middleware o librerías del
ecosistema `net/http` estándar.

## Ejecutarlo

```bash
cd examples/miniproyectos/todo-api-comparativa
go run .
```

Código fuente: [`examples/miniproyectos/todo-api-comparativa/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/miniproyectos/todo-api-comparativa)
(carpetas `ginapi/`, `echoapi/`, `fiberapi/`, `chiapi/`).
