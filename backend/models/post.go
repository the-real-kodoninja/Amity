package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	Username  string `json:"username" bson:"username"`
	Content   string `json:"content" bson:"content"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
}

type Post struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ActivityID   string             `json:"activity_id" bson:"activity_id"`
	ActivityType string             `json:"activity_type" bson:"activity_type"`
	Username     string             `json:"username" bson:"username"`
	Content      string             `json:"content" bson:"content"`
	ImageURL     string             `json:"image_url" bson:"image_url"`
	Hashtags     []string           `json:"hashtags" bson:"hashtags"`
	Timestamp    string             `json:"timestamp" bson:"timestamp"`
	Likes        int                `json:"likes" bson:"likes"`
	Shares       int                `json:"shares" bson:"shares"`
	Comments     []Comment          `json:"comments" bson:"comments"`
	IsNSFW       bool               `json:"is_nsfw" bson:"is_nsfw"`
	IsMinted     bool               `json:"is_minted" bson:"is_minted"`
	MintedToken  string             `json:"minted_token" bson:"minted_token"` // Simulated NFT token ID
}
