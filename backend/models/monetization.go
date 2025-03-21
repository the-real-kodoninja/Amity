package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Monetization struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username      string             `json:"username" bson:"username"`
	TotalEarnings float64            `json:"total_earnings" bson:"total_earnings"`
	AdEarnings    float64            `json:"ad_earnings" bson:"ad_earnings"`
	NFTEarnings   float64            `json:"nft_earnings" bson:"nft_earnings"`
}
