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
	// TODO:
	// A Go map is NOT safe for concurrent writes — two goroutines writing at
	// once will panic. Protect the write with the mutex.
	// Hint: s.mu.Lock() ... defer s.mu.Unlock(), then assign into s.files.
}

// get returns the content for name and whether it was found.
func (s *fileStore) get(name string) ([]byte, bool) {
	// TODO:
	// A read can use the read-lock (allows many concurrent readers).
	// Hint: s.mu.RLock() ... defer s.mu.RUnlock().
	// Hint: the comma-ok form `v, ok := s.files[name]` tells you if it exists.
	return nil, false
}

// delete removes name from the store and reports whether it existed.
func (s *fileStore) delete(name string) bool {
	// TODO:
	// This is a write, so use the full Lock (not RLock).
	// Check existence first (comma-ok), then call the builtin delete(s.files, name).
	// Return whether the key existed BEFORE deletion — the handler needs this
	// to decide 204 vs 404.
	return false
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
		// TODO:
		// 1. Read the file name from the query string: r.URL.Query().Get("name").
		//    If it is empty, respond http.StatusBadRequest (400) and return.
		// 2. Read the whole body: body, err := io.ReadAll(r.Body).
		//    (We knowingly load it all into RAM here; a later milestone switches
		//    to streaming so large files do not blow up memory.)
		//    On error, respond http.StatusInternalServerError (500) and return.
		// 3. store.put(name, body), then respond http.StatusCreated (201).
		_ = io.ReadAll // remove this line once you use io.ReadAll
	})

	// --- Retrieve by name. ---
	mux.HandleFunc("GET /files/{name}", func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		// 1. name := r.PathValue("name")
		// 2. content, ok := store.get(name)
		// 3. if !ok -> respond http.StatusNotFound (404) and return.
		// 4. otherwise w.Write(content) (status defaults to 200).
	})

	// --- Delete by name. ---
	mux.HandleFunc("DELETE /files/{name}", func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		// 1. name := r.PathValue("name")
		// 2. existed := store.delete(name)
		// 3. Decide the semantics you documented:
		//    - if !existed -> http.StatusNotFound (404)
		//    - otherwise   -> http.StatusNoContent (204)
		//    (This is the idempotency choice we discussed — own it and document it.)
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
