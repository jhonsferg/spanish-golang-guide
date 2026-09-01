// Package entrepo recrea, a mano y en miniatura, el ESTILO de la API
// fluida que Ent genera automáticamente (ver Capítulo 63 para la
// explicación completa de por qué no es código realmente generado en
// este ejemplo). Sirve para comparar la ergonomía del patrón
// Client.Producto.Create()... frente a GORM y sqlc con el mismo modelo.
package entrepo

import "sort"

type Producto struct {
	ID     int
	Nombre string
	Precio float64
}

type Client struct {
	productos   map[int]Producto
	siguienteID int
}

func Nuevo() *Client {
	return &Client{productos: make(map[int]Producto)}
}

type ProductoCreate struct {
	client *Client
	nombre string
	precio float64
}

func (c *Client) CreateProducto() *ProductoCreate {
	return &ProductoCreate{client: c}
}

func (pc *ProductoCreate) SetNombre(n string) *ProductoCreate {
	pc.nombre = n
	return pc
}

func (pc *ProductoCreate) SetPrecio(p float64) *ProductoCreate {
	pc.precio = p
	return pc
}

func (pc *ProductoCreate) Save() (Producto, error) {
	pc.client.siguienteID++
	p := Producto{ID: pc.client.siguienteID, Nombre: pc.nombre, Precio: pc.precio}
	pc.client.productos[p.ID] = p
	return p, nil
}

func (c *Client) QueryProductos() []Producto {
	lista := make([]Producto, 0, len(c.productos))
	for _, p := range c.productos {
		lista = append(lista, p)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].Precio > lista[j].Precio })
	return lista
}
