// Command webcam serves a WebRTC stream from a camera.
package main

import (
	"flag"
	"log"
	"net/http"

	"0f.io/webcam"
)

var (
	addr     = flag.String("addr", ":8003", "http address to listen on")
	dataRoot = flag.String("data", "./ui", "data dir")
)

func main() {
	flag.Parse()
	http.HandleFunc("/cam/", webcam.ServeHTTP)
	http.Handle("/", http.FileServer(http.Dir(*dataRoot)))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
