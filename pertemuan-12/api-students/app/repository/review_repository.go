package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"api-students/app/model"
)

var (
	ErrNotOwner     = errors.New("bukan pemilik resource")
	ErrOwnReview    = errors.New("tidak dapat membalas atau menilai ulasan sendiri")
	ErrAlreadyVoted = errors.New("sudah memberikan vote")
	ErrReplyLimit   = errors.New("batas balasan tercapai")
)

type ReviewRepository interface {
	EnsureIndexes(ctx context.Context) error
	Create(ctx context.Context, r model.Review) (model.Review, error)
	FindByCourse(ctx context.Context, courseID int) ([]model.Review, error)
}

type reviewRepository struct {
	coll *mongo.Collection
}

func NewReviewRepository(db *mongo.Database) ReviewRepository {
	return &reviewRepository{
		coll: db.Collection("reviews"),
	}
}

func (r *reviewRepository) EnsureIndexes(ctx context.Context) error {
	// Unique compound index: satu user hanya boleh memberi 1 review per course
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "course_id", Value: 1},
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := r.coll.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (r *reviewRepository) Create(ctx context.Context, rev model.Review) (model.Review, error) {
	now := time.Now()
	rev.CreatedAt = now
	rev.UpdatedAt = now

	res, err := r.coll.InsertOne(ctx, rev)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return model.Review{}, ErrDuplicate
		}
		return model.Review{}, err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		rev.ID = oid.Hex()
	}

	return rev, nil
}

func (r *reviewRepository) FindByCourse(ctx context.Context, courseID int) ([]model.Review, error) {
	filter := bson.M{"course_id": courseID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []model.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	if reviews == nil {
		reviews = []model.Review{}
	}

	return reviews, nil
}