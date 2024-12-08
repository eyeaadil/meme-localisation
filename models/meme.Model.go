package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Translation represents a single translated text for a meme.
type Translation struct {
	Language       string `bson:"language" json:"language" binding:"required"`
	TranslatedText string `bson:"translated_text" json:"translated_text" binding:"required"`
}

// LocalizedImage represents a translated image for a meme.
type LocalizedImage struct {
	Language  string `bson:"language" json:"language" binding:"required"`
	ImagePath string `bson:"image_path" json:"image_path" binding:"required"`
}

// Meme represents a meme uploaded by a user.
type Meme struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	OriginalImage   string             `bson:"original_image" json:"original_image" binding:"required"`
	ExtractedText   string             `bson:"extracted_text" json:"extracted_text,omitempty"`
	Translations    []Translation      `bson:"translations" json:"translations,omitempty"`
	LocalizedImages []LocalizedImage   `bson:"localized_images" json:"localized_images,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
}

// MemeCreateRequest provides a strict input validation structure
type MemeCreateRequest struct {
	UserID          string             `json:"user_id" binding:"required"`
	OriginalImage   string             `json:"original_image" binding:"required,url"`
	ExtractedText   string             `json:"extracted_text,omitempty"`
	Translations    []Translation      `json:"translations,omitempty" binding:"omitempty,dive"`
	LocalizedImages []LocalizedImage   `json:"localized_images,omitempty" binding:"omitempty,dive"`
}