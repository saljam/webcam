// Command webcam serves a web stream from a camera.
package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"

	"code.google.com/p/go.net/websocket"
)

var (
	addr     = flag.String("addr", ":8003", "http address to listen on")
	dataRoot = flag.String("data", "./ui", "data dir")

	c, a     *websocket.Conn
	wg, done sync.WaitGroup
)

type candidateMsg struct {
	Sdp  string `json:"candidate"`
	Mid  string `json:"sdpMid"`
	Line int    `json:"sdpMLineIndex"`
	Type string `json:"type"`
}

type offerMsg struct {
	Sdp  string `json:"sdp"`
	Type string `json:"type"`
}

var (
	candidate = make(chan candidateMsg)
	offer     = make(chan string)
)

func readMsgs(r io.Reader) {
	runtime.LockOSThread()
	
	dec := json.NewDecoder(r)
	var msg map[string]interface{}
	MakePeerConnection()
	for {
		err := dec.Decode(&msg)
		if err != nil {
			return
		}
		
		switch msg["type"] {
			case "candidate":
				log.Println("got candidate")
				Candidate(
					msg["candidate"].(string),
					msg["sdpMid"].(string),
					int(msg["sdpMLineIndex"].(float64)),
				)
				log.Println("done candidate")
			case "answer":
				log.Println("got answer")
				Answer(msg["sdp"].(string))
				log.Println("done answer")
			default:
				log.Println("got unknow json message:", msg)
		}
	}
}

func call(ws *websocket.Conn) {
	log.Println(ws.Request().RemoteAddr, "connected")
	defer log.Println(ws.Request().RemoteAddr, "disconnected")

	enc := json.NewEncoder(ws)
	go readMsgs(ws)
	runtime.LockOSThread()

	for {
		select {
		case c := <-candidate:
			log.Println("local candidate")
			c.Type = "candidate"
			enc.Encode(c)
		case sdp := <-offer:
			log.Println("local offer")
			enc.Encode(offerMsg{
				Sdp:  sdp,
				Type: "offer",
			})
		}
	}
}

func main() {
	flag.Parse()
	http.Handle("/call", websocket.Handler(call))
	http.Handle("/", http.FileServer(http.Dir(*dataRoot)))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
