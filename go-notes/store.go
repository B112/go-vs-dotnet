package main

import (
	"sync"
	"time"
)

type NoteStore struct {
	mu     sync.RWMutex
	notes  []Note
	nextID int
}

func NewNoteStore() *NoteStore {
	s := &NoteStore{nextID: 1}
	for _, n := range []CreateNoteRequest{
		{"Setup Balena fleet", "Configure fleet variables for CM5 devices", "work"},
		{"Buy groceries", "Milk, bread, cheese, coffee", "personal"},
		{"WireGuard VPN config", "Generate new peer keys for test devices", "work"},
	} {
		s.Add(n)
	}
	return s
}

func (s *NoteStore) All() []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Note, len(s.notes))
	copy(result, s.notes)
	return result
}

func (s *NoteStore) ByID(id int) (Note, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, n := range s.notes {
		if n.ID == id {
			return n, true
		}
	}
	return Note{}, false
}

func (s *NoteStore) Add(req CreateNoteRequest) Note {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := Note{
		ID:        s.nextID,
		Title:     req.Title,
		Content:   req.Content,
		Category:  req.Category,
		CreatedAt: time.Now(),
	}
	s.nextID++
	s.notes = append(s.notes, n)
	return n
}

func (s *NoteStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, n := range s.notes {
		if n.ID == id {
			s.notes = append(s.notes[:i], s.notes[i+1:]...)
			return true
		}
	}
	return false
}

func (s *NoteStore) Stats(uptime string) Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cats := make(map[string]int)
	for _, n := range s.notes {
		cats[n.Category]++
	}
	return Stats{
		Total:      len(s.notes),
		ByCategory: cats,
		Uptime:     uptime,
	}
}
