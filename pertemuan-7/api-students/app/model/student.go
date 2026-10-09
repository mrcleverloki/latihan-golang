package model

import "time"

// Student merepresentasikan entitas mahasiswa di database
type Student struct {
	ID        int       `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     string    `json:"grade"`
	IsActive  bool      `json:"is_active"`
	OwnerID   *int      `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Request struct dengan validasi deklaratif (Tugas D.2)
type CreateStudentRequest struct {
	NIM      string `json:"nim" validate:"required,nim_format"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Grade    string `json:"grade" validate:"required"`
	IsActive *bool  `json:"is_active"`
}

type ReplaceStudentRequest struct {
	NIM      string `json:"nim" validate:"required,nim_format"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Grade    string `json:"grade" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type PatchStudentRequest struct {
	NIM      *string `json:"nim,omitempty" validate:"omitnil,nim_format"`
	Name     *string `json:"name,omitempty" validate:"omitnil,min=3,max=100"`
	Grade    *string `json:"grade,omitempty" validate:"omitnil"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// Model untuk Cursor Keyset Pagination (Tugas D.3)
type Cursor struct {
	CreatedAt time.Time
	ID        int
}

type CursorQuery struct {
	Limit    int
	After    *Cursor
	Search   string
	IsActive *bool
}

type CursorMeta struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

// Struktur query parameter lama (Offset pagination)
type ListQuery struct {
	Page     int
	Limit    int
	Sort     string
	Order    string
	Search   string
	IsActive *bool
}

func (q ListQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}

type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
