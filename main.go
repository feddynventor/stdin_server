package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const bufferSize = 5 * 1024 * 1024 // 5MB

type Server struct {
	buffer      *bytes.Buffer
	bufferLock  sync.Mutex
	dataUpdated *sync.Cond
}

func NewServer() *Server {
	server := &Server{
		buffer: &bytes.Buffer{},
	}
	server.dataUpdated = sync.NewCond(&server.bufferLock)
	return server
}

func (s *Server) writeToBuffer(data []byte) {
	s.bufferLock.Lock()
	defer s.bufferLock.Unlock()

	// Add data to the buffer
	s.buffer.Write(data)

	// Truncate buffer if it exceeds bufferSize
	if s.buffer.Len() > bufferSize {
		overflow := s.buffer.Len() - bufferSize
		_, _ = s.buffer.Read(make([]byte, overflow))
	}

	// Notify readers
	s.dataUpdated.Broadcast()
}

func (s *Server) streamHandler(w http.ResponseWriter, r *http.Request) {
	startParam := r.URL.Query().Get("start")
	start := 0
	if startParam != "" {
		var err error
		start, err = strconv.Atoi(startParam)
		if err != nil || start < 0 {
			http.Error(w, "Invalid start parameter", http.StatusBadRequest)
			return
		}
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")

	// Stream data
	s.bufferLock.Lock()
	defer s.bufferLock.Unlock()
	for {
		// Check if there's enough data to start from the requested offset
		if start < s.buffer.Len() {
			_, err := w.Write(s.buffer.Bytes()[start:])
			if err != nil {
				log.Println("Error writing to response:", err)
				return
			}
			start = s.buffer.Len()
			flusher.Flush()
		}

		// Wait for new data
		s.dataUpdated.Wait()
	}
}

func main() {
	server := NewServer()

	// Start reading from stdin
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("Error reading stdin:", err)
				}
				break
			}
			server.writeToBuffer(buf[:n])
		}
	}()

	http.HandleFunc("/stream", server.streamHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
