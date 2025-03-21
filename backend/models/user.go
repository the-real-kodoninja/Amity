package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username      string             `json:"username" bson:"username"`
	Email         string             `json:"email" bson:"email"`
	Password      string             `json:"password" bson:"password"` // Hashed password
	Location      string             `json:"location" bson:"location"`
	Followers     int                `json:"followers" bson:"followers"`
	Following     []string           `json:"following" bson:"following"`
	Friends       []string           `json:"friends" bson:"friends"`
	BlockedUsers  []string           `json:"blocked_users" bson:"blocked_users"`
	ProfilePhoto  string             `json:"profile_photo" bson:"profile_photo"`
	BannerPhoto   string             `json:"banner_photo" bson:"banner_photo"`
	NSFWSettings  bool               `json:"nsfw_settings" bson:"nsfw_settings"` // true = show NSFW content
	Notifications []Notification     `json:"notifications" bson:"notifications"`
}

type Notification struct {
	ID        string `json:"id" bson:"id"`
	Type      string `json:"type" bson:"type"` // e.g., "like", "comment", "follow", "friend_request"
	From      string `json:"from" bson:"from"` // Username of the user who triggered the notification
	Content   string `json:"content" bson:"content"`
	Timestamp string `json:"timestamp" bson:"timestamp"`
	Read      bool   `json:"read" bson:"read"`
}
