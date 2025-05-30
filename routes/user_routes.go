package routes

import (
	"user-api/controllers"
	"user-api/handlers"

	"github.com/gorilla/mux"
)

func UserRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/createUser", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/getAllUsers", controllers.GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/auth/google", handlers.GoogleLoginHandler).Methods("POST")


	return router
}
