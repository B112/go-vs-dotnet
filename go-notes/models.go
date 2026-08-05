package main

import "time"

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateNoteRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

type ImportRequest struct {
	FilePath string `json:"filePath"`
}

type Stats struct {
	Total      int            `json:"total"`
	ByCategory map[string]int `json:"byCategory"`
	Uptime     string         `json:"uptime"`
}
