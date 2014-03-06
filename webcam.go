// Command webcam serves a web stream from a camera.
package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
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

func readMsgs(r io.Reader, pc PeerConn) {
	dec := json.NewDecoder(r)
	var msg map[string]interface{}
	for {
		err := dec.Decode(&msg)
		if err != nil {
			return
		}

		switch msg["type"] {
		case "candidate":
			log.Println("got candidate")
			pc.AddCandidate(
				msg["candidate"].(string),
				msg["sdpMid"].(string),
				int(msg["sdpMLineIndex"].(float64)),
			)
		case "answer":
			log.Println("got answer")
			pc.AddAnswer(msg["sdp"].(string))
		default:
			log.Println("got unknow json message:", msg)
		}
	}
}

func call(ws *websocket.Conn) {
	log.Println(ws.Request().RemoteAddr, "connected")
	defer log.Println(ws.Request().RemoteAddr, "disconnected")

	enc := json.NewEncoder(ws)
	pc := Offer()
	go readMsgs(ws, pc)

	for {
		select {
		case c := <-pc.Candidate:
			log.Println("sending candidate")
			c.Type = "candidate"
			enc.Encode(c)
		case sdp := <-pc.Offer:
			log.Println("sending offer")
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
