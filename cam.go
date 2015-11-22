// Package webcam provides an http.Handler to serve a WebRTC stream from a camera,
// using the native API of the reference WebRTC implementation.
package webcam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
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

// ServeHTTP handles WebRTC requests over websockets.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "webrtc-stream.html") {
		p := r.URL.Path[:len(r.URL.Path)-len("webrtc-stream.html")]
		fmt.Fprintf(w, template, p)
		return
	}
	websocket.Handler(call).ServeHTTP(w, r)
}

var template = `<polymer-element name="webrtc-stream">
<template>
<video id=cam autoplay>
</template>
<style>
video {
	width: 100%%;
}
</style>
<script>
function error(err) {
	console.log("err:", err)
}

Polymer('webrtc-stream', {
	ready: function() {
		delete URL; // Polymer's platform.js overwrites URL with one that doesn't have createObjectURL().
		
		var ws = new WebSocket('ws://' + location.host + %s)
		var cfg = {"iceServers": [{"url": "stun:stun.l.google.com:19302"}]};
		pc = new RTCPeerConnection(cfg, {optional: [{RtpDataChannels: true}]});
		
		ws.onmessage = function(m) {
			var msg = JSON.parse(m.data);
			if (msg.type === 'offer') {
				pc.setRemoteDescription(new RTCSessionDescription(msg), function(){}, error);
				pc.createAnswer(function(desc) {
					pc.setLocalDescription(desc);
					ws.send(JSON.stringify(desc));
				}, error);
			} else if (msg.type === 'candidate') {
				pc.addIceCandidate(new RTCIceCandidate(msg));
			} else {
				console.log("what's this?", msg)
			}
		}
		
		pc.onicecandidate = function(e) {
			if (e.candidate) {
				e.candidate.type="candidate";
				console.log(e.candidate);
				ws.send(JSON.stringify(e.candidate));
			}
		}
		pc.onaddstream = function (e) {
			var vid = document.getElementById("cam");
			attachMediaStream(vid, e.stream);
		};
	}
});
</script>
</polymer-element>


<title>webcam</title>
<body>
`