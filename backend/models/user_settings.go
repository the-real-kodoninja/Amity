package models

type UserSettings struct {
	NSFWEnabled       bool   `json:"nsfw_enabled" bson:"nsfw_enabled"`
	ProfileVisibility string `json:"profile_visibility" bson:"profile_visibility"` // "public", "friends", "private"
	Messaging         string `json:"messaging" bson:"messaging"`                   // "everyone", "friends", "none"
}
