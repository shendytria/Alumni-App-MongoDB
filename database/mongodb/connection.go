package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectMongo() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	// ⚠️ Sesuaikan dengan DB yang berisi koleksi `users`.
	// Dari screenshot Compass: DB = "user"
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "user" // default agar cocok dengan Compass kamu
	}

	clientOpts := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("❌ Gagal konek MongoDB:", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB tidak merespon:", err)
	}

	DB = client.Database(dbName)
	log.Println("✅ MongoDB berhasil terhubung! DB =", dbName)
}
