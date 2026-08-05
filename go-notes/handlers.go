package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var startTime = time.Now()

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func main() {
	store := NewNoteStore()
	mux := http.NewServeMux()

	// GET /api/notes?category=work
	mux.HandleFunc("GET /api/notes", func(w http.ResponseWriter, r *http.Request) {
		cat := r.URL.Query().Get("category")
		notes := store.All()
		if cat != "" {
			filtered := notes[:0]
			for _, n := range notes {
				if strings.EqualFold(n.Category, cat) {
					filtered = append(filtered, n)
				}
			}
			notes = filtered
		}
		writeJSON(w, http.StatusOK, notes)
	})

	// GET /api/notes/{id}
	mux.HandleFunc("GET /api/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		n, ok := store.ByID(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, n)
	})

	// POST /api/notes
	mux.HandleFunc("POST /api/notes", func(w http.ResponseWriter, r *http.Request) {
		var req CreateNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		if req.Title == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title required"})
			return
		}
		if req.Category == "" {
			req.Category = "general"
		}
		writeJSON(w, http.StatusCreated, store.Add(req))
	})

	// DELETE /api/notes/{id}
	mux.HandleFunc("DELETE /api/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		if !store.Delete(id) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// GET /api/stats
	mux.HandleFunc("GET /api/stats", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(startTime).Round(time.Second).String()
		writeJSON(w, http.StatusOK, store.Stats(uptime))
	})

	// POST /api/notes/import — error handling example
	// Go has no try/catch. Every call that can fail returns (result, error).
	// You check it immediately — the happy path and error path live side by side.
	mux.HandleFunc("POST /api/notes/import", func(w http.ResponseWriter, r *http.Request) {
		var body ImportRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("invalid request body: %v", err),
			})
			return
		}

		data, err := os.ReadFile(body.FilePath)
		if err != nil {
			if os.IsNotExist(err) {
				writeJSON(w, http.StatusNotFound, map[string]string{
					"error": fmt.Sprintf("file not found: %s", body.FilePath),
				})
			} else {
				writeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": fmt.Sprintf("failed to read file: %v", err),
				})
			}
			return
		}

		var notes []CreateNoteRequest
		if err := json.Unmarshal(data, &notes); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("invalid JSON in file: %v", err),
			})
			return
		}

		imported := 0
		var errors []string
		for i, n := range notes {
			if n.Title == "" {
				errors = append(errors, fmt.Sprintf("note %d: title required", i))
				continue
			}
			if n.Category == "" {
				n.Category = "general"
			}
			store.Add(n)
			imported++
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"imported": imported,
			"errors":   errors,
		})
	})

	// Serve UI
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, indexHTML)
	})

	addr := ":5080"
	log.Printf("Go Notes running on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
