package models

import "time"

// Translation represents a single translated text for a meme.
type Translation struct {
    Language       string `bson:"language" json:"language"`
    TranslatedText string `bson:"translated_text" json:"translated_text"`
}

// LocalizedImage represents a translated image for a meme.
type LocalizedImage struct {
    Language  string `bson:"language" json:"language"`
    ImagePath string `bson:"image_path" json:"image_path"`
}

// Meme represents a meme uploaded by a user.
type Meme struct {
    ID              string           `bson:"_id,omitempty" json:"id"`
    UserID          string           `bson:"user_id" json:"user_id"`
    OriginalImage   string           `bson:"original_image" json:"original_image"`
    ExtractedText   string           `bson:"extracted_text" json:"extracted_text"`
    Translations    []Translation    `bson:"translations" json:"translations"`
    LocalizedImages []LocalizedImage `bson:"localized_images" json:"localized_images"`
    CreatedAt       time.Time        `bson:"created_at" json:"created_at"`
}
