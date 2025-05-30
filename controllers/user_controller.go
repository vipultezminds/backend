package controllers

import (
	"encoding/json"
	"net/http"
	"time"
	"user-api/models"
	"user-api/services"
	"user-api/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func CreateUser(w http.ResponseWriter, r *http.Request) {
//     err := r.ParseForm()
//     if err != nil {
//         utils.RespondError(w, http.StatusBadRequest, "Invalid form data")
//         return
//     }

//     user := models.User{
//         Name:       r.FormValue("name"),
//         Email:      r.FormValue("email"),
//         Role:       r.FormValue("role"),
//         EmployeeID: r.FormValue("employee_id"),
//     }

//     user.ID = primitive.NewObjectID()
//     now := primitive.NewDateTimeFromTime(time.Now())
//     user.CreatedAt = now
//     user.UpdatedAt = now

//     err = services.CreateUser(user)
//     if err != nil {
//         utils.RespondError(w, http.StatusInternalServerError, "User creation failed")
//         return
//     }

//     utils.RespondJSON(w, http.StatusCreated, user)
// }

func CreateUser(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        utils.RespondError(w, http.StatusBadRequest, "Invalid form data")
        return
    }

    user := models.User{
        Name:       r.FormValue("name"),
        Email:      r.FormValue("email"),
        Role:       r.FormValue("role"),
        EmployeeID: r.FormValue("employee_id"),
    }

    // Check if email already exists
    exists, err := services.IsEmailExists(user.Email)
    if err != nil {
        utils.RespondError(w, http.StatusInternalServerError, "Failed to check email existence")
        return
    }
    if exists {
        utils.RespondError(w, http.StatusConflict, "Email already exists")
        return
    }

    user.ID = primitive.NewObjectID()
    now := primitive.NewDateTimeFromTime(time.Now())
    user.CreatedAt = now
    user.UpdatedAt = now

    err = services.CreateUser(user)
    if err != nil {
        utils.RespondError(w, http.StatusInternalServerError, "User creation failed")
        return
    }

    utils.RespondJSON(w, http.StatusCreated, user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := services.GetAllUsers()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Could not fetch users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := services.GetUserByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	err := services.UpdateUser(id, user)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Update failed")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "User updated"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := services.DeleteUser(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Delete failed")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "User deleted"})
}
