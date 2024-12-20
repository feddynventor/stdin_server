package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
)

func main() {
	// Create a server on port 8000
	// Exactly how you would run an HTTP/1.1 server
	srv := &http.Server{Addr: ":8000", Handler: http.HandlerFunc(handle)}

	// Start the server with TLS, since we are running HTTP/2 it must be
	// run with TLS.
	// Exactly how you would run an HTTP/1.1 server with TLS connection.
	log.Printf("Serving %s on https://0.0.0.0:8000", os.Args[1])
	log.Fatal(srv.ListenAndServeTLS("tls/server.crt", "tls/server.key"))
}

func handle(w http.ResponseWriter, r *http.Request) {
	// Log the request protocol
	log.Printf("Got connection: %s", r.Proto)
	// Headers
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Transfer-Encoding", "chunked")
	// Open file
	dat, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewScanner(dat)
	reader.Split(bufio.ScanBytes)
	// Buffered sender
	for reader.Scan() {
		w.Write([]byte(reader.Bytes()))
	}
}
