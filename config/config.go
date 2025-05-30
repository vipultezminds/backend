package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	ctx := context.TODO()

    fmt.Println("MONGO_URI",os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://vipul:vipul010399@cluster0.wz5k24d.mongodb.net/"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Optional: Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	DB = client.Database("userdb")
	log.Println("✅ MongoDB connected")
}

func GetCollection(collectionName string) *mongo.Collection {
    fmt.Println(collectionName)
	if DB == nil {
		// log.Fatal("❌ MongoDB not connected. Call ConnectDB() first.")
        ConnectDB()
	}
	return DB.Collection(collectionName)
}
