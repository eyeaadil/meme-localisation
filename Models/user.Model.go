package models

import "time"

// User represents a user of the platform.
type User struct {
    ID               string    `bson:"_id,omitempty" json:"id"`
    Name             string    `bson:"name" json:"name"`
    Email            string    `bson:"email" json:"email"`
    Password         string    `bson:"password" json:"-"` // Hashed password
    PreferredLanguage string    `bson:"preferred_language" json:"preferred_language"`
    CreatedAt        time.Time `bson:"created_at" json:"created_at"`
}
