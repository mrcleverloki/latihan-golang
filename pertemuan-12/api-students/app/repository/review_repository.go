package repository

import (
	"context"
	"errors"

	"api-students/app/model"
)

var (
	ErrNotOwner     = errors.New("bukan pemilik resource")
	ErrOwnReview    = errors.New("tidak dapat membalas atau menilai ulasan sendiri")
	ErrAlreadyVoted = errors.New("sudah memberikan vote")
	ErrReplyLimit   = errors.New("batas balasan tercapai")
)

type ReviewRepository interface {
	Create(ctx context.Context, r model.Review) (model.Review, error)
	FindByCourse(ctx context.Context, courseID int) ([]model.Review, error)
}