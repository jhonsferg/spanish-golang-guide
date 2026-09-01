# Recursos adicionales

Enlaces oficiales y de la comunidad, curados - no una lista exhaustiva de
"Awesome X", sino lo que de verdad conviene tener a mano.

## Oficiales

- [go.dev](https://go.dev) - sitio oficial del lenguaje.
- [go.dev/doc/effective_go](https://go.dev/doc/effective_go) - la guía de
  estilo e idiomatismo escrita por el propio equipo de Go; sigue siendo
  vigente.
- [go.dev/doc/devel/release](https://go.dev/doc/devel/release) - índice de
  todas las release notes (la fuente de la sección
  [Versiones de Go](../10-versiones-de-go/index.md) de esta guía).
- [pkg.go.dev](https://pkg.go.dev) - documentación de cualquier paquete
  público, incluida la librería estándar.
- [go.dev/play](https://go.dev/play) - Go Playground, para probar snippets
  sin instalar nada.
- [go.dev/blog](https://go.dev/blog) - anuncios y artículos técnicos del
  equipo de Go.

## Filosofía y estilo

- [Go Proverbs](https://go-proverbs.github.io/) - los aforismos de Rob Pike
  que resumen la filosofía de diseño del lenguaje.
- [Google Go Style Guide](https://google.github.io/styleguide/go/) - guía
  de estilo usada internamente en Google, pública y muy referenciada.
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) -
  otra guía de estilo influyente, con más recomendaciones concretas de
  performance y concurrencia.

## Comunidad

- [Awesome Go](https://github.com/avelino/awesome-go) - listado curado de
  librerías y frameworks por categoría.
- [r/golang](https://www.reddit.com/r/golang/) - comunidad activa en
  Reddit.
- [Gophers Slack](https://invite.slack.golangbridge.org/) - Slack de la
  comunidad, con canales por tema (`#general`, `#performance`,
  `#kubernetes-golang`, etc.).

## Herramientas relacionadas con esta guía

- [staticcheck](https://staticcheck.io/) - linter estático más allá de
  `go vet`, muy usado en CI.
- [golangci-lint](https://golangci-lint.run/) - agregador de linters, útil
  para estandarizar checks de calidad en un proyecto.
- [Delve](https://github.com/go-delve/delve) - debugger para Go.

## Frameworks y librerías cubiertos en la guía

Si quieres profundizar más allá de lo que cubre cada capítulo, la
documentación oficial de cada proyecto:

- [Gin](https://gin-gonic.com/) · [Echo](https://echo.labstack.com/) ·
  [Fiber](https://gofiber.io/) · [Chi](https://go-chi.io/) - frameworks web,
  ver [Parte VIII](../08-frameworks-web/56-http-frameworks-overview.md).
- [GORM](https://gorm.io/) · [sqlc](https://sqlc.dev/) ·
  [Ent](https://entgo.io/) - acceso a datos, ver
  [Parte IX](../09-bases-de-datos-y-orms/61-gorm-orm.md).
