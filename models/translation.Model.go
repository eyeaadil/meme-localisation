package models

import "time"

// TranslationLog represents a log of a translation request.
type TranslationLog struct {
    ID             string    `bson:"_id,omitempty" json:"id"`
    MemeID         string    `bson:"meme_id" json:"meme_id"`
    SourceText     string    `bson:"source_text" json:"source_text"`
    SourceLanguage string    `bson:"source_language" json:"source_language"`
    TargetLanguage string    `bson:"target_language" json:"target_language"`
    TranslatedText string    `bson:"translated_text" json:"translated_text"`
    CreatedAt      time.Time `bson:"created_at" json:"created_at"`
}
