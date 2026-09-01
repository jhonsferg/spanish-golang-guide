// Mini-proyecto de la Parte VIII (Frameworks web): la MISMA API de
// TODOs (crear, listar, completar, borrar) implementada cuatro veces —
// una por framework, en ginapi/, echoapi/, fiberapi/ y chiapi/ — para
// comparar directamente la sintaxis de cada uno resolviendo idéntico
// problema. Este main.go ejecuta la misma secuencia de requests contra
// las cuatro implementaciones.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jhonsferg/spanish-golang-guide/examples/miniproyectos/todo-api-comparativa/chiapi"
	"github.com/jhonsferg/spanish-golang-guide/examples/miniproyectos/todo-api-comparativa/echoapi"
	"github.com/jhonsferg/spanish-golang-guide/examples/miniproyectos/todo-api-comparativa/fiberapi"
	"github.com/jhonsferg/spanish-golang-guide/examples/miniproyectos/todo-api-comparativa/ginapi"
)

// ejercitar corre la misma secuencia CRUD contra cualquier http.Handler.
func ejercitar(nombre string, handler http.Handler) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println(nombre, "error:", err)
		return
	}
	server := &http.Server{Handler: handler}
	go server.Serve(ln)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	base := "http://" + ln.Addr().String()
	cliente := &http.Client{}

	respCrear, _ := cliente.Post(base+"/todos", "application/json", strings.NewReader(`{"titulo":"aprender frameworks"}`))
	var creado map[string]any
	json.NewDecoder(respCrear.Body).Decode(&creado)
	respCrear.Body.Close()

	respRechazo, _ := cliente.Post(base+"/todos", "application/json", strings.NewReader(`{"titulo":""}`))
	statusRechazo := respRechazo.StatusCode
	respRechazo.Body.Close()

	id := int(creado["id"].(float64))
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/todos/%d", base, id), nil)
	respPut, _ := cliente.Do(req)
	statusPut := respPut.StatusCode
	respPut.Body.Close()

	respLista, _ := cliente.Get(base + "/todos")
	var lista []map[string]any
	json.NewDecoder(respLista.Body).Decode(&lista)
	respLista.Body.Close()

	reqDel, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/todos/%d", base, id), nil)
	respDel, _ := cliente.Do(reqDel)
	statusDel := respDel.StatusCode
	respDel.Body.Close()

	fmt.Printf("%-6s crear=%d rechazo-validación=%d completar=%d listar=%d(%d items) borrar=%d\n",
		nombre, respCrear.StatusCode, statusRechazo, statusPut, respLista.StatusCode, len(lista), statusDel)
}

func main() {
	ejercitar("Gin", ginapi.NewRouter())
	ejercitar("Echo", echoapi.NewRouter())
	ejercitar("Chi", chiapi.NewRouter())

	// Fiber no implementa http.Handler (corre sobre fasthttp), así que
	// se prueba con su propio helper .Test() en vez de net.Listener.
	app := fiberapi.NewApp()
	respCrear, _ := app.Test(mustRequest(http.MethodPost, "/todos", `{"titulo":"aprender fiber"}`))
	var creado map[string]any
	json.NewDecoder(respCrear.Body).Decode(&creado)

	respLista, _ := app.Test(mustRequest(http.MethodGet, "/todos", ""))
	var lista []map[string]any
	json.NewDecoder(respLista.Body).Decode(&lista)

	fmt.Printf("%-6s crear=%d listar=%d(%d items)\n", "Fiber", respCrear.StatusCode, respLista.StatusCode, len(lista))
}

func mustRequest(metodo, ruta, cuerpo string) *http.Request {
	var req *http.Request
	if cuerpo != "" {
		req, _ = http.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(metodo, ruta, nil)
	}
	return req
}
