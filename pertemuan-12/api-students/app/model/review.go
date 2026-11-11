package model

import "time"

type Reply struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    int       `json:"user_id" bson:"user_id"`
	Comment   string    `json:"comment" bson:"comment"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Review struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	CourseID  int       `json:"course_id" bson:"course_id"`
	UserID    int       `json:"user_id" bson:"user_id"`
	Rating    int       `json:"rating" bson:"rating"`
	Comment   string    `json:"comment" bson:"comment"`
	Tags      []string  `json:"tags" bson:"tags"`
	Replies   []Reply   `json:"replies" bson:"replies"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateReviewRequest struct {
	Rating  int      `json:"rating" validate:"required,min=1,max=5"`
	Comment string   `json:"comment" validate:"required,min=5,max=1000"`
	Tags    []string `json:"tags" validate:"omitempty,dive,max=30"`
}