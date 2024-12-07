package models

import "time"

// AuthToken represents a user's authentication token.
type AuthToken struct {
    UserID    string    `bson:"user_id" json:"user_id"`
    Token     string    `bson:"token" json:"token"`
    ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
}
