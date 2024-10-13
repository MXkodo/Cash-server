package models

import "time"

type Document struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Mime      string    `json:"mime"`
	FilePath  string    `json:"file_path"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
}
