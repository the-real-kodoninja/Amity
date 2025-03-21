package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username   string             `json:"username" bson:"username"`
	Content    string             `json:"content" bson:"content"` // Limited to 280 characters
	Media      []Media            `json:"media" bson:"media"`     // Photos, videos, files
	Likes      int                `json:"likes" bson:"likes"`
	Reactions  map[string]int     `json:"reactions" bson:"reactions"` // e.g., {"heart": 5, "kiss": 3}
	Shares     int                `json:"shares" bson:"shares"`
	Comments   []Comment          `json:"comments" bson:"comments"`
	Timestamp  string             `json:"timestamp" bson:"timestamp"`
	HiddenBy   []string           `json:"hidden_by" bson:"hidden_by"`   // Users who hid this post
	IsShort    bool               `json:"is_short" bson:"is_short"`     // Indicates if this is a short
	Visibility string             `json:"visibility" bson:"visibility"` // "public", "friends", "private"
}

type Media struct {
	Type string `json:"type" bson:"type"` // "photo", "video", "file"
	URL  string `json:"url" bson:"url"`
	Size int64  `json:"size" bson:"size"` // Size in bytes
}

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Content   string             `json:"content" bson:"content"`
	Timestamp string             `json:"timestamp" bson:"timestamp"`
	Emoji     string             `json:"emoji" bson:"emoji"`
	GIF       string             `json:"gif" bson:"gif"`         // URL of the GIF
	Replies   []Comment          `json:"replies" bson:"replies"` // Nested replies
}
