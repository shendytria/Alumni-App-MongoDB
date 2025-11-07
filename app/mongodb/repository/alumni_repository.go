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

type AlumniRepository struct{}

func NewAlumniRepository() *AlumniRepository {
	return &AlumniRepository{}
}

func (r *AlumniRepository) GetAll(search string, sortBy string, order string, page, limit int) ([]model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"nama": bson.M{"$regex": search, "$options": "i"}},
			{"jurusan": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		},
	}
	// exclude soft-deleted
	filter["$and"] = []bson.M{{"deleted_at": bson.M{"$exists": false}}}

	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}

	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cur, err := database.DB.Collection("alumni").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Alumni
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *AlumniRepository) Count(search string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"nama": bson.M{"$regex": search, "$options": "i"}},
			{"jurusan": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		},
	}
	filter["$and"] = []bson.M{{"deleted_at": bson.M{"$exists": false}}}

	return database.DB.Collection("alumni").CountDocuments(ctx, filter)
}

func (r *AlumniRepository) GetByID(id primitive.ObjectID) (model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var a model.Alumni
	err := database.DB.Collection("alumni").
		FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$exists": false}}).
		Decode(&a)
	return a, err
}

func (r *AlumniRepository) GetByUserID(userID primitive.ObjectID) (model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var a model.Alumni
	err := database.DB.Collection("alumni").
		FindOne(ctx, bson.M{"user_id": userID, "deleted_at": bson.M{"$exists": false}}).
		Decode(&a)
	return a, err
}

func (r *AlumniRepository) Create(a *model.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.ID = primitive.NewObjectID()
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	_, err := database.DB.Collection("alumni").InsertOne(ctx, a)
	return err
}

func (r *AlumniRepository) Update(id primitive.ObjectID, a *model.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.UpdatedAt = time.Now()
	update := bson.M{"$set": a}

	_, err := database.DB.Collection("alumni").UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *AlumniRepository) SoftDelete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	update := bson.M{"$set": bson.M{"deleted_at": now}}

	_, err := database.DB.Collection("alumni").UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *AlumniRepository) GetAllByUserID(userID primitive.ObjectID) ([]model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := database.DB.Collection("alumni").Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Alumni
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// func (r *AlumniRepository) UpdateFieldByID(ctx context.Context, id primitive.ObjectID, field string, value any) error {
// 	_, err := database.DB.Collection("alumni").UpdateOne(
// 		ctx,
// 		bson.M{"_id": id},
// 		bson.M{"$set": bson.M{
// 			field:       value,
// 			"updated_at": time.Now(),
// 		}},
// 	)
// 	return err
// }

// func (r *AlumniRepository) UpdateFieldByHex(idHex, field string, value any) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	oid, err := primitive.ObjectIDFromHex(idHex)
// 	if err != nil {
// 		return err
// 	}
// 	return r.UpdateFieldByID(ctx, oid, field, value)
// }