package repository

import (
	"context"
	"alumni-app/app/mongodb/model"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository interface {
	Create(file *model.File) error
	FindAll() ([]model.File, error)
	FindByID(id string) (*model.File, error)
	Delete(id string) error
}

type fileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository(db *mongo.Database) FileRepository {
	return &fileRepository{collection: db.Collection("files")}
}

func (r *fileRepository) Create(file *model.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	file.UploadedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}
	file.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *fileRepository) FindAll() ([]model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var files []model.File
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileRepository) FindByID(id string) (*model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var file model.File
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
