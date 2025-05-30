package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"user-api/models"
	"user-api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/idtoken"
)

var UserCollection *mongo.Collection

type GoogleAuthRequest struct {
	IDToken string `json:"id_token"`
}

type GoogleAuthResponse struct {
	User models.User `json:"user"`
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req GoogleAuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.IDToken) == "" {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		utils.RespondError(w, http.StatusInternalServerError, "Missing GOOGLE_CLIENT_ID env variable")
		return
	}

	ctx := context.Background()
	payload, err := idtoken.Validate(ctx, req.IDToken, clientID)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid Google ID token")
		return
	}

	email, ok1 := payload.Claims["email"].(string)
	name, ok2 := payload.Claims["name"].(string)
	if !ok1 || !ok2 {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to parse user info from token")
		return
	}

	// Check if user exists
	var user models.User
	err = UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		// New user: create
		user = models.User{
			ID:         primitive.NewObjectID(),
			Name:       name,
			Email:      email,
			Role:       "employee",
			Projects:   []primitive.ObjectID{},
			EmployeeID: "",
			CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt:  primitive.NewDateTimeFromTime(time.Now()),
		}

		_, err = UserCollection.InsertOne(ctx, user)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "Error creating user")
			return
		}
	} else if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Database error")
		return
	}

	utils.RespondJSON(w, http.StatusOK, GoogleAuthResponse{User: user})
}

