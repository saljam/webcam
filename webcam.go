// Command webcam serves a web stream from a V4L camera
package main

import (
	"flag"
	"log"
	"net/http"
	"io"
	"sync"
	
	"code.google.com/p/go.net/websocket"
)

var (
	addr     = flag.String("addr", ":8003", "http address to listen on")
	dataRoot = flag.String("data", "./ui", "data dir")

	c, a *websocket.Conn
	wg, done sync.WaitGroup
)

func call(ws *websocket.Conn) {
	log.Println(ws.Request().RemoteAddr, "connected")
	defer log.Println(ws.Request().RemoteAddr, "disconnected")

	c = ws
	wg.Done()
	wg.Wait()
	done.Add(1)
	log.Println("copying")
	io.Copy(a, c)
	wg.Add(1)
	done.Done()
	done.Wait()
}

func answer(ws *websocket.Conn) {
	log.Println(ws.Request().RemoteAddr, "connected")
	defer log.Println(ws.Request().RemoteAddr, "disconnected")

	a = ws
	wg.Done()
	wg.Wait()
	done.Add(1)
	log.Println("copying")
	io.Copy(c, a)
	done.Done()
	done.Wait()
}

func main() {
	flag.Parse()
	wg.Add(2)
	http.Handle("/call", websocket.Handler(call))
	http.Handle("/answer", websocket.Handler(answer))
	http.Handle("/", http.FileServer(http.Dir(*dataRoot)))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
