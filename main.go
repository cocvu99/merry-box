package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

// fileStore is an in-memory store for milestone M1.
// NOTE: data lives only in RAM and is lost on restart — that is expected
// at this milestone. Persistence (disk, then PostgreSQL) comes later.
type fileStore struct {
	mu    sync.RWMutex
	files map[string][]byte
}

func newFileStore() *fileStore {
	return &fileStore{
		files: make(map[string][]byte),
	}
}

// put stores content under name, overwriting any existing entry.
func (s *fileStore) put(name string, content []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.files[name] = content
}

// get returns the content for name and whether it was found.
func (s *fileStore) get(name string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	content, ok := s.files[name]
	return content, ok
}

// delete removes name from the store and reports whether it existed.
func (s *fileStore) delete(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, existed := s.files[name]
	delete(s.files, name)
	return existed
}

func main() {
	store := newFileStore()
	mux := http.NewServeMux()

	// --- Health check: fully implemented as a reference for the handlers below. ---
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// --- Upload: raw body + ?name=... (multipart parsing comes at a later milestone). ---
	mux.HandleFunc("POST /files", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing required query parameter: name", http.StatusBadRequest)
			return
		}
		// We knowingly load the whole body into RAM here; a later milestone
		// switches to streaming so large files do not blow up memory.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusInternalServerError)
			return
		}
		store.put(name, body)
		w.WriteHeader(http.StatusCreated)
	})

	// --- Retrieve by name. ---
	mux.HandleFunc("GET /files/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		content, ok := store.get(name)
		if !ok {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		_, _ = w.Write(content)
	})

	// --- Delete by name. ---
	// DELETE is NOT silently idempotent here: deleting a missing file returns
	// 404 so clients can tell "already gone" from "never existed". Documented
	// choice per the challenge ("idempotent or not? pick and document").
	mux.HandleFunc("DELETE /files/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if !store.delete(name) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
