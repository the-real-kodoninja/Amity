package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string             `json:"name" bson:"name"`
    Description string             `json:"description" bson:"description"`
    Creator     string             `json:"creator" bson:"creator"`
    Members     []string           `json:"members" bson:"members"`
    Posts       []Post             `json:"posts" bson:"posts"`
    CreatedAt   string             `json:"created_at" bson:"created_at"`
}