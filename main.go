package main

import (
	"log"
	"net/http"
	"os"
	"user-api/config"
	"user-api/handlers"
	"user-api/routes"
	"user-api/utils"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to DB before using collection
	config.ConnectDB()

	// Get users collection
	userCollection := config.GetCollection("users")

	handlers.UserCollection = userCollection

	// Create unique index on email field
	err = utils.CreateUserEmailUniqueIndex(userCollection)
	if err != nil {
		log.Fatal("Failed to create unique index:", err)
	}

	router := routes.UserRoutes()

	// Setup CORS options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // your React app origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap your router with CORS middleware
	handler := corsHandler.Handler(router)

	port := os.Getenv("PORT")
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
