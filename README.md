# Guía de Go en Español

Guía completa de Go (Golang) en español: 66 capítulos organizados en 9 partes,
desde "Hola, mundo" hasta microservicios, frameworks web y ORMs en producción,
con ejemplos de código verificables, mini-proyectos y una sección dedicada a
qué trajo cada versión de Go (hasta Go 1.27).

Sitio publicado: ver `docs/` renderizado con [mkdocs-material](https://squidfunk.github.io/mkdocs-material/).

## Estructura del repositorio

```
docs/                    # Contenido de la guía (fuente del sitio mkdocs)
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
  11-recursos/              # Buenas prácticas, glosario, enlaces
examples/                 # Código Go ejecutable, un módulo por el repo entero
  <NN>-<capitulo>/         # Ejemplo mínimo por capítulo
  miniproyectos/           # Mini-proyectos por parte
mkdocs.yml                 # Configuración del sitio
```

Cada archivo dentro de `docs/` conserva su número de capítulo global (p.ej.
`docs/03-concurrencia/21-goroutines.md` es el Capítulo 21), porque el texto se
refiere entre capítulos por número.

## Levantar el sitio localmente

Requiere Python 3.

```bash
pip install -r requirements.txt
mkdocs serve
```

Abre `http://127.0.0.1:8000`.

## Ejecutar y verificar los ejemplos de código

```bash
cd examples
go build ./...
go vet ./...
```

## Licencia

El contenido de la guía (`docs/`) y los ejemplos de código (`examples/`) se
distribuyen bajo [CC BY-SA 4.0](LICENSE).

## Contribuir

Este repositorio sigue reglas estrictas de commits y ramas descritas en
`.claude/GIT_RULES.md` (commits atómicos por archivo, ramas por mejora, sin
push directo a `main`). Si vas a contribuir, léelas antes de abrir un PR.
