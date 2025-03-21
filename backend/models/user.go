package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username         string             `json:"username" bson:"username"`
	Email            string             `json:"email" bson:"email"`
	Password         string             `json:"password" bson:"password"`
	Location         string             `json:"location" bson:"location"`
	Followers        int                `json:"followers" bson:"followers"`
	Following        []string           `json:"following" bson:"following"`
	Friends          []string           `json:"friends" bson:"friends"`
	BlockedUsers     []string           `json:"blocked_users" bson:"blocked_users"`
	ProfilePhoto     string             `json:"profile_photo" bson:"profile_photo"`
	BannerPhoto      string             `json:"banner_photo" bson:"banner_photo"`
	Settings         UserSettings       `json:"settings" bson:"settings"`
	Notifications    []Notification     `json:"notifications" bson:"notifications"`
	Verified         bool               `json:"verified" bson:"verified"`
	IsAdmin          bool               `json:"is_admin" bson:"is_admin"`
	Banned           bool               `json:"banned" bson:"banned"`
	PinnedPostID     string             `json:"pinned_post_id" bson:"pinned_post_id"`
	TotalNFTEarnings float64            `json:"total_nft_earnings" bson:"total_nft_earnings"` // New field
}
