package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Activity struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username   string             `json:"username" bson:"username"`
	ActivityID string             `json:"activity_id" bson:"activity_id"`
	Type       string             `json:"type" bson:"type"`
	Actor      string             `json:"actor" bson:"actor"`
	Object     string             `json:"object" bson:"object"`
	ReceivedAt string             `json:"received_at" bson:"received_at"`
}
