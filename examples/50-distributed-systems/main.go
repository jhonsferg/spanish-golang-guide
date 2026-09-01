// Ejemplo del Capítulo 50: consistent hashing — la técnica que usan
// sistemas distribuidos (caches, bases de datos particionadas) para
// repartir claves entre nodos sin remapear casi todo cuando un nodo
// entra o sale.
package main

import (
	"fmt"
	"hash/fnv"
	"sort"
)

type AnilloConsistente struct {
	nodosVirtuales int
	anillo         map[uint32]string
	posiciones     []uint32
}

func NuevoAnillo(nodosVirtuales int) *AnilloConsistente {
	return &AnilloConsistente{
		nodosVirtuales: nodosVirtuales,
		anillo:         make(map[uint32]string),
	}
}

func hashear(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (a *AnilloConsistente) AgregarNodo(nodo string) {
	for i := 0; i < a.nodosVirtuales; i++ {
		clave := hashear(fmt.Sprintf("%s#%d", nodo, i))
		a.anillo[clave] = nodo
		a.posiciones = append(a.posiciones, clave)
	}
	sort.Slice(a.posiciones, func(i, j int) bool { return a.posiciones[i] < a.posiciones[j] })
}

func (a *AnilloConsistente) QuitarNodo(nodo string) {
	nuevasPosiciones := a.posiciones[:0]
	for _, p := range a.posiciones {
		if a.anillo[p] == nodo {
			delete(a.anillo, p)
			continue
		}
		nuevasPosiciones = append(nuevasPosiciones, p)
	}
	a.posiciones = nuevasPosiciones
}

// NodoPara retorna el nodo responsable de una clave: el primer punto del
// anillo en sentido horario a partir del hash de la clave.
func (a *AnilloConsistente) NodoPara(clave string) string {
	if len(a.posiciones) == 0 {
		return ""
	}
	h := hashear(clave)
	idx := sort.Search(len(a.posiciones), func(i int) bool { return a.posiciones[i] >= h })
	if idx == len(a.posiciones) {
		idx = 0 // el anillo da la vuelta
	}
	return a.anillo[a.posiciones[idx]]
}

func main() {
	anillo := NuevoAnillo(100) // 100 nodos virtuales por nodo real: mejor distribución
	anillo.AgregarNodo("cache-1")
	anillo.AgregarNodo("cache-2")
	anillo.AgregarNodo("cache-3")

	claves := []string{"usuario:42", "usuario:7", "sesión:abc", "carrito:99"}
	fmt.Println("--- antes de escalar ---")
	asignacionOriginal := map[string]string{}
	for _, c := range claves {
		nodo := anillo.NodoPara(c)
		asignacionOriginal[c] = nodo
		fmt.Printf("%-14s -> %s\n", c, nodo)
	}

	// Quitar un nodo: solo las claves que apuntaban a ÉL se reasignan,
	// no todas — esa es la ventaja frente a un hash % N ingenuo.
	anillo.QuitarNodo("cache-2")

	fmt.Println("--- después de quitar cache-2 ---")
	movidas := 0
	for _, c := range claves {
		nodo := anillo.NodoPara(c)
		if nodo != asignacionOriginal[c] {
			movidas++
		}
		fmt.Printf("%-14s -> %s\n", c, nodo)
	}
	fmt.Printf("claves reasignadas: %d de %d\n", movidas, len(claves))
}
