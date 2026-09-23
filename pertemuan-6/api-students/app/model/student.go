package model

import "time"

// Student merepresentasikan entitas mahasiswa di database
type Student struct {
	ID        int       `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     string    `json:"grade"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Request struct untuk validasi payload input
type CreateStudentRequest struct {
	NIM      string `json:"nim"`
	Name     string `json:"name"`
	Grade    string `json:"grade"`
	IsActive *bool  `json:"is_active"`
}

type ReplaceStudentRequest struct {
	NIM      string `json:"nim"`
	Name     string `json:"name"`
	Grade    string `json:"grade"`
	IsActive bool   `json:"is_active"`
}

type PatchStudentRequest struct {
	NIM      *string `json:"nim"`
	Name     *string `json:"name"`
	Grade    *string `json:"grade"`
	IsActive *bool   `json:"is_active"`
}

// Struktur query parameter untuk filtering, sorting, pagination
type ListQuery struct {
	Page     int
	Limit    int
	Sort     string
	Order    string
	Search   string
	IsActive *bool
}

// Offset menghitung baris yang dilewati untuk klausa SQL OFFSET
func (q ListQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}

// Format respons standar API
type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}