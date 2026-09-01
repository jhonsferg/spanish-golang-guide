// Ejemplo del Capítulo 52: una caché LRU (Least Recently Used) desde
// cero con container/list - el mismo principio detrás de cachés en
// producción, sin depender de una librería externa.
package main

import (
	"container/list"
	"fmt"
)

type entrada struct {
	clave string
	valor any
}

type CacheLRU struct {
	capacidad int
	orden     *list.List // el frente es "más reciente"
	elementos map[string]*list.Element
}

func NuevaCacheLRU(capacidad int) *CacheLRU {
	return &CacheLRU{
		capacidad: capacidad,
		orden:     list.New(),
		elementos: make(map[string]*list.Element),
	}
}

func (c *CacheLRU) Get(clave string) (any, bool) {
	el, ok := c.elementos[clave]
	if !ok {
		return nil, false
	}
	c.orden.MoveToFront(el)
	return el.Value.(*entrada).valor, true
}

func (c *CacheLRU) Put(clave string, valor any) {
	if el, ok := c.elementos[clave]; ok {
		el.Value.(*entrada).valor = valor
		c.orden.MoveToFront(el)
		return
	}

	if c.orden.Len() >= c.capacidad {
		masAntiguo := c.orden.Back()
		if masAntiguo != nil {
			c.orden.Remove(masAntiguo)
			delete(c.elementos, masAntiguo.Value.(*entrada).clave)
		}
	}

	el := c.orden.PushFront(&entrada{clave: clave, valor: valor})
	c.elementos[clave] = el
}

func main() {
	cache := NuevaCacheLRU(3)

	cache.Put("a", 1)
	cache.Put("b", 2)
	cache.Put("c", 3)

	cache.Get("a") // "a" pasa a ser la más recientemente usada

	cache.Put("d", 4) // la caché está llena: expulsa la menos usada ("b")

	for _, clave := range []string{"a", "b", "c", "d"} {
		valor, ok := cache.Get(clave)
		fmt.Printf("%s -> valor=%v presente=%v\n", clave, valor, ok)
	}
}
