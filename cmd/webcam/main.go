// Command webcam serves a WebRTC stream from a camera.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"0f.io/webcam"
)

var (
	addr     = flag.String("addr", ":8003", "http address to listen on")
	dataRoot = flag.String("data", "./ui", "data dir")
)

func handleSession(w http.ResponseWriter, r *http.Request) {
	offer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Couldn't read offer body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	s := webcam.NewSession()
	err = s.Remote(offer)
	if err != nil {
		log.Printf("Couldn't set remote description: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	answer, err := s.Description()
	if err != nil {
		log.Printf("Couldn't get local description: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "%s", answer)
}

func main() {
	flag.Parse()
	http.HandleFunc("/session", handleSession)
	http.Handle("/", http.FileServer(http.Dir(*dataRoot)))
	log.Fatal(http.ListenAndServe(*addr, nil))
}
