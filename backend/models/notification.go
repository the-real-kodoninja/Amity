package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Type      string             `json:"type" bson:"type"` // "like", "share", "comment", "friend_request", "follow", etc.
	From      string             `json:"from" bson:"from"`
	Message   string             `json:"message" bson:"message"`
	RelatedID string             `json:"related_id" bson:"related_id"` // Post ID, Group ID, etc.
	IsRead    bool               `json:"is_read" bson:"is_read"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
}
