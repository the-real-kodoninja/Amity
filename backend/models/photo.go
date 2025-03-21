package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Photo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Caption     string             `json:"caption" bson:"caption"`
	Timestamp   string             `json:"timestamp" bson:"timestamp"`
	IsMinted    bool               `json:"is_minted" bson:"is_minted"`
	MintedToken string             `json:"minted_token" bson:"minted_token"` // Simulated NFT token ID
}
