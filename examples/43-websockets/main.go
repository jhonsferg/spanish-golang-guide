// Ejemplo del Capítulo 43: un servidor y un cliente WebSocket mínimos
// con github.com/coder/websocket - comunicación bidireccional sobre una
// sola conexión persistente, a diferencia del modelo request/response de
// HTTP normal.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

func servidorEco() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer conn.CloseNow()

		ctx := r.Context()
		for {
			tipo, msg, err := conn.Read(ctx)
			if err != nil {
				return // el cliente cerró la conexión
			}
			eco := append([]byte("eco: "), msg...)
			if err := conn.Write(ctx, tipo, eco); err != nil {
				return
			}
		}
	})
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Handler: servidorEco()}
	go server.Serve(ln)
	defer server.Close()

	wsURL := "ws://" + ln.Addr().String()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.CloseNow()

	mensajes := []string{"hola", "mundo", "desde websockets"}
	for _, m := range mensajes {
		if err := conn.Write(ctx, websocket.MessageText, []byte(m)); err != nil {
			log.Fatal(err)
		}
		_, respuesta, err := conn.Read(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("enviado=%q recibido=%q\n", m, string(respuesta))
	}

	conn.Close(websocket.StatusNormalClosure, "listo")
}
