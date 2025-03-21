package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notification struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Type      string             `json:"type" bson:"type"`       // e.g., "like", "comment", "follow", "friend_request"
    From      string             `json:"from" bson:"from"`       // Username of the user who triggered the notification
    Content   string             `json:"content" bson:"content"` // Notification message
    RelatedID string             `json:"related_id" bson:"related_id"` // ID of the related post, group, etc.
    Timestamp string             `json:"timestamp" bson:"timestamp"`
    Read      bool               `json:"read" bson:"read"`
}