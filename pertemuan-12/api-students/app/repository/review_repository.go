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

type replyEntity struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    int           `bson:"user_id"`
	Comment   string        `bson:"comment"`
	CreatedAt time.Time     `bson:"created_at"`
}

type reviewEntity struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	CourseID  int           `bson:"course_id"`
	UserID    int           `bson:"user_id"`
	Rating    int           `bson:"rating"`
	Comment   string        `bson:"comment"`
	Tags      []string      `bson:"tags"`
	Replies   []replyEntity `bson:"replies"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func toReviewModel(e reviewEntity) model.Review {
	var replies []model.Reply
	if len(e.Replies) > 0 {
		replies = make([]model.Reply, len(e.Replies))
		for i, rep := range e.Replies {
			replies[i] = model.Reply{
				ID:        rep.ID.Hex(),
				UserID:    rep.UserID,
				Comment:   rep.Comment,
				CreatedAt: rep.CreatedAt,
			}
		}
	}

	return model.Review{
		ID:        e.ID.Hex(),
		CourseID:  e.CourseID,
		UserID:    e.UserID,
		Rating:    e.Rating,
		Comment:   e.Comment,
		Tags:      e.Tags,
		Replies:   replies,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func NewReviewRepository(db *mongo.Database) ReviewRepository {
	return &reviewRepository{
		coll: db.Collection("reviews"),
	}
}

func (r *reviewRepository) EnsureIndexes(ctx context.Context) error {
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

	var entities []reviewEntity
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}

	reviews := make([]model.Review, 0, len(entities))
	for _, e := range entities {
		reviews = append(reviews, toReviewModel(e))
	}

	return reviews, nil
}
