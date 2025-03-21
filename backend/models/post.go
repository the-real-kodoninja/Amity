package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Content      string             `json:"content" bson:"content"`
	Media        []Media            `json:"media" bson:"media"`
	Timestamp    string             `json:"timestamp" bson:"timestamp"`
	Likes        int                `json:"likes" bson:"likes"`
	Reactions    map[string]int     `json:"reactions" bson:"reactions"`
	Shares       int                `json:"shares" bson:"shares"`
	Comments     []Comment          `json:"comments" bson:"comments"`
	HiddenBy     []string           `json:"hidden_by" bson:"hidden_by"`
	IsShort      bool               `json:"is_short" bson:"is_short"`
	Deleted      bool               `json:"deleted" bson:"deleted"`
	Sponsored    bool               `json:"sponsored" bson:"sponsored"`
	Live         bool               `json:"live" bson:"live"`
	NFTAddress   string             `json:"nft_address" bson:"nft_address"`     // New field
	MintEarnings float64            `json:"mint_earnings" bson:"mint_earnings"` // New field
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
