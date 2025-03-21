package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdminMessage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	From      string             `json:"from" bson:"from"`
	Content   string             `json:"content" bson:"content"`
	Timestamp string             `json:"timestamp" bson:"timestamp"`
	Read      bool               `json:"read" bson:"read"`
}

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	From      string             `bson:"from"`
	To        string             `bson:"to"`
	Content   string             `bson:"content"`
	Timestamp string             `bson:"timestamp"`
}
