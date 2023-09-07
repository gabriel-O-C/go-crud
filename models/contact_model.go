package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Contact struct {
    Id       primitive.ObjectID `json:"id,omitempty"`
    Name     string             `json:"name,omitempty" validate:"required"`
    Email 	 string             `json:"email,omitempty" validate:"required"`
}