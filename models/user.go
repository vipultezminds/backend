package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	Email      string               `json:"email" bson:"email"`
	Role       string               `json:"role" bson:"role"`
	Projects   []primitive.ObjectID `bson:"projects" json:"projects"`
    EmployeeID string             `json:"employee_id" bson:"employee_id"`
	CreatedAt  primitive.DateTime   `bson:"created_at"`
	UpdatedAt  primitive.DateTime   `bson:"updated_at"`
}
