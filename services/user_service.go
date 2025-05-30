package services

import (
	"context"
	"fmt"
	"time"
	"user-api/config"
	"user-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection = config.GetCollection("users")

func CreateUser(user models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if user.ID.IsZero() {
        user.ID = primitive.NewObjectID()
    }
    now := time.Now()
    user.CreatedAt = primitive.NewDateTimeFromTime(now)
    user.UpdatedAt = primitive.NewDateTimeFromTime(now)

    _, err := userCollection.InsertOne(ctx, user)
    if err != nil {
        fmt.Println("InsertOne error:", err)
        return err
    }
    return nil
}



func GetAllUsers() ([]models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var users []models.User
    cursor, err := userCollection.Find(ctx, bson.M{})
    if err != nil {
        return users, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var user models.User
        cursor.Decode(&user)
        users = append(users, user)
    }

    return users, nil
}

func GetUserByID(id string) (models.User, error) {
    var user models.User
    objID, _ := primitive.ObjectIDFromHex(id)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    return user, err
}

func UpdateUser(id string, user models.User) error {
    objID, _ := primitive.ObjectIDFromHex(id)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := userCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
        "$set": user,
    })

    return err
}

func DeleteUser(id string) error {
    objID, _ := primitive.ObjectIDFromHex(id)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := userCollection.DeleteOne(ctx, bson.M{"_id": objID})
    return err
}

func IsEmailExists(email string) (bool, error) {
    collection := config.GetCollection("users")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"email": email}

    var user models.User
    err := collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return false, nil // email does not exist
        }
        return false, err // some other error
    }
    return true, nil // found a user with this email
}