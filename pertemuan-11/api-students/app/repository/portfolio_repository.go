package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"api-students/app/model"
)

var (
	ErrInvalidID         = errors.New("format id tidak valid")
	ErrPortfolioNotFound = errors.New("portofolio tidak ditemukan")
	ErrProjectNotFound   = errors.New("proyek tidak ditemukan")
	ErrProjectLimit      = errors.New("jumlah proyek sudah mencapai batas maksimal 20")
)

const MaxProjectsPerPortfolio = 20

type PortfolioRepository interface {
	EnsureIndexes(ctx context.Context) error
	FindByStudentID(ctx context.Context, studentID int) (model.Portfolio, error)
	Upsert(ctx context.Context, studentID int, req model.UpsertPortfolioRequest) (model.Portfolio, error)
	AddProject(ctx context.Context, studentID int, p model.Project) (model.Portfolio, error)
	DeleteProject(ctx context.Context, studentID int, projectID string) (model.Portfolio, error)
	TopTechnologies(ctx context.Context) ([]model.TechnologyStat, error)
}

type portfolioDocument struct {
	ID           bson.ObjectID         `bson:"_id,omitempty"`
	StudentID    int                   `bson:"student_id"`
	Headline     string                `bson:"headline"`
	Bio          string                `bson:"bio"`
	Projects     []projectDocument     `bson:"projects"`
	Certificates []certificateDocument `bson:"certificates"`
	CreatedAt    time.Time             `bson:"created_at"`
	UpdatedAt    time.Time             `bson:"updated_at"`
}

type projectDocument struct {
	ID           bson.ObjectID `bson:"_id"`
	Title        string        `bson:"title"`
	Technologies []string      `bson:"technologies"`
	Link         string        `bson:"link"`
	Year         int           `bson:"year"`
}

type certificateDocument struct {
	Name   string    `bson:"name"`
	Issuer string    `bson:"issuer"`
	Date   time.Time `bson:"date"`
}

func (d portfolioDocument) toModel() model.Portfolio {
	projects := make([]model.Project, 0, len(d.Projects))
	for _, p := range d.Projects {
		techs := p.Technologies
		if techs == nil {
			techs = []string{}
		}
		projects = append(projects, model.Project{
			ID:           p.ID.Hex(),
			Title:        p.Title,
			Technologies: techs,
			Link:         p.Link,
			Year:         p.Year,
		})
	}

	certs := make([]model.Certificate, 0, len(d.Certificates))
	for _, c := range d.Certificates {
		certs = append(certs, model.Certificate{
			Name:   c.Name,
			Issuer: c.Issuer,
			Date:   c.Date,
		})
	}

	return model.Portfolio{
		ID:           d.ID.Hex(),
		StudentID:    d.StudentID,
		Headline:     d.Headline,
		Bio:          d.Bio,
		Projects:     projects,
		Certificates: certs,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

type portfolioMongoRepository struct {
	coll *mongo.Collection
}

func NewPortfolioRepository(db *mongo.Database) PortfolioRepository {
	return &portfolioMongoRepository{coll: db.Collection("student_portfolios")}
}

func nowMillis() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}

func parseMongoID(id string) (bson.ObjectID, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, ErrInvalidID
	}
	return oid, nil
}

// EnsureIndexes menjamin aturan satu mahasiswa satu portofolio dan mempercepat filter
func (r *portfolioMongoRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("student_unique_portfolio"),
		},
		{
			Keys:    bson.D{{Key: "projects.technologies", Value: 1}},
			Options: options.Index().SetName("projects_technologies_idx"),
		},
	})
	if err != nil {
		return fmt.Errorf("membuat index student_portfolios: %w", err)
	}
	return nil
}

func (r *portfolioMongoRepository) FindByStudentID(ctx context.Context, studentID int) (model.Portfolio, error) {
	var doc portfolioDocument
	err := r.coll.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Portfolio{}, ErrPortfolioNotFound
		}
		return model.Portfolio{}, fmt.Errorf("mengambil portofolio: %w", err)
	}
	return doc.toModel(), nil
}

