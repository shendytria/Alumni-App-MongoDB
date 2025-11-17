package repository

import (
	"alumni-app/app/mongodb/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface {
	GetByUsername(username string) (model.User, error)
	GetByID(id interface{}) (model.User, error)
	Create(user *model.User) error
}

type UserRepository struct {
	Col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepositoryInterface {
	return &UserRepository{
		Col: db.Collection("users"),
	}
}

func (r *UserRepository) GetByUsername(username string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := r.Col.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return user, err
}

func (r *UserRepository) GetByID(id interface{}) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (r *UserRepository) Create(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	_, err := r.Col.InsertOne(ctx, user)
	return err
}
