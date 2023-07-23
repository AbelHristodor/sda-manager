package models

import "gorm.io/gorm"

const (
	// Hymn categories
	Standard = "Standard"
	Old = "Old"
	Others = "Others"
)

type Hymn struct {
	gorm.Model
	Title string `gorm:"type:text;default:''"`
	Verses []Verse // One-to-many relationship
	Category string `gorm:"type:text;default:''"`
	Path string `gorm:"type:text;default:''"`
	Number string
}

type Verse struct {
	gorm.Model
	HymnID uint // Foreign key
	Text string `gorm:"type:text;default:''"`
	Stanza int `gorm:"default:1"`
	IsChorus bool `gorm:"default:false"`
}