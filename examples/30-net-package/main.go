// Ejemplo del Capítulo 30: un servidor TCP mínimo y un cliente que se
// conecta a él, en el mismo proceso para que el ejemplo sea autocontenido.
package main

import (
	"bufio"
	"fmt"
	"net"
)

func iniciarServidorEco(listo chan<- string) {
	ln, err := net.Listen("tcp", "127.0.0.1:0") // puerto 0: el SO elige uno libre
	if err != nil {
		fmt.Println("error al escuchar:", err)
		return
	}
	defer ln.Close()

	listo <- ln.Addr().String()

	conn, err := ln.Accept()
	if err != nil {
		return
	}
	defer conn.Close()

	linea, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}
	fmt.Fprintf(conn, "eco: %s", linea)
}

func main() {
	direccion := make(chan string)
	go iniciarServidorEco(direccion)

	addr := <-direccion
	fmt.Println("servidor escuchando en", addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("error al conectar:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "hola desde el cliente\n")

	respuesta, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("error al leer respuesta:", err)
		return
	}
	fmt.Print("cliente recibió: ", respuesta)
}
