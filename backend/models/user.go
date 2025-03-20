package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"password"` // Hashed password
	ProfilePic  string             `bson:"profilePic" json:"profilePic"`
	CoverPhoto  string             `bson:"coverPhoto" json:"coverPhoto"`
	Location    string             `bson:"location" json:"location"`
	LastActive  time.Time          `bson:"lastActive" json:"lastActive"`
	Connections []string           `bson:"connections" json:"connections"`
	Followers   []string           `bson:"followers" json:"followers"`
	Following   []string           `bson:"following" json:"following"`
	Photos      []string           `bson:"photos" json:"photos"`
	Posts       []string           `bson:"posts" json:"posts"`
	IsAnonymous bool               `bson:"isAnonymous" json:"isAnonymous"`
}
