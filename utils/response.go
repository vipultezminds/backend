package utils

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// APIResponse defines the standard API response structure
type APIResponse struct {
    Data      any `json:"data"`
    ErrorMsg  string      `json:"errorMsg"`
    ErrorCode int         `json:"errorCode"`
}

// RespondJSON sends a JSON response with the standard API response format
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
    response := APIResponse{
        Data:      data,
        ErrorMsg:  "",
        ErrorCode: 0,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}

// RespondError sends an error response with the standard API response format
func RespondError(w http.ResponseWriter, status int, message string) {
    response := APIResponse{
        Data:      nil,
        ErrorMsg:  message,
        ErrorCode: 1,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}

// Call this after you connect to MongoDB
func CreateUserEmailUniqueIndex(userCollection *mongo.Collection) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    indexModel := mongo.IndexModel{
        Keys:    bson.D{{Key: "email", Value: 1}}, // index on email
        Options: options.Index().SetUnique(true),
    }

    _, err := userCollection.Indexes().CreateOne(ctx, indexModel)
    if err != nil {
        return err
    }

    log.Println("✅ Unique index created on 'email' field")
    return nil
}
