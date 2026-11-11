package model

import "time"

// Project disimpan secara embedding di dalam Portfolio
type Project struct {
	ID           string   `json:"id"`
	Title        string   `json:"title" validate:"required,min=3,max=100"`
	Technologies []string `json:"technologies" validate:"required,min=1,dive,min=1,max=30"`
	Link         string   `json:"link" validate:"omitempty,url"`
	Year         int      `json:"year" validate:"required,min=2000,max=2100"`
}

// Certificate disimpan secara embedding di dalam Portfolio
type Certificate struct {
	Name   string    `json:"name" validate:"required,min=3,max=100"`
	Issuer string    `json:"issuer" validate:"required,min=2,max=100"`
	Date   time.Time `json:"date" validate:"required"`
}

// Portfolio adalah entitas utama dokumen portofolio mahasiswa
type Portfolio struct {
	ID           string        `json:"id"`
	StudentID    int           `json:"student_id"`
	Headline     string        `json:"headline"`
	Bio          string        `json:"bio"`
	Projects     []Project     `json:"projects"`
	Certificates []Certificate `json:"certificates"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// UpsertPortfolioRequest digunakan untuk PUT /students/:id/portfolio
type UpsertPortfolioRequest struct {
	Headline     string        `json:"headline" validate:"required,max=150"`
	Bio          string        `json:"bio" validate:"required,max=1000"`
	Certificates []Certificate `json:"certificates" validate:"omitempty,dive"`
}

// AddProjectRequest digunakan untuk POST /students/:id/portfolio/projects
type AddProjectRequest struct {
	Title        string   `json:"title" validate:"required,min=3,max=100"`
	Technologies []string `json:"technologies" validate:"required,min=1,max=10,dive,min=1,max=30"`
	Link         string   `json:"link" validate:"omitempty,url"`
	Year         int      `json:"year" validate:"required,min=2000,max=2100"`
}

// TechnologyStat untuk aggregation top 10 teknologi
type TechnologyStat struct {
	Technology string `json:"technology"`
	Total      int    `json:"total"`
}