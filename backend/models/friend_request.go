package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FriendRequest struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	From      string             `json:"from" bson:"from"`
	To        string             `json:"to" bson:"to"`
	Status    string             `json:"status" bson:"status"` // "pending", "accepted", "rejected"
	CreatedAt string             `json:"created_at" bson:"created_at"`
}
