package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type List struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Name      string             `json:"name" bson:"name"`
	Items     []ListItem         `json:"items" bson:"items"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
}

type ListItem struct {
	Type    string `json:"type" bson:"type"` // "post", "photo", "group", etc.
	ItemID  string `json:"item_id" bson:"item_id"`
	AddedAt string `json:"added_at" bson:"added_at"`
}
