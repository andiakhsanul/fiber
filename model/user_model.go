package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Username string             `json:"username" validate:"required"`
    Name     string             `json:"name" validate:"required"`
    Location string             `json:"location" validate:"required"`
    Title    string             `json:"title" validate:"required"`
    Password string             `json:"password" validate:"required"`
}

