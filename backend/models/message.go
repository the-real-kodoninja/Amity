package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	From      string             `json:"from" bson:"from"`
	To        string             `json:"to" bson:"to"`
	Content   string             `json:"content" bson:"content"`
	IsAI      bool               `json:"is_ai" bson:"is_ai"` // For AI IM
	CreatedAt string             `json:"created_at" bson:"created_at"`
}
