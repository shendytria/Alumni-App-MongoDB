package repository

import (
	"alumni-app/app/mongodb/model"
	"alumni-app/database/mongodb"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PekerjaanRepository struct{}

func NewPekerjaanRepository() *PekerjaanRepository { return &PekerjaanRepository{} }

func (r *PekerjaanRepository) GetAll(search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
		},
	}
	filter["$and"] = []bson.M{{"deleted_at": bson.M{"$exists": false}}}

	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}
	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cur, err := database.DB.Collection("pekerjaan").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result []model.PekerjaanAlumni
	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PekerjaanRepository) Count(search string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
		},
	}
	filter["$and"] = []bson.M{{"deleted_at": bson.M{"$exists": false}}}

	return database.DB.Collection("pekerjaan").CountDocuments(ctx, filter)
}

func (r *PekerjaanRepository) GetByID(id primitive.ObjectID) (model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var p model.PekerjaanAlumni
	err := database.DB.Collection("pekerjaan").FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	return p, err
}

func (r *PekerjaanRepository) GetByAlumniID(alumniID primitive.ObjectID) ([]model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := database.DB.Collection("pekerjaan").Find(ctx, bson.M{
		"alumni_id":  alumniID,
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.PekerjaanAlumni
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PekerjaanRepository) Create(p *model.PekerjaanAlumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p.ID = primitive.NewObjectID()
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	_, err := database.DB.Collection("pekerjaan").InsertOne(ctx, p)
	return err
}

func (r *PekerjaanRepository) Update(id primitive.ObjectID, p *model.PekerjaanAlumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p.UpdatedAt = time.Now()
	update := bson.M{"$set": p}

	_, err := database.DB.Collection("pekerjaan").UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *PekerjaanRepository) SoftDelete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	update := bson.M{"$set": bson.M{"deleted_at": now}}
	_, err := database.DB.Collection("pekerjaan").UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *PekerjaanRepository) GetTrashed() ([]model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := database.DB.Collection("pekerjaan").Find(ctx, bson.M{"deleted_at": bson.M{"$exists": true}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.PekerjaanAlumni
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PekerjaanRepository) GetTrashedByAlumniIDs(alumniIDs []primitive.ObjectID) ([]model.PekerjaanAlumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"alumni_id":  bson.M{"$in": alumniIDs},
		"deleted_at": bson.M{"$exists": true},
	}

	cur, err := database.DB.Collection("pekerjaan").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.PekerjaanAlumni
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PekerjaanRepository) Restore(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$unset": bson.M{"deleted_at": ""}}
	_, err := database.DB.Collection("pekerjaan").UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *PekerjaanRepository) HardDelete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := database.DB.Collection("pekerjaan").DeleteOne(ctx, bson.M{"_id": id})
	return err
}
