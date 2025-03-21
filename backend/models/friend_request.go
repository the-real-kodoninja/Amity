package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FriendRequest struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	From      string             `bson:"from"`
	To        string             `bson:"to"`
	Status    string             `json:"status" bson:"status"` // "pending", "accepted", "rejected"
	CreatedAt string             `json:"created_at" bson:"created_at"`
	Timestamp string             `bson:"timestamp"`
}