// Upsert membuat portofolio baru jika belum ada, atau mengupdate headline/bio/certificates jika sudah ada
func (r *portfolioMongoRepository) Upsert(ctx context.Context, studentID int, req model.UpsertPortfolioRequest) (model.Portfolio, error) {
	t := nowMillis()

	certs := make([]certificateDocument, 0, len(req.Certificates))
	for _, c := range req.Certificates {
		certs = append(certs, certificateDocument{
			Name:   c.Name,
			Issuer: c.Issuer,
			Date:   c.Date.UTC().Truncate(time.Millisecond),
		})
	}

	filter := bson.M{"student_id": studentID}
	update := bson.M{
		"$set": bson.M{
			"headline":     req.Headline,
			"bio":          req.Bio,
			"certificates": certs,
			"updated_at":   t,
		},
		"$setOnInsert": bson.M{
			"student_id": studentID,
			"projects":   []projectDocument{},
			"created_at": t,
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var doc portfolioDocument
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&doc)
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("upsert portofolio: %w", err)
	}

	return doc.toModel(), nil
}

// AddProject menambah satu proyek secara atomik dengan syarat projects.19 belum ada (maks 20)
func (r *portfolioMongoRepository) AddProject(ctx context.Context, studentID int, p model.Project) (model.Portfolio, error) {
	limitField := "projects." + strconv.Itoa(MaxProjectsPerPortfolio-1)

	doc := projectDocument{
		ID:           bson.NewObjectID(),
		Title:        p.Title,
		Technologies: p.Technologies,
		Link:         p.Link,
		Year:         p.Year,
	}

	var updated portfolioDocument
	err := r.coll.FindOneAndUpdate(ctx,
		bson.M{
			"student_id": studentID,
			limitField:   bson.M{"$exists": false},
		},
		bson.M{
			"$push": bson.M{"projects": doc},
			"$set":  bson.M{"updated_at": nowMillis()},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)

	if err == nil {
		return updated.toModel(), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.Portfolio{}, fmt.Errorf("menambah proyek: %w", err)
	}

	// Cari penyebab kenapa tidak match
	count, err := r.coll.CountDocuments(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("memeriksa portofolio: %w", err)
	}
	if count == 0 {
		return model.Portfolio{}, ErrPortfolioNotFound
	}
	return model.Portfolio{}, ErrProjectLimit
}

// DeleteProject menghapus satu proyek di dalam array projects menggunakan $pull
func (r *portfolioMongoRepository) DeleteProject(ctx context.Context, studentID int, projectID string) (model.Portfolio, error) {
	oid, err := parseMongoID(projectID)
	if err != nil {
		return model.Portfolio{}, err
	}

	var updated portfolioDocument
	err = r.coll.FindOneAndUpdate(ctx,
		bson.M{
			"student_id":   studentID,
			"projects._id": oid,
		},
		bson.M{
			"$pull": bson.M{"projects": bson.M{"_id": oid}},
			"$set":  bson.M{"updated_at": nowMillis()},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)

	if err == nil {
		return updated.toModel(), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.Portfolio{}, fmt.Errorf("menghapus proyek: %w", err)
	}

	count, err := r.coll.CountDocuments(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("memeriksa portofolio: %w", err)
	}
	if count == 0 {
		return model.Portfolio{}, ErrPortfolioNotFound
	}
	return model.Portfolio{}, ErrProjectNotFound
}

// TopTechnologies menjalankan aggregation pipeline 2x $unwind untuk menghasilkan 10 teknologi terpopuler
func (r *portfolioMongoRepository) TopTechnologies(ctx context.Context) ([]model.TechnologyStat, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$projects"}},
		{{Key: "$unwind", Value: "$projects.technologies"}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$projects.technologies",
			"total": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "total", Value: -1},
			{Key: "_id", Value: 1},
		}}},
		{{Key: "$limit", Value: 10}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregasi teknologi: %w", err)
	}
	defer cursor.Close(ctx)

	var raw []struct {
		Tech  string `bson:"_id"`
		Total int    `bson:"total"`
	}
	if err := cursor.All(ctx, &raw); err != nil {
		return nil, fmt.Errorf("membaca hasil aggregasi: %w", err)
	}

	stats := make([]model.TechnologyStat, 0, len(raw))
	for _, item := range raw {
		stats = append(stats, model.TechnologyStat{
			Technology: item.Tech,
			Total:      item.Total,
		})
	}
	return stats, nil
}