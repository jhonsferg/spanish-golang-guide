# Mini-proyecto: el mismo modelo con GORM, sqlc y Ent

El mismo modelo — un `Producto` con nombre y precio — creado y
consultado con tres enfoques de acceso a datos distintos, para comparar
directamente cuánto código y qué forma toma cada uno resolviendo
exactamente lo mismo.

## Qué comparar al leer el código

- **`gormrepo/`** (GORM): el modelo ES la definición de la tabla
  (`gorm.Model` + tags). `AutoMigrate` crea el schema desde el struct.
  Cero SQL visible.
- **`sqlcrepo/`** (estilo sqlc): el schema vive en SQL explícito, y cada
  función mapea 1:1 a una consulta con `Scan()` manual. Más líneas, cero
  magia — lo que ves es exactamente lo que se ejecuta.
- **`entrepo/`** (estilo Ent): una API fluida por *builders*
  (`Create().SetX().Save()`), la que Ent generaría automáticamente a
  partir de un schema declarativo — aquí escrita a mano solo para
  ilustrar la forma (ver [Capítulo 63](63-ent-orm.md) para la aclaración
  completa).

## La pregunta que este mini-proyecto ayuda a responder

¿Prefieres que el ORM **infiera** el schema desde tus structs (GORM),
que **generes código** desde SQL explícito (sqlc), o que **generes
código** desde un schema declarativo en Go (Ent)? Los tres llegan al
mismo resultado con el mismo modelo — la diferencia es dónde vive la
fuente de verdad y cuánto SQL terminas escribiendo o leyendo a mano.

## Una nota técnica: por qué `sqlcrepo` no importa su propio driver

`gormrepo` y `sqlcrepo` usan el mismo motor SQLite puro-Go por debajo
(`glebarez/go-sqlite`, que también usa `modernc.org/sqlite` internamente).
Un driver de `database/sql` solo puede registrarse una vez por nombre en
todo el binario — así que `sqlcrepo` reutiliza el driver `"sqlite"` que
ya registra `gormrepo` en vez de duplicar la importación. Ver el
comentario en `sqlcrepo/sqlcrepo.go`.

## Ejecutarlo

```bash
cd examples/miniproyectos/mismo-modelo-tres-orms
go run .
```

Código fuente: [`examples/miniproyectos/mismo-modelo-tres-orms/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/miniproyectos/mismo-modelo-tres-orms)
(carpetas `gormrepo/`, `sqlcrepo/`, `entrepo/`).
