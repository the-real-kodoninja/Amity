package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hangout struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	Creator      string             `json:"creator" bson:"creator"`
	Participants []string           `json:"participants" bson:"participants"`
	Date         string             `json:"date" bson:"date"`
	Location     string             `json:"location" bson:"location"`
	CreatedAt    string             `json:"created_at" bson:"created_at"`
}
