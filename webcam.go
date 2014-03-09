// Package webcam provides an http.Handler to serve a WebRTC stream from a camera.
package webcam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"code.google.com/p/go.net/websocket"
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

type Webcam struct {
	wsHandler http.Handler
}

func (s *Webcam) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "webrtc-stream.html") {
		p := r.URL.Path[:len(r.URL.Path)-len("webrtc-stream.html")]
		fmt.Fprintf(w, template, p)
		return
	}
	s.wsHandler.ServeHTTP(w, r)
}

func NewWebcam() *Webcam {
	return &Webcam{websocket.Handler(call)}
}
