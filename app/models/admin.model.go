package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdminModel struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Username string             `json:"username,omitempty" validate:"required"`
	Password string             `json:"passwrod,omitempty" validate:"required"`
}
