package model

import "time"

type Review struct {
	ID        string    `json:"id"`
	NovelID   string    `json:"novel_id"`
	AccountID string    `json:"account_id"`
	Rating    int       `json:"rating"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsSpoiler bool      `json:"is_spoiler"`
	Upvotes   int32     `json:"upvotes"`
	Downvotes int32     `json:"downvotes"`
	CreatedAt time.Time `json:"created_at"`
}