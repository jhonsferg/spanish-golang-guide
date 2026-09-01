// Ejemplo del Capítulo 63: una recreación MINIATURA, escrita a mano, del
// estilo de API fluida que `entc generate` produce automáticamente para
// Ent a partir de un schema declarativo (ent/schema/persona.go en un
// proyecto real). El código generado real de Ent es mucho más extenso
// (miles de líneas por entidad); esto ilustra el PATRÓN —
// Client.Persona.Create().Set...().Save(ctx) y builders de consulta
// encadenables — sin requerir el paso de generación de código.
package main

import (
	"context"
	"errors"
	"fmt"
)

type Persona struct {
	ID     int
	Nombre string
	Edad   int
}

// Client es el análogo simplificado del *ent.Client generado.
type Client struct {
	Persona *PersonaClient
}

func NuevoClient() *Client {
	return &Client{Persona: &PersonaClient{datos: map[int]Persona{}}}
}

type PersonaClient struct {
	datos       map[int]Persona
	siguienteID int
}

func (c *PersonaClient) Create() *PersonaCreate {
	return &PersonaCreate{client: c}
}

func (c *PersonaClient) Query() *PersonaQuery {
	return &PersonaQuery{client: c}
}

// PersonaCreate es el análogo del builder *ent.PersonaCreate generado.
type PersonaCreate struct {
	client *PersonaClient
	nombre string
	edad   int
}

func (pc *PersonaCreate) SetNombre(n string) *PersonaCreate {
	pc.nombre = n
	return pc
}

func (pc *PersonaCreate) SetEdad(e int) *PersonaCreate {
	pc.edad = e
	return pc
}

func (pc *PersonaCreate) Save(ctx context.Context) (Persona, error) {
	pc.client.siguienteID++
	p := Persona{ID: pc.client.siguienteID, Nombre: pc.nombre, Edad: pc.edad}
	pc.client.datos[p.ID] = p
	return p, nil
}

// PersonaQuery es el análogo del builder *ent.PersonaQuery generado.
type PersonaQuery struct {
	client   *PersonaClient
	edadMin  int
	filtrar  bool
}

func (pq *PersonaQuery) WhereEdadMayorQue(edad int) *PersonaQuery {
	pq.edadMin = edad
	pq.filtrar = true
	return pq
}

func (pq *PersonaQuery) All(ctx context.Context) ([]Persona, error) {
	var resultado []Persona
	for _, p := range pq.client.datos {
		if pq.filtrar && p.Edad <= pq.edadMin {
			continue
		}
		resultado = append(resultado, p)
	}
	return resultado, nil
}

var ErrNoEncontrado = errors.New("persona no encontrada")

func (pq *PersonaQuery) OnlyID(ctx context.Context, id int) (Persona, error) {
	p, ok := pq.client.datos[id]
	if !ok {
		return Persona{}, ErrNoEncontrado
	}
	return p, nil
}

func main() {
	ctx := context.Background()
	client := NuevoClient()

	p1, _ := client.Persona.Create().SetNombre("Grace Hopper").SetEdad(85).Save(ctx)
	fmt.Printf("creada: %+v\n", p1)

	client.Persona.Create().SetNombre("Alan Turing").SetEdad(41).Save(ctx)
	client.Persona.Create().SetNombre("Margaret Hamilton").SetEdad(50).Save(ctx)

	mayores, _ := client.Persona.Query().WhereEdadMayorQue(45).All(ctx)
	fmt.Println("personas con edad > 45:")
	for _, p := range mayores {
		fmt.Printf("  %s (%d)\n", p.Nombre, p.Edad)
	}

	if _, err := client.Persona.Query().OnlyID(ctx, 999); errors.Is(err, ErrNoEncontrado) {
		fmt.Println("error esperado:", err)
	}
}
