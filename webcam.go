// Command webcam serves a web stream from a V4L camera
package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	addr     = flag.String("addr", ":8003", "http address to listen on")
	dataRoot = flag.String("data", "./ui", "data dir")
)

func stream(w http.ResponseWriter, r *http.Request) {
	log.Println("Stream:", r.RemoteAddr, "connected")
	w.Header().Add("Content-Type", "video/x-vnd.on2.vp8")

	streamVid(w)
	log.Println("Stream:", r.RemoteAddr, "disconnected")
}

func main() {
	flag.Parse()
	http.HandleFunc("/cam", stream)
	http.Handle("/", http.FileServer(http.Dir(*dataRoot)))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
