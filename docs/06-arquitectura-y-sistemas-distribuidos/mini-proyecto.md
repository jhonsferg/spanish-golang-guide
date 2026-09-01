# Mini-proyecto: API con Clean Architecture

Una API HTTP completa (CRUD de notas) organizada en capas — dominio,
casos de uso y repositorio — donde el transporte HTTP es apenas un
adaptador reemplazable. Expande la idea del
[Capítulo 45](45-clean-architecture.md) a un CRUD real, con Dockerfile
de despliegue incluido.

## Qué demuestra

- **El caso de uso (`ServicioNotas`) depende de una interfaz**
  (`RepositorioNotas`), no de una implementación — cambiar de memoria a
  PostgreSQL más adelante no tocaría ni una línea de `ServicioNotas`.
- **Las reglas de negocio viven en el caso de uso, no en el handler**:
  la validación "el contenido no puede estar vacío" está en
  `CrearNota`, no en el handler HTTP — así que también aplicaría si
  mañana esta misma lógica se expone por gRPC o por un worker de cola.
- **El handler HTTP solo traduce**: decodifica JSON, llama al caso de
  uso, traduce el resultado a un status code. Cero lógica de negocio.
- Un `Dockerfile` multi-stage listo para producción (build sin cgo,
  imagen final `distroless`).

## Diagrama de dependencias

```
HTTP handler  --depende de-->  ServicioNotas (caso de uso)
                                     |
                                     | depende de la INTERFAZ
                                     v
                              RepositorioNotas
                                     ^
                                     | implementa
                                     |
                            repoEnMemoria (adaptador)
```

Las flechas de dependencia siempre apuntan hacia el dominio — nunca al
revés. Eso es Clean Architecture en una frase.

## Ejecutarlo

```bash
cd examples/miniproyectos/mini-api-clean-architecture
go run .
```

Código fuente: [`examples/miniproyectos/mini-api-clean-architecture/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/miniproyectos/mini-api-clean-architecture)
