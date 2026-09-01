# Mini-proyecto: sistema de inventario

Un inventario con productos físicos (con stock) y digitales (sin límite
de stock), que comparten comportamiento vía interfaces pequeñas y
embedding de structs - el patrón central de la [Parte II](13-interfaces.md).

## Qué demuestra

- **Embedding**: `ProductoFisico` y `ProductoDigital` embeben
  `productoBase` para heredar `Precio()` sin reescribirlo.
- **Interfaces pequeñas y compuestas**: `ConStock` extiende `Vendible`
  agregando `Reservar`/`StockDisponible` - solo los productos físicos la
  implementan.
- **Type assertions** para pedir "más capacidad" a un valor que solo se
  conoce por su interfaz base (`p.(ConStock)`), el patrón detrás de cómo
  Go maneja capacidades opcionales sin herencia.
- Errores envueltos con contexto (`fmt.Errorf(...%w..., ErrStockInsuficiente)`).

## Por qué esta separación importa

`Inventario` solo conoce `Vendible` - no le importa si un producto es
físico o digital para calcular el valor total. Pero `procesarVenta`
necesita explícitamente `ConStock`, porque reservar stock no tiene
sentido para un producto digital. Esa es la diferencia entre "lo que
todos comparten" y "lo que algunos pueden hacer además".

## Ejecutarlo

```bash
cd examples/miniproyectos/inventario
go run .
```

Código fuente: [`examples/miniproyectos/inventario/`](https://github.com/jhonsferg/spanish-golang-guide/tree/main/examples/miniproyectos/inventario)
