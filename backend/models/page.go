package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Page struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string             `json:"name" bson:"name"`
    Description string             `json:"description" bson:"description"`
    Creator     string             `json:"creator" bson:"creator"`
    Followers   []string           `json:"followers" bson:"followers"`
    Posts       []Post             `json:"posts" bson:"posts"`
    CreatedAt   string             `json:"created_at" bson:"created_at"`
}