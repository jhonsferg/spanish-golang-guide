# Guía de Go en Español

> Guía completa de Go (Golang) en español: 66 capítulos organizados en 9 partes,
> desde "Hola, mundo" hasta microservicios, frameworks web y ORMs en producción.
> Con ejemplos de código verificables, mini-proyectos y una sección dedicada a
> qué trajo cada versión de Go (Go 1.18 → 1.27).
>
> 🌐 Navegable como sitio en **[jhonsferg.github.io/spanish-golang-guide](https://jhonsferg.github.io/spanish-golang-guide/)** o como Markdown plano directamente en `docs/`.

---

## ¿Qué es esta guía?

Esta guía está escrita para desarrolladores hispanohablantes que quieren dominar Go de verdad: no solo aprender la sintaxis, sino entender el _por qué_ detrás de cada decisión de diseño del lenguaje. Cada concepto se explica como lo haría un senior a un colega: directo, honesto y con ejemplos que resuelven problemas reales.

No es un tutorial de "hola mundo" ni una traducción de la documentación oficial. Es una guía de referencia progresiva que construye conocimiento de forma acumulativa, desde los fundamentos del lenguaje hasta patrones de arquitectura, concurrencia, producción y sistemas distribuidos.

---

## Versiones cubiertas

| Versión | Fecha    | Lo que se cubre en esta guía                                        |
|---------|----------|---------------------------------------------------------------------|
| Go 1.18 | Mar 2022 | Generics, fuzzing nativo, `any`, workspace mode                     |
| Go 1.19 | Ago 2022 | Mejoras en docs, `atomic.Int64/Uint64`, soft memory limit           |
| Go 1.20 | Feb 2023 | `errors.Join`, `http.ResponseController`, slices helpers            |
| Go 1.21 | Ago 2023 | `slices`, `maps`, `cmp`, `log/slog`, `min/max` builtins             |
| Go 1.22 | Feb 2024 | Variables de loop por iteración, `math/rand/v2`, routing mejorado   |
| Go 1.23 | Ago 2024 | Iteradores con `range` sobre funciones, `unique` package            |
| Go 1.24 | Feb 2025 | `generic type aliases`, `weak` package, mejoras en testing          |
| Go 1.25 | Ago 2025 | Swiss tables map, mejoras en el GC, `sync.Map` genérico             |
| Go 1.26 | Feb 2026 | `os/user` portable, mejoras en `net/http`, `crypto` modernizado     |
| Go 1.27 | Ago 2026 | Profile-guided optimization estable, `encoding/json/v2` (preview)  |

---

## Estructura del repositorio

```
docs/                    # Contenido de la guía (fuente del sitio MkDocs)
  01-fundamentos/
  02-tipos-y-composicion/
  03-concurrencia/
  04-libreria-estandar/
  05-produccion-y-herramientas/
  06-arquitectura-y-sistemas-distribuidos/
  07-proyectos-integrados/
  08-frameworks-web/
  09-bases-de-datos-y-orms/
  10-versiones-de-go/      # Changelog por versión de Go
  11-recursos/             # Buenas prácticas, glosario, enlaces
examples/                 # Código Go ejecutable (un módulo por el repo)
  <NN>-<capitulo>/        # Ejemplo mínimo por capítulo
  miniproyectos/          # Mini-proyectos por parte
mkdocs.yml               # Configuración del sitio
pyproject.toml           # Dependencias Python (uv)
```

---

## Levantar el sitio localmente

Requiere [uv](https://docs.astral.sh/uv/):

```bash
# Instalar uv (si no lo tienes)
curl -LsSf https://astral.sh/uv/install.sh | sh

# Instalar dependencias y levantar el servidor
uv sync
uv run mkdocs serve
# → http://127.0.0.1:8000
```

Alternativamente, con Docker (no requiere Python):

```bash
docker compose up
# → http://localhost:8080
```

---

## Ejecutar y verificar los ejemplos de código

```bash
cd examples
go build ./...
go vet ./...
```

---

## Contribuir

Las contribuciones son bienvenidas. Lee la [Guía de contribución](CONTRIBUTING.md) para saber cómo empezar.

---

## Licencia

El contenido de la guía (`docs/`) y los ejemplos de código (`examples/`) se distribuyen bajo [CC BY-SA 4.0](LICENSE).
