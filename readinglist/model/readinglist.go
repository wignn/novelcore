package model

import "time"

type ReadingListEntry struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"account_id"`
	NovelID        string    `json:"novel_id"`
	Status         string    `json:"status"`
	CurrentChapter float64   `json:"current_chapter"`
	Rating         int32     `json:"rating"`
	Notes          string    `json:"notes"`
	IsFavorite     bool      `json:"is_favorite"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
