package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	mu     sync.Mutex
	tasks  []Task
	nextID int
}

func NewStore() *Store { return &Store{tasks: []Task{}, nextID: 1} }

func (s *Store) list() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]Task, len(s.tasks))
	copy(result, s.tasks)
	return result
}

func (s *Store) create(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	task := Task{ID: s.nextID, Title: title, CreatedAt: time.Now().UTC()}
	s.nextID++
	s.tasks = append(s.tasks, task)
	return task
}

func (s *Store) delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return true
		}
	}
	return false
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func newHandler(store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, store.list())
	})
	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		var input struct { Title string `json:"title"` }
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		input.Title = strings.TrimSpace(input.Title)
		if input.Title == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
			return
		}
		writeJSON(w, http.StatusCreated, store.create(input.Title))
	})
	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
			return
		}
		if !store.delete(id) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}

func main() {
	log.Println("Task API running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", newHandler(NewStore())))
}
