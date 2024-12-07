package models

import "time"

// Analytics represents metrics for a meme.
type Analytics struct {
    ID                   string    `bson:"_id,omitempty" json:"id"`
    MemeID               string    `bson:"meme_id" json:"meme_id"`
    Views                int       `bson:"views" json:"views"`
    TranslationsRequested int       `bson:"translations_requested" json:"translations_requested"`
    Likes                int       `bson:"likes" json:"likes"`
    CreatedAt            time.Time `bson:"created_at" json:"created_at"`
}
